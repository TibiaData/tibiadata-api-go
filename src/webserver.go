package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/mantyr/go-charset/data"
	"golang.org/x/text/encoding/charmap"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// TibiadataDefaultVoc - default vocation when not specified in request
var TibiadataDefaultVoc string = "all"

// Tibiadata app flags for running
var TibiadataAPIversion int = 3
var TibiadataDebug bool

// Tibiadata app user-agent
var TibiadataUserAgent string

// Tibiadata app details set to release/build on GitHub
var TibiadataBuildRelease = "unknown"     // will be set by GitHub Actions (to release number)
var TibiadataBuildBuilder = "manual"      // will be set by GitHub Actions
var TibiadataBuildCommit = "-"            // will be set by GitHub Actions (to git commit)
var TibiadataBuildEdition = "open-source" //

// Information - child of JSONData
type Information struct {
	APIVersion int    `json:"api_version"`
	Timestamp  string `json:"timestamp"`
}

func main() {
	// logging start of TibiaData
	log.Printf("[info] TibiaData API starting..")

	// running the TibiaDataInitializer function
	TibiaDataInitializer()

	// logging build information
	log.Printf("[info] TibiaData API release: %s", TibiadataBuildRelease)
	log.Printf("[info] TibiaData API build: %s", TibiadataBuildBuilder)
	log.Printf("[info] TibiaData API commit: %s", TibiadataBuildCommit)
	log.Printf("[info] TibiaData API edition: %s", TibiadataBuildEdition)

	// setting gin-application to certain mode if GIN_MODE is set to release, test or debug (default is release)
	switch ginMode := getEnv("GIN_MODE", "release"); ginMode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	// logging the gin.mode
	log.Printf("[info] TibiaData API gin-mode: %s", gin.Mode())

	// setting tibiadata-application to log much less if DEBUG_MODE is false (default is false)
	if !getEnvAsBool("DEBUG_MODE", false) {
		log.Printf("[info] TibiaData API debug-mode: disabled")
	} else {
		// setting debug to true for more logging
		TibiadataDebug = true
		log.Printf("[info] TibiaData API debug-mode: enabled")

		// logging user-agent string
		log.Printf("[debug] TIbiaData API User-Agent: %s", TibiadataUserAgent)
	}

	router := gin.Default()

	// disable proxy feature of gin
	router.SetTrustedProxies(nil)

	// Ping-endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// TibiaData API version 3
	v3 := router.Group("/v3")
	{
		// Tibia characters
		v3.GET("/characters/character/:character", TibiaCharactersCharacterV3)

		// Tibia creatures
		v3.GET("/creatures", TibiaCreaturesOverviewV3)
		v3.GET("/creatures/creature/:race", TibiaCreaturesCreatureV3)

		// Tibia fansites
		v3.GET("/fansites", TibiaFansitesV3)

		// Tibia highscores
		v3.GET("/highscores/world/:world", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/experience/"+TibiadataDefaultVoc)
		})
		v3.GET("/highscores/world/:world/:category", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/"+c.Param("category")+"/"+TibiadataDefaultVoc)
		})
		v3.GET("/highscores/world/:world/:category/:vocation", TibiaHighscoresV3)

		// Tibia killstatistics
		v3.GET("/killstatistics/world/:world", TibiaKillstatisticsV3)

		// Tibia spells
		v3.GET("/spells", TibiaSpellsOverviewV3)
		v3.GET("/spells/spell/:spell", TibiaSpellsSpellV3)
		v3.GET("/spells/vocation/:vocation", TibiaSpellsOverviewV3)

		// Tibia worlds
		v3.GET("/worlds", TibiaWorldsOverviewV3)
		v3.GET("/worlds/world/:world", TibiaWorldsWorldV3)
	}

	// container version details endpoint
	router.GET("/versions", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"release": TibiadataBuildRelease,
			"build":   TibiadataBuildBuilder,
			"commit":  TibiadataBuildCommit,
			"edition": TibiadataBuildEdition,
		})
	})

	// Start the router
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// TibiaDataInitializer func - init things at beginning
func TibiaDataInitializer() {
	// setting TibiadataBuildEdition
	if isEnvExist("TIBIADATA_EDITION") {
		TibiadataBuildEdition = getEnv("TIBIADATA_EDITION", "open-source")
	}

	// generating TibiadataUserAgent with TibiadataUserAgentGenerator function
	TibiadataUserAgent = TibiadataUserAgentGenerator(TibiadataAPIversion)
}

