package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mantyr/go-charset/data"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	// Tibiadata app resty vars
	TibiadataUserAgent, TibiadataProxyDomain string
	TibiadataRequest                         = TibiadataRequestStruct{
		Method:   resty.MethodGet,
		URL:      "",
		FormData: make(map[string]string),
	}

	TibiadataHost string // set through env TIBIADATA_HOST
)

// TibiadataRequest - struct of request information
type TibiadataRequestStruct struct {
	Method   string            `json:"method"`    // Request method (default: GET)
	URL      string            `json:"url"`       // Request URL
	FormData map[string]string `json:"form_data"` // Request form content (used when POST)
}

func startWebServer() {
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

	router := gin.Default()

	// disable proxy feature of gin
	_ = router.SetTrustedProxies(nil)

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

		// Tibia guilds
		v3.GET("/guilds/guild/:guild", TibiaGuildsGuildV3)
		//v3.GET("/guilds/guild/:guild/events",TibiaGuildsGuildEventsV3)
		//v3.GET("/guilds/guild/:guild/wars",TibiaGuildsGuildWarsV3)
		v3.GET("/guilds/world/:world", TibiaGuildsOverviewV3)

		// Tibia highscores
		v3.GET("/highscores/world/:world", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/world/"+c.Param("world")+"/experience/"+TibiadataDefaultVoc)
		})
		v3.GET("/highscores/world/:world/:category", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/world/"+c.Param("world")+"/"+c.Param("category")+"/"+TibiadataDefaultVoc)
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
	_ = router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
*/

// TibiaDataAPIHandleOtherResponse func - handling of responses..
func TibiaDataAPIHandleOtherResponse(c *gin.Context, httpCode int, s string, j interface{}) {
	if TibiadataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response (with specific status code)
	c.JSON(httpCode, j)
}

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
func TibiadataHTMLDataCollectorV3(TibiadataRequest TibiadataRequestStruct) (string, error) {
	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiadataDebug var
	if TibiadataDebug {
		client.SetDebug(true)
		client.EnableTrace()
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

	// Replace domain with proxy if env TIBIADATA_PROXY set
	if TibiadataProxyDomain != "" {
		TibiadataRequest.URL = strings.ReplaceAll(TibiadataRequest.URL, "https://www.tibia.com/", TibiadataProxyDomain)
	}

	// defining values for request
	var res *resty.Response
	var err error

	switch TibiadataRequest.Method {
	case resty.MethodPost:
		res, err = client.R().
			SetFormData(TibiadataRequest.FormData).
			Post(TibiadataRequest.URL)
	default:
		res, err = client.R().Get(TibiadataRequest.URL)

	}

	if TibiadataDebug {
		// logging trace information for resty
		TibiadataRequestTraceLogger(res, err)
	}

	if err != nil {
		log.Printf("[error] TibiadataHTMLDataCollectorV3 (Status: %s, URL: %s) in resp1: %s", res.Status(), TibiadataRequest.URL, err)

		var LogMessage string
		switch res.StatusCode() {
		case http.StatusForbidden:
			// throttled request
			LogMessage = "request throttled due to rate-limitation on tibia.com"
			log.Printf("[warning] TibiadataHTMLDataCollectorV3: %s!", LogMessage)
			return "", errors.New(LogMessage)

		case http.StatusFound:
			// Check if page is in maintenance mode
			location, _ := res.RawResponse.Location()
			if location.Host == "maintenance.tibia.com" {
				LogMessage := "maintenance mode detected on tibia.com"
				log.Printf("[info] TibiadataHTMLDataCollectorV3: %s!", LogMessage)
				return "", errors.New(LogMessage)
			}
			fallthrough

		default:
			LogMessage = "unknown error occurred on tibia.com"
			log.Printf("[error] TibiadataHTMLDataCollectorV3: %s!", LogMessage)
			return "", errors.New(LogMessage)
		}
	}

	// Convert body to io.Reader
	resIo := bytes.NewReader(res.Body())

	// wrap reader in a converting reader from ISO 8859-1 to UTF-8
	resIo2 := TibiaDataConvertEncodingtoUTF8(resIo)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resIo2)
	if err != nil {
		log.Printf("[error] TibiadataHTMLDataCollectorV3 (URL: %s) error: %s", TibiadataRequest.URL, err)
	}

	// Find of this to get div with class BoxContent
	data, err := doc.Find(".Border_2 .Border_3").Html()
	if err != nil {
		log.Fatal(err)
	}

	// Return of extracted html to functions..
	return data, nil
}

// TibiadataRequestTraceLogger func - prints out trace information to log
func TibiadataRequestTraceLogger(res *resty.Response, err error) {
	log.Println("TRACE RESTY",
		"\n~~~ TRACE INFO ~~~",
		"\nDNSLookup      :", res.Request.TraceInfo().DNSLookup,
		"\nConnTime       :", res.Request.TraceInfo().ConnTime,
		"\nTCPConnTime    :", res.Request.TraceInfo().TCPConnTime,
		"\nTLSHandshake   :", res.Request.TraceInfo().TLSHandshake,
		"\nServerTime     :", res.Request.TraceInfo().ServerTime,
		"\nResponseTime   :", res.Request.TraceInfo().ResponseTime,
		"\nTotalTime      :", res.Request.TraceInfo().TotalTime,
		"\nIsConnReused   :", res.Request.TraceInfo().IsConnReused,
		"\nIsConnWasIdle  :", res.Request.TraceInfo().IsConnWasIdle,
		"\nConnIdleTime   :", res.Request.TraceInfo().ConnIdleTime,
		"\nRequestAttempt :", res.Request.TraceInfo().RequestAttempt,
		"\nRemoteAddr     :", res.Request.TraceInfo().RemoteAddr.String(),
		"\n==============================================================================")
}
