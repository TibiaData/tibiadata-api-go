package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/TibiaData/tibiadata-api-go/src/validation"
	_ "github.com/mantyr/go-charset/data"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	// TibiaData app resty vars
	TibiaDataUserAgent, TibiaDataProxyDomain string

	// ErrorNotFound will be returned if the requests ends up in a 404
	ErrorNotFound = errors.New("page not found")
)

// DebugOutInformation wraps OutInformation with some debug info
type DebugOutInformation struct {
	Information Information `json:"information"`
	Debug       Debug       `json:"debug"`
}

// OutInformation wraps Information in other for all json outputs be consistent
type OutInformation struct {
	Information Information `json:"information"`
}

// Information stores some API related data
type Information struct {
	APIVersion int    `json:"api_version"`
	Timestamp  string `json:"timestamp"`
	Status     Status `json:"status"`
}

// Status stores information about the response
type Status struct {
	HTTPCode int    `json:"http_code"`
	Error    int    `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
}

// TibiaDataRequest is the struct of request information
type TibiaDataRequestStruct struct {
	Method   string            `json:"method"`    // Request method (default: GET)
	URL      string            `json:"url"`       // Request URL
	FormData map[string]string `json:"form_data"` // Request form content (used when POST)
}

// RunWebServer starts the gin server
// It blocks the code and will only finish execution on shutdown
func runWebServer() {
	// Setting gin-application to certain mode if GIN_MODE is set to release, test or debug (default is release)
	switch ginMode := getEnv("GIN_MODE", "release"); ginMode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Logging the gin.mode
	log.Printf("[info] TibiaData API gin-mode: %s", gin.Mode())

	// Starting an Engine instance
	router := gin.Default()

	// Gin middleware to enable GZIP support
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Set 404 not found page
	router.NoRoute(func(c *gin.Context) {
		TibiaDataErrorHandler(
			c,
			ErrorNotFound,
			http.StatusNotFound,
		)
	})

	// Disable proxy feature of gin
	_ = router.SetTrustedProxies(nil)

	// Set the ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		data := Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
				Message:  "pong",
			},
		}

		var output OutInformation
		output.Information = data

		c.JSON(http.StatusOK, output)
	})

	// Set the health endpoint
	router.GET("/health", func(c *gin.Context) {
		//TODO: Make this actually return problems
		data := Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
				Message:  "UP",
			},
		}

		var output OutInformation
		output.Information = data

		c.JSON(http.StatusOK, output)
	})

	// Set the debug endpoint
	router.GET("/debug", debugHandler)

	// TibiaData API version 3 endpoints
	v3 := router.Group("/v3")
	{
		// Tibia characters
		v3.GET("/character/:name", tibiaCharactersCharacterV3)

		// Tibia creatures
		v3.GET("/creature/:race", tibiaCreaturesCreatureV3)
		v3.GET("/creatures", tibiaCreaturesOverviewV3)

		// Tibia fansites
		v3.GET("/fansites", tibiaFansitesV3)

		// Tibia guilds
		v3.GET("/guild/:name", tibiaGuildsGuildV3)
		// v3.GET("/guild/:name/events",TibiaGuildsGuildEventsV3)
		// v3.GET("/guild/:name/wars",TibiaGuildsGuildWarsV3)
		v3.GET("/guilds/:world", tibiaGuildsOverviewV3)

		// Tibia highscores
		v3.GET("/highscores/:world", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/experience/"+TibiaDataDefaultVoc)
		})
		v3.GET("/highscores/:world/:category", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/"+c.Param("category")+"/"+TibiaDataDefaultVoc)
		})
		v3.GET("/highscores/:world/:category/:vocation", tibiaHighscoresV3)

		// Tibia houses
		v3.GET("/house/:world/:house_id", tibiaHousesHouseV3)
		v3.GET("/houses/:world/:town", tibiaHousesOverviewV3)

		// Tibia killstatistics
		v3.GET("/killstatistics/:world", tibiaKillstatisticsV3)

		// Tibia news
		v3.GET("/news/archive", tibiaNewslistV3)       // all categories (default 90 days)
		v3.GET("/news/archive/:days", tibiaNewslistV3) // all categories
		v3.GET("/news/id/:news_id", tibiaNewsV3)       // shows one news entry
		v3.GET("/news/latest", tibiaNewslistV3)        // only news and articles
		v3.GET("/news/newsticker", tibiaNewslistV3)    // only news_ticker

		// Tibia spells
		v3.GET("/spell/:spell_id", tibiaSpellsSpellV3)
		v3.GET("/spells", tibiaSpellsOverviewV3)

		// Tibia worlds
		v3.GET("/world/:name", tibiaWorldsWorldV3)
		v3.GET("/worlds", tibiaWorldsOverviewV3)
	}

	// Container version details endpoint
	router.GET("/versions", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"release": TibiaDataBuildRelease,
			"build":   TibiaDataBuildBuilder,
			"commit":  TibiaDataBuildCommit,
			"edition": TibiaDataBuildEdition,
		})
	})

	// Build the http server
	server := &http.Server{
		Addr:    ":8080", // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		Handler: router,
	}

	// Prepare for a graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Run a go routine that will receive the shutdown input
	go func() {
		<-quit
		log.Println("[info] TibiaData API received shutdown input")
		if err := server.Close(); err != nil {
			log.Fatal("[error] TibiaData API server close error:", err)
		}
	}()

	log.Println("[info] TibiaData API starting webserver")

	// Run the server
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("[info] TibiaData API server gracefully shut down")
		} else {
			log.Fatal("[error] TibiaData API server closed unexpectedly")
		}
	}
}

// Character godoc
// @Summary      Show one character
// @Description  Show all information about one character available
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        name path string true "The character name" extensions(x-example=Trollefar)
// @Success      200  {object}  CharacterResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/character/{name} [get]
func tibiaCharactersCharacterV3(c *gin.Context) {
	// Getting params from URL
	name := c.Param("name")

	// Validate the name
	err := validation.IsCharacterNameValid(name)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	// Build the request structure
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=characters&name=" + TibiaDataQueryEscapeStringV3(name),
	}

	// Handle the request
	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCharactersCharacterV3Impl(BoxContentHTML)
		},
		"TibiaCharactersCharacterV3")
}

// Creatures godoc
// @Summary      List of creatures
// @Description  Show all creatures listed
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Success      200  {object}  CreaturesOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/creatures [get]
func tibiaCreaturesOverviewV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCreaturesOverviewV3Impl(BoxContentHTML)
		},
		"TibiaCreaturesOverviewV3")
}

// Creature godoc
// @Summary      Show one creature
// @Description  Show all information about one creature
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Param        race path string true "The race of creature" extensions(x-example=nightmare)
// @Success      200  {object}  CreatureResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/creature/{race} [get]
func tibiaCreaturesCreatureV3(c *gin.Context) {
	// getting params from URL
	race := c.Param("race")

	// Validate the race
	endpoint, err := validation.IsCreatureNameValid(race)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures&race=" + endpoint,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCreaturesCreatureV3Impl(race, BoxContentHTML)
		},
		"TibiaCreaturesCreatureV3")
}

// Fansites godoc
// @Summary      Promoted and supported fansites
// @Description  List of all promoted and supported fansites
// @Tags         fansites
// @Accept       json
// @Produce      json
// @Success      200  {object}  FansitesResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/fansites [get]
func tibiaFansitesV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=fansites",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaFansitesV3Impl(BoxContentHTML)
		},
		"TibiaFansitesV3")
}

// Guild godoc
// @Summary      Show one guild
// @Description  Show all information about one guild
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of guild" extensions(x-example=Elysium)
// @Success      200  {object}  GuildResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/guild/{name} [get]
func tibiaGuildsGuildV3(c *gin.Context) {
	// getting params from URL
	guild := c.Param("name")

	// Validate the name
	err := validation.IsGuildNameValid(guild)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=" + TibiaDataQueryEscapeStringV3(guild),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaGuildsGuildV3Impl(guild, BoxContentHTML)
		},
		"TibiaGuildsGuildV3")
}

// Guilds godoc
// @Summary      List all guilds from a world
// @Description  Show all guilds on a certain world
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        world path string true "The world" extensions(x-example=Antica)
// @Success      200  {object}  GuildsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/guilds/{world} [get]
func tibiaGuildsOverviewV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaGuildsOverviewV3Impl(world, BoxContentHTML)
		},
		"TibiaGuildsOverviewV3")
}

// Highscores godoc
// @Summary      Highscores of tibia
// @Description  Show all highscores of tibia
// @Tags         highscores
// @Accept       json
// @Produce      json
// @Param        world    path string true "The world" default(all) extensions(x-example=Antica)
// @Param        category path string true "The category" default(experience) Enums(achievements, axefighting, charmpoints, clubfighting, distancefighting, experience, fishing, fistfighting, goshnarstaint, loyaltypoints, magiclevel, shielding, swordfighting, dromescore) extensions(x-example=fishing)
// @Param        vocation path string true "The vocation" default(all) Enums(all, knights, paladins, sorcerers, druids) extensions(x-example=knights)
// @Success      200  {object}  HighscoresResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/highscores/{world}/{category}/{vocation} [get]
func tibiaHighscoresV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	category := c.Param("category")
	vocation := c.Param("vocation")

	// Check if vocation is valid
	err := validation.IsVocationValid(vocation)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	if strings.EqualFold(world, "all") {
		world = ""
	} else {
		world = TibiaDataStringWorldFormatToTitleV3(world)
	}

	if world != "" {
		// Check if world exists
		exists, err := validation.WorldExists(world)
		if err != nil {
			TibiaDataErrorHandler(c, err, 0)
			return
		}

		if !exists {
			TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
			return
		}
	}

	if category != "" {
		err = validation.IsHighscoreCategoryValid(category)
		if err != nil {
			TibiaDataErrorHandler(c, validation.ErrorHighscoreCategoryDoesNotExist, http.StatusBadRequest)
			return
		}
	}

	highscoreCategory := validation.HighscoreCategoryFromString(category)

	// Sanitize of vocation input
	vocationName, vocationid := TibiaDataVocationValidator(vocation)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=highscores&world=" + TibiaDataQueryEscapeStringV3(world) + "&category=" + strconv.Itoa(int(highscoreCategory)) + "&profession=" + TibiaDataQueryEscapeStringV3(vocationid) + "&currentpage=400000000000000",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaHighscoresV3Impl(world, highscoreCategory, vocationName, BoxContentHTML)
		},
		"TibiaHighscoresV3")
}

// House godoc
// @Summary      House view
// @Description  Show all information about one house
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world     path string true "The world to show" extensions(x-example=Antica)
// @Param        house_id  path int    true "The ID of the house" extensions(x-example=35019)
// @Success      200  {object}  HouseResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/house/{world}/{house_id} [get]
func tibiaHousesHouseV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	houseidStr := c.Param("house_id")

	houseid, err := strconv.Atoi(houseidStr)
	if err != nil {
		TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// check if house exists
	exists, err = validation.HouseExistsRaw(houseid)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorHouseDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=houses&page=view&world=" + TibiaDataQueryEscapeStringV3(world) + "&houseid=" + TibiaDataQueryEscapeStringV3(houseidStr),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaHousesHouseV3Impl(houseid, BoxContentHTML)
		},
		"TibiaHousesHouseV3")
}

// Houses godoc
// @Summary      List of houses
// @Description  Show all houses filtered on world and town
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Param        town  path string true "The town to show" extensions(x-example=Venore)
// @Success      200  {object}  HousesOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/houses/{world}/{town} [get]
//TODO: This API needs to be refactored somehow to use tibiaDataRequestHandler
func tibiaHousesOverviewV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	town := c.Param("town")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)
	town = TibiaDataStringWorldFormatToTitleV3(town)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// Check if town exists
	exists, err = validation.TownExists(town)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorTownDoesNotExist, http.StatusBadRequest)
		return
	}

	jsonData, err := TibiaHousesOverviewV3Impl(c, world, town, TibiaDataHTMLDataCollectorV3)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// return jsonData
	TibiaDataAPIHandleResponse(c, "TibiaHousesOverviewV3", jsonData)
}

// Killstatistics godoc
// @Summary      The killstatistics
// @Description  Show all killstatistics filtered on world
// @Tags         killstatistics
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Success      200  {object}  KillStatisticsResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/killstatistics/{world} [get]
func tibiaKillstatisticsV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaKillstatisticsV3Impl(world, BoxContentHTML)
		},
		"TibiaKillstatisticsV3")
}

// News archive godoc
// @Summary      Show news archive (90 days)
// @Description  Show news archive with a filtering on 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/news/archive [get]
func tibiaNewslistArchiveV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// News archive (with day filter) godoc
// @Summary      Show news archive (with days filter)
// @Description  Show news archive with a filtering option on days
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        days path int true "The number of days to show" default(90) minimum(1) extensions(x-example=30)
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/news/archive/{days} [get]
func tibiaNewslistArchiveDaysV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// Latest news godoc
// @Summary      Show newslist (90 days)
// @Description  Show newslist with filtering on articles and news of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/news/latest [get]
func tibiaNewslistLatestV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// News ticker godoc
// @Summary      Show news tickers (90 days)
// @Description  Show news of type news tickers of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/news/newsticker [get]
func tibiaNewslistV3(c *gin.Context) {
	// getting params from URL
	daysStr := c.Param("days")

	var (
		days int
		err  error
	)

	if daysStr != "" {
		// convert param to int
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
			return
		}
	}

	if days == 0 {
		days = 90 // default for recent posts
	}

	// generating dates to pass to FormData
	DateBegin := time.Now().AddDate(0, 0, -days)
	DateEnd := time.Now()

	tibiadataRequest := TibiaDataRequestStruct{
		Method: http.MethodPost,
		URL:    "https://www.tibia.com/news/?subtopic=newsarchive",
		FormData: map[string]string{
			"filter_begin_day":   strconv.Itoa(DateBegin.UTC().Day()),        // period
			"filter_begin_month": strconv.Itoa(int(DateBegin.UTC().Month())), // period
			"filter_begin_year":  strconv.Itoa(DateBegin.UTC().Year()),       // period
			"filter_end_day":     strconv.Itoa(DateEnd.UTC().Day()),          // period
			"filter_end_month":   strconv.Itoa(int(DateEnd.UTC().Month())),   // period
			"filter_end_year":    strconv.Itoa(DateEnd.UTC().Year()),         // period
			"filter_cipsoft":     "cipsoft",                                  // category
			"filter_community":   "community",                                // category
			"filter_development": "development",                              // category
			"filter_support":     "support",                                  // category
			"filter_technical":   "technical",                                // category
		},
	}

	if c.Request != nil {
		// getting type of news list
		switch tmp := strings.Split(c.Request.URL.Path, "/"); tmp[3] {
		case "newsticker":
			tibiadataRequest.FormData["filter_ticker"] = "ticker"
		case "latest":
			tibiadataRequest.FormData["filter_article"] = "article"
			tibiadataRequest.FormData["filter_news"] = "news"
		case "archive":
			tibiadataRequest.FormData["filter_ticker"] = "ticker"
			tibiadataRequest.FormData["filter_article"] = "article"
			tibiadataRequest.FormData["filter_news"] = "news"
		}
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaNewslistV3Impl(days, BoxContentHTML)
		},
		"TibiaNewslistV3")
}

// News entry godoc
// @Summary      Show one news entry
// @Description  Show one news entry
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        news_id path int true "The ID of news entry" extensions(x-example=6512)
// @Success      200  {object}  NewsResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/news/id/{news_id} [get]
func tibiaNewsV3(c *gin.Context) {
	// getting params from URL
	newsIDStr := c.Param("news_id")

	// convert param to int
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
		return
	}

	// checking the NewsID provided
	err = validation.IsNewsIDValid(newsID)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/news/?subtopic=newsarchive&id=" + newsIDStr,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaNewsV3Impl(newsID, tibiadataRequest.URL, BoxContentHTML)
		},
		"TibiaNewsV3")
}

// Spells godoc
// @Summary      List all spells
// @Description  Show all spells
// @Tags         spells
// @Accept       json
// @Produce      json
// @Success      200  {object}  SpellsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/spells [get]
func tibiaSpellsOverviewV3(c *gin.Context) {
	// getting params from URL
	vocation := c.Param("vocation")
	if vocation == "" {
		vocation = TibiaDataDefaultVoc
	}

	err := validation.IsVocationValid(vocation)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// Sanitize of vocation input
	vocationName, _ := TibiaDataVocationValidator(vocation)
	if vocationName == "all" || vocationName == "none" {
		vocationName = ""
	} else {
		// removes the last letter (s) from the string (required for spells page)
		vocationName = strings.TrimSuffix(vocationName, "s")
		// setting string to first upper case
		vocationName = strings.Title(strings.ToLower(vocationName))
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiaDataQueryEscapeStringV3(vocationName),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaSpellsOverviewV3Impl(vocationName, BoxContentHTML)
		},
		"TibiaSpellsOverviewV3")
}

// Spell godoc
// @Summary      Show one spell
// @Description  Show all information about one spell
// @Tags         spells
// @Accept       json
// @Produce      json
// @Param        spell_id path string true "The name of spell" extensions(x-example=stronghaste)
// @Success      200  {object}  SpellInformationResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/spell/{spell_id} [get]
func tibiaSpellsSpellV3(c *gin.Context) {
	// getting params from URL
	spellRaw := c.Param("spell_id")

	spell, err := validation.IsSpellNameOrFormulaValid(spellRaw)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&spell=" + spell,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaSpellsSpellV3Impl(spell, BoxContentHTML)
		},
		"TibiaSpellsSpellV3")
}

// Worlds godoc
// @Summary      List of all worlds
// @Description  Show all worlds of Tibia
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Success      200  {object}  WorldsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/worlds [get]
func tibiaWorldsOverviewV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaWorldsOverviewV3Impl(BoxContentHTML)
		},
		"TibiaWorldsOverviewV3")
}

// World godoc
// @Summary      Show one world
// @Description  Show all information about one world
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of world" extensions(x-example=Antica)
// @Success      200  {object}  WorldResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v3/world/{name} [get]
func tibiaWorldsWorldV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("name")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaWorldsWorldV3Impl(world, BoxContentHTML)
		},
		"TibiaWorldsWorldV3")
}

func TibiaDataErrorHandler(c *gin.Context, err error, httpCode int) {
	if err == nil {
		panic(errors.New("TibiaDataErrorHandler called with nil err"))
	}

	info := Information{
		APIVersion: TibiaDataAPIversion,
		Timestamp:  TibiaDataDatetimeV3(""),
		Status: Status{
			HTTPCode: httpCode,
		},
	}

	switch t := err.(type) {
	case validation.Error:
		if httpCode == 0 {
			if t.Code() == 10 || t.Code() == 11 {
				httpCode = http.StatusInternalServerError
				info.Status.HTTPCode = httpCode
			} else {
				httpCode = http.StatusBadRequest
				info.Status.HTTPCode = httpCode
			}
		}

		info.Status.Error = t.Code()
		info.Status.Message = t.Error()
	case error:
		if httpCode == 0 {
			httpCode = http.StatusBadGateway
			info.Status.HTTPCode = httpCode
		}

		info.Status.Message = err.Error()

		log.Printf("[TibiaDataErrorHandler] HTTPCode: %d], Message: %s", info.Status.HTTPCode, info.Status.Message)
	}

	var output OutInformation
	output.Information = info

	c.JSON(httpCode, output)
}

func tibiaDataRequestHandler(c *gin.Context, tibiaDataRequest TibiaDataRequestStruct, requestHandler func(string) (interface{}, error), handlerName string) {
	BoxContentHTML, err := TibiaDataHTMLDataCollectorV3(tibiaDataRequest)
	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadGateway)
		return
	}

	jsonData, err := requestHandler(BoxContentHTML)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// return jsonData
	TibiaDataAPIHandleResponse(c, handlerName, jsonData)
}

// TibiaDataAPIHandleResponse func - handling of responses..
// This should NOT be invoked if an error occured
func TibiaDataAPIHandleResponse(c *gin.Context, s string, j interface{}) {
	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] " + s + " - (" + c.Request.RequestURI + ") returned data:")
		js, err := json.Marshal(j)
		log.Printf("[debug] %s\n", js)
		if err != nil {
			log.Printf("[debug] the above had an error: %s\n", err)
		}
	}

	if TibiaDataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response
	c.JSON(http.StatusOK, j)
}

// TibiadataUserAgentGenerator func - creates User-Agent for requests
func TibiaDataUserAgentGenerator(version int) string {
	// setting product name
	useragent := "TibiaData-API/v" + strconv.Itoa(version)

	// adding details in parenthesis
	useragentDetails := []string{
		"release/" + TibiaDataBuildRelease,
		"build/" + TibiaDataBuildBuilder,
		"commit/" + TibiaDataBuildCommit,
		"edition/" + TibiaDataBuildEdition,
		TibiaDataHost,
	}
	useragent += " (" + strings.Join(useragentDetails, "; ") + ")"

	return useragent
}

// TibiaDataHTMLDataCollectorV3 func
func TibiaDataHTMLDataCollectorV3(TibiaDataRequest TibiaDataRequestStruct) (string, error) {
	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiaDataDebug var
	if TibiaDataDebug {
		client.SetDebug(true)
		client.EnableTrace()
	}

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   TibiaDataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Replace domain with proxy if env TIBIADATA_PROXY set
	if TibiaDataProxyDomain != "" {
		TibiaDataRequest.URL = strings.ReplaceAll(TibiaDataRequest.URL, "https://www.tibia.com/", TibiaDataProxyDomain)
	}

	// defining values for request
	var (
		res        *resty.Response
		err        error
		LogMessage string
	)

	switch TibiaDataRequest.Method {
	case resty.MethodPost:
		res, err = client.R().
			SetFormData(TibiaDataRequest.FormData).
			Post(TibiaDataRequest.URL)
	default:
		res, err = client.R().Get(TibiaDataRequest.URL)
	}

	if TibiaDataDebug {
		// logging trace information for resty
		TibiaDataRequestTraceLogger(res, err)
	}

	if err != nil {
		log.Printf("[error] TibiaDataHTMLDataCollectorV3 (Status: %s, URL: %s) in resp1: %s", res.Status(), res.Request.URL, err)

		switch res.StatusCode() {
		case http.StatusForbidden:
			// throttled request
			LogMessage = "request throttled due to rate-limitation on tibia.com"
			log.Printf("[warning] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
			return "", err

		case http.StatusFound:
			// Check if page is in maintenance mode
			location, _ := res.RawResponse.Location()
			if location.Host == "maintenance.tibia.com" {
				LogMessage := "maintenance mode detected on tibia.com"
				log.Printf("[info] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
				return "", err
			}
			fallthrough

		default:
			LogMessage = "unknown error occurred on tibia.com"
			log.Printf("[error] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
			return "", err
		}
	}

	// Convert body to io.Reader
	resIo := bytes.NewReader(res.Body())

	// wrap reader in a converting reader from ISO 8859-1 to UTF-8
	resIo2 := TibiaDataConvertEncodingtoUTF8(resIo)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resIo2)
	if err != nil {
		log.Printf("[error] TibiaDataHTMLDataCollectorV3 (URL: %s) error: %s", res.Request.URL, err)
	}

	// Find of this to get div with class BoxContent
	data, err := doc.Find(".Border_2 .Border_3").Html()
	if err != nil {
		return "", err
	}

	// Return of extracted html to functions..
	return data, nil
}