/*
// TibiaDataAPIHandleErrorResponse func - handling of responses..
func TibiaDataAPIHandleErrorResponse(c *gin.Context, s1 string, s2 string, s3 string) {
	if TibiadataDebug {
		log.Println("[error] " + s1 + " - (" + c.Request.RequestURI + "). " + s2 + "; " + s3)
	}

	// return error response
	c.JSON(http.StatusOK, gin.H{"error": s2})
}

// TibiaDataAPIHandleOtherResponse func - handling of responses..
func TibiaDataAPIHandleOtherResponse(c *gin.Context, httpCode int, s string, j interface{}) {
	if TibiadataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response (with specific status code)
	c.JSON(httpCode, j)
}
*/

// TibiaDataAPIHandleSuccessResponse func - handling of responses..
func TibiaDataAPIHandleSuccessResponse(c *gin.Context, s string, j interface{}) {
	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] " + s + " - (" + c.Request.RequestURI + ") returned data:")
		js, _ := json.Marshal(j)
		log.Printf("[debug] %s\n", js)
	}

	if TibiadataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response
	c.JSON(http.StatusOK, j)
}

// TibiadataUserAgentGenerator func - creates User-Agent for requests
func TibiadataUserAgentGenerator(version int) string {

	// setting product name
	useragent := "TibiaData-API/v" + strconv.Itoa(version)

	// adding information of host
	TibiadataHost := getEnv("TIBIADATA_UA_HOSTNAME", "")
	if TibiadataHost != "" {
		TibiadataHost = "+https://" + TibiadataHost
	}

	// adding details in parenthesis
	useragentDetails := []string{
		"release/" + TibiadataBuildRelease,
		"build/" + TibiadataBuildBuilder,
		"commit/" + TibiadataBuildCommit,
		"edition/" + TibiadataBuildEdition,
		TibiadataHost,
	}
	useragent += " (" + strings.Join(useragentDetails, "; ") + ")"

	return useragent
}

// TibiadataHTMLDataCollectorV3 func
func TibiadataHTMLDataCollectorV3(TibiaURL string) string {

	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiadataDebug var
	if TibiadataDebug {
		client.SetDebug(true)
	}

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   TibiadataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	res, err := client.R().Get(TibiaURL)
	if err != nil {
		log.Printf("[error] TibiadataHTMLDataCollectorV3 (URL: %s) in resp1: %s", TibiaURL, err)
	}

	// Checking if status is something else than 200
	if res.StatusCode() != 200 {
		log.Printf("[warni] TibiadataHTMLDataCollectorV3 (URL: %s) status code: %s", TibiaURL, res.Status())

		// Check if page is in maintenance mode
		if res.StatusCode() == 302 {
			log.Printf("[info] TibiadataHTMLDataCollectorV3 (URL: %s): Page tibia.com returns 302, probably maintenance mode enabled. ", TibiaURL)

			// TODO
			// do response with maintenance mode..
		}
	}

	// Convert string to io.Reader
	res_io := strings.NewReader(res.String())

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res_io)
	if err != nil {
		log.Printf("[error] TibiadataHTMLDataCollectorV3 (URL: %s) error: %s", TibiaURL, err)
	}

	// Find of this to get div with class BoxContent
	data, err := doc.Find(".Border_2 .Border_3").Html()
	if err != nil {
		log.Fatal(err)
	}

	// convert string from eg "&nbsp;" to " "
	data = html.UnescapeString(data)
	data = strings.ReplaceAll(data, "&nbsp;", " ")

	// convert string from ISO 8859-1 to UTF-8
	data, _ = TibiaDataConvertEncodingtoUTF8(data)

	// Return of extracted html to functions..
	return string(data)
}

// TibiadataHTMLRemoveLinebreaksV3 func
func TibiadataHTMLRemoveLinebreaksV3(data string) string {
	return string(strings.ReplaceAll(data, "\n", ""))
}

// TibiadataRemoveURLsV3 func
func TibiadataRemoveURLsV3(data string) string {
	// prepare return value
	var returnData string

	// convert string from UTF8 to ISO88591
	data, _ = TibiaDataConvertEncodingtoISO88591(data)

	// Regex to remove URLs
	regex := regexp.MustCompile(`<a.*>(.*)<\/a>`)
	result := regex.FindAllStringSubmatch(data, -1)

	if len(result) > 0 {
		returnData = result[0][1]
	} else {
		returnData = ""
	}
	return string(returnData)
}

// TibiadataStringWorldFormatToTitleV3 func
func TibiadataStringWorldFormatToTitleV3(world string) string {
	return string(strings.Title(strings.ToLower(world)))
}

// TibiadataUnescapeStringV3 func
func TibiadataUnescapeStringV3(data string) string {
	//	data, _ = TibiaDataConvertEncodingtoUTF8(data)
	return string(html.UnescapeString(data))
}

// TibiadataQueryEscapeStringV3 func
func TibiadataQueryEscapeStringV3(data string) string {
	data, _ = TibiaDataConvertEncodingtoISO88591(data)
	return string(url.QueryEscape(data))
}

// TibiadataDatetimeV3 func
func TibiadataDatetimeV3(date string) string {

	var returnDate string

	// we need to use TibiaDataConvertEncodingtoISO88591 so that the parser doens't complain
	date, _ = TibiaDataConvertEncodingtoISO88591(date)

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now().Format(time.RFC3339)
	} else {
		// Converting: Jan 02 2007, 19:20:30 CET -> RFC1123 -> RFC3339

		// regex to exact values..
		regex1 := regexp.MustCompile(`(.*).([0-9][0-9]).([0-9][0-9][0-9][0-9]),.([0-9][0-9]:[0-9][0-9]:[0-9][0-9]).(.*)`)
		subma1 := regex1.FindAllStringSubmatch(date, -1)

		if len(subma1) > 0 {
			// Adding fake-Sun for valid RFC1123 convertion..
			dateDate, err := time.Parse(time.RFC1123, "Sun, "+subma1[0][2]+" "+subma1[0][1]+" "+subma1[0][3]+" "+subma1[0][4]+" "+subma1[0][5])
			if err != nil {
				// log.Fatal(err)
				log.Println(err)
			}

			// Set data to return
			returnDate = dateDate.Format(time.RFC3339)

		} else {
			// Format not defined yet..
			log.Println("Incoming date: " + date)
			log.Println("UNKNOWN FORMAT YET!")

			// Parse the given string to be formatted correct later
			dateDate, err := time.Parse(time.RFC3339, string(date))
			if err != nil {
				log.Fatal(err)
			}

			// Set data to return
			returnDate = dateDate.Format(time.RFC3339)

		}
	}

	// Return of formatted date and time string to functions..
	return returnDate

}

// TibiadataDateV3 func
func TibiadataDateV3(date string) string {
	// we need to use TibiaDataConvertEncodingtoISO88591 so that the parser doens't complain
	date, _ = TibiaDataConvertEncodingtoISO88591(date)

	// use regex to skip weird formatting on "spaces"
	regex1 := regexp.MustCompile(`([a-zA-Z]{3}).*([0-9]{2}).*([0-9]{4})`)
	subma1 := regex1.FindAllStringSubmatch(date, -1)
	date = (subma1[0][1] + " " + subma1[0][2] + " " + subma1[0][3])

	// parsing and setting format of return
	tmpDate, _ := time.Parse("Jan 02 2006", date)
	date = tmpDate.Format("2006-01-02")

	return date
}

// TibiadataStringToIntegerV3 func
func TibiadataStringToIntegerV3(data string) int {

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(data, "")
	returnData, _ := strconv.Atoi(processedString)

	// Return of formatted date and time string to functions..
	return int(returnData)
}

// match html tag and replace it with ""
func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

// TibiaDataConvertEncodingtoISO88591 func - convert string from UTF-8 to latin1 (ISO 8859-1)
func TibiaDataConvertEncodingtoISO88591(data string) (string, error) {
	data, err := charmap.ISO8859_1.NewEncoder().String(data)
	return data, err
}

// TibiaDataConvertEncodingtoUTF8 func - convert string from latin1 (ISO 8859-1) to UTF-8
func TibiaDataConvertEncodingtoUTF8(data string) (string, error) {
	data, err := charmap.ISO8859_1.NewDecoder().String(data)
	return data, err
}

// isEnvExist func - check if environment var is set
func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsBool func - read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

/*
// getEnvAsFloat func - read an environment variable into a float64 or return default value
func getEnvAsFloat(name string, defaultVal float64) float64 {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

// getEnvAsInt func - read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
*/
