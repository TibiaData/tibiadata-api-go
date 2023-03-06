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

	_ "github.com/mantyr/go-charset/data"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	// TibiaData app resty vars
	TibiaDataUserAgent, TibiaDataProxyDomain string
)

// Information - child of JSONData
type Information struct {
	APIVersion int    `json:"api_version"`
	Timestamp  string `json:"timestamp"`
}

// TibiaDataRequest - struct of request information
type TibiaDataRequestStruct struct {
	Method   string            `json:"method"`    // Request method (default: GET)
	URL      string            `json:"url"`       // Request URL
	FormData map[string]string `json:"form_data"` // Request form content (used when POST)
}

// runWebServer starts the gin server
func runWebServer() {
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

	// gin middleware to enable GZIP support
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// set 404 not found page
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// Set proxy feature of gin
	trustedProxies := getEnv("GIN_TRUSTED_PROXIES", "")
	if len(trustedProxies) > 0 {
		_ = router.SetTrustedProxies(strings.Split(trustedProxies, ","))
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	// Ping-endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// health endpoints for kubernetes
	router.GET("/health", healthz)
	router.GET("/healthz", healthz)
	router.GET("/readyz", readyz)

	// TibiaData API version 3
	v3 := router.Group("/v3")
	{
		// Tibia characters
		v3.GET("/boostablebosses", tibiaBoostableBossesV3)

		// Tibia characters
		v3.GET("/character/:name", tibiaCharactersCharacterV3)

		// Tibia creatures
		v3.GET("/creature/:race", tibiaCreaturesCreatureV3)
		v3.GET("/creatures", tibiaCreaturesOverviewV3)

		// Tibia fansites
		v3.GET("/fansites", tibiaFansitesV3)

		// Tibia guilds
		v3.GET("/guild/:name", tibiaGuildsGuildV3)
		//v3.GET("/guild/:name/events",TibiaGuildsGuildEventsV3)
		//v3.GET("/guild/:name/wars",TibiaGuildsGuildWarsV3)
		v3.GET("/guilds/:world", tibiaGuildsOverviewV3)

		// Tibia highscores
		v3.GET("/highscores/:world", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/experience/"+TibiaDataDefaultVoc+"/1")
		})
		v3.GET("/highscores/:world/:category", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v3.BasePath()+"/highscores/"+c.Param("world")+"/"+c.Param("category")+"/"+TibiaDataDefaultVoc+"/1")
		})
		v3.GET("/highscores/:world/:category/:vocation", tibiaHighscoresV3)
		v3.GET("/highscores/:world/:category/:vocation/:page", tibiaHighscoresV3)

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

	// container version details endpoint
	router.GET("/versions", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"release": TibiaDataBuildRelease,
			"build":   TibiaDataBuildBuilder,
			"commit":  TibiaDataBuildCommit,
			"edition": TibiaDataBuildEdition,
		})
	})

	// build the http server
	server := &http.Server{
		Addr:    ":8080", // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		Handler: router,
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// we run a go routine that will receive the shutdown input
	go func() {
		<-quit
		log.Println("[info] TibiaData API received shutdown input")
		if err := server.Close(); err != nil {
			log.Fatal("[error] TibiaData API server close error:", err)
		}
	}()

	// run the server
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("[info] TibiaData API server gracefully shut down")
		} else {
			log.Fatal("[error] TibiaData API server closed unexpectedly")
		}
	}
}

// BoostableBosses godoc
// @Summary      List of boostable bosses
// @Description  Show all boostable bosses listed
// @Tags         boostable bosses
// @Accept       json
// @Produce      json
// @Success      200  {object}  BoostableBossesOverviewResponse
// @Router       /v3/boostablebosses [get]
func tibiaBoostableBossesV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=boostablebosses",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaBoostableBossesOverviewV3Impl(BoxContentHTML), http.StatusOK
		},
		"TibiaBoostableBossesV3")
}

// Character godoc
// @Summary      Show one character
// @Description  Show all information about one character available
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        name path string true "The character name" extensions(x-example=Trollefar)
// @Success      200  {object}  CharacterResponse
// @Router       /v3/character/{name} [get]
func tibiaCharactersCharacterV3(c *gin.Context) {
	// getting params from URL
	name := c.Param("name")

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=characters&name=" + TibiaDataQueryEscapeStringV3(name),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaCharactersCharacterV3Impl(BoxContentHTML), http.StatusOK
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
// @Router       /v3/creatures [get]
func tibiaCreaturesOverviewV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaCreaturesOverviewV3Impl(BoxContentHTML), http.StatusOK
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
// @Router       /v3/creature/{race} [get]
func tibiaCreaturesCreatureV3(c *gin.Context) {
	// getting params from URL
	race := c.Param("race")

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures&race=" + TibiaDataQueryEscapeStringV3(race),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaCreaturesCreatureV3Impl(race, BoxContentHTML), http.StatusOK
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
// @Router       /v3/fansites [get]
func tibiaFansitesV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=fansites",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaFansitesV3Impl(BoxContentHTML), http.StatusOK
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
// @Router       /v3/guild/{name} [get]
func tibiaGuildsGuildV3(c *gin.Context) {
	// getting params from URL
	guild := c.Param("name")

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=" + TibiaDataQueryEscapeStringV3(guild),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaGuildsGuildV3Impl(guild, BoxContentHTML), http.StatusOK
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
// @Router       /v3/guilds/{world} [get]
func tibiaGuildsOverviewV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaGuildsOverviewV3Impl(world, BoxContentHTML), http.StatusOK
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
// @Param        category path string true "The category" default(experience) Enums(achievements, axefighting, charmpoints, clubfighting, distancefighting, experience, fishing, fistfighting, goshnarstaint, loyaltypoints, magiclevel, shielding, swordfighting, dromescore, bosspoints) extensions(x-example=fishing)
// @Param        vocation path string true "The vocation" default(all) Enums(all, knights, paladins, sorcerers, druids) extensions(x-example=knights)
// @Param        page     path int    true "The current page" default(1) minimum(1) maximum(20) extensions(x-example=1)
// @Success      200  {object}  HighscoresResponse
// @Router       /v3/highscores/{world}/{category}/{vocation}/{page} [get]
func tibiaHighscoresV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	category := c.Param("category")
	vocation := c.Param("vocation")
	page := c.Param("page")

	// maybe return error on faulty vocation value?!

	// Adding fix for First letter to be upper and rest lower
	if strings.EqualFold(world, "all") {
		world = ""
	} else {
		world = TibiaDataStringWorldFormatToTitleV3(world)
	}

	highscoreCategory := HighscoreCategoryFromString(category)

	// Sanitize of vocation input
	vocationName, vocationid := TibiaDataVocationValidator(vocation)

	// checking the page provided
	if page == "" {
		page = "1"
	}
	if TibiaDataStringToIntegerV3(page) < 1 || TibiaDataStringToIntegerV3(page) > 23 {
		TibiaDataAPIHandleResponse(c, http.StatusBadRequest, "TibiaHighscoresV3", gin.H{"error": "page needs to be from 1 to 20 (possible until 23)"})
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=highscores&world=" + TibiaDataQueryEscapeStringV3(world) + "&category=" + strconv.Itoa(int(highscoreCategory)) + "&profession=" + TibiaDataQueryEscapeStringV3(vocationid) + "&currentpage=" + TibiaDataQueryEscapeStringV3(page),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaHighscoresV3Impl(world, highscoreCategory, vocationName, TibiaDataStringToIntegerV3(page), BoxContentHTML), http.StatusOK
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
// @Router       /v3/house/{world}/{house_id} [get]
func tibiaHousesHouseV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	houseid := c.Param("house_id")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=houses&page=view&world=" + TibiaDataQueryEscapeStringV3(world) + "&houseid=" + TibiaDataQueryEscapeStringV3(houseid),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaHousesHouseV3Impl(houseid, BoxContentHTML), http.StatusOK
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
// @Router       /v3/houses/{world}/{town} [get]
// TODO: This API needs to be refactored somehow to use tibiaDataRequestHandler
func tibiaHousesOverviewV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	town := c.Param("town")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)
	town = TibiaDataStringWorldFormatToTitleV3(town)

	jsonData := TibiaHousesOverviewV3Impl(c, world, town, TibiaDataHTMLDataCollectorV3)

	// return jsonData
	TibiaDataAPIHandleResponse(c, http.StatusOK, "TibiaHousesOverviewV3", jsonData)
}

// Killstatistics godoc
// @Summary      The killstatistics
// @Description  Show all killstatistics filtered on world
// @Tags         killstatistics
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Success      200  {object}  KillStatisticsResponse
// @Router       /v3/killstatistics/{world} [get]
func tibiaKillstatisticsV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaKillstatisticsV3Impl(world, BoxContentHTML), http.StatusOK
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
// @Router       /v3/news/newsticker [get]
func tibiaNewslistV3(c *gin.Context) {
	// getting params from URL
	days := TibiaDataStringToIntegerV3(c.Param("days"))
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
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaNewslistV3Impl(days, BoxContentHTML), http.StatusOK
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
// @Router       /v3/news/id/{news_id} [get]
func tibiaNewsV3(c *gin.Context) {
	// getting params from URL
	NewsID := TibiaDataStringToIntegerV3(c.Param("news_id"))

	// checking the NewsID provided
	if NewsID <= 0 {
		TibiaDataAPIHandleResponse(c, http.StatusBadRequest, "TibiaNewsV3", gin.H{"error": "no valid news_id provided"})
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/news/?subtopic=newsarchive&id=" + strconv.Itoa(NewsID),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaNewsV3Impl(NewsID, tibiadataRequest.URL, BoxContentHTML), http.StatusOK
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
// @Router       /v3/spells [get]
func tibiaSpellsOverviewV3(c *gin.Context) {
	// getting params from URL
	vocation := c.Param("vocation")
	if vocation == "" {
		vocation = TibiaDataDefaultVoc
	}

	// Sanitize of vocation input
	vocationName, _ := TibiaDataVocationValidator(vocation)
	if vocationName == "all" || vocationName == "none" {
		vocationName = ""
	} else {
		// removes the last letter (s) from the string (required for spells page)
		vocationName = strings.TrimSuffix(vocationName, "s")
		// setting string to first upper case
		vocationName = strings.Title(vocationName)
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiaDataQueryEscapeStringV3(vocationName),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaSpellsOverviewV3Impl(vocationName, BoxContentHTML), http.StatusOK
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
// @Router       /v3/spell/{spell_id} [get]
func tibiaSpellsSpellV3(c *gin.Context) {
	// getting params from URL
	spell := c.Param("spell_id")

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&spell=" + TibiaDataQueryEscapeStringV3(spell),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaSpellsSpellV3Impl(spell, BoxContentHTML), http.StatusOK
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
// @Router       /v3/worlds [get]
func tibiaWorldsOverviewV3(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaWorldsOverviewV3Impl(BoxContentHTML), http.StatusOK
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
// @Router       /v3/world/{name} [get]
func tibiaWorldsWorldV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("name")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitleV3(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds&world=" + TibiaDataQueryEscapeStringV3(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, int) {
			return TibiaWorldsWorldV3Impl(world, BoxContentHTML), http.StatusOK
		},
		"TibiaWorldsWorldV3")
}

func tibiaDataRequestHandler(c *gin.Context, tibiaDataRequest TibiaDataRequestStruct, requestHandler func(string) (interface{}, int), handlerName string) {
	BoxContentHTML, err := TibiaDataHTMLDataCollectorV3(tibiaDataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleResponse(c, http.StatusBadGateway, handlerName, gin.H{"error": err.Error()})
	}

	jsonData, httpStatusCode := requestHandler(BoxContentHTML)

	// return jsonData
	TibiaDataAPIHandleResponse(c, httpStatusCode, handlerName, jsonData)
}

// TibiaDataAPIHandleResponse func - handling of responses..
func TibiaDataAPIHandleResponse(c *gin.Context, httpCode int, s string, j interface{}) {
	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] " + s + " - (" + c.Request.RequestURI + ") returned data:")
		js, _ := json.Marshal(j)
		log.Printf("[debug] %s\n", js)
	}

	if TibiaDataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response
	c.JSON(httpCode, j)
}

// TibiaDataUserAgentGenerator func - creates User-Agent for requests
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
		log.Printf("[error] TibiaDataHTMLDataCollectorV3 (Status: %s, URL: %s) in resp1: %s", res.Status(), TibiaDataRequest.URL, err)

		switch res.StatusCode() {
		case http.StatusForbidden:
			// throttled request
			LogMessage = "request throttled due to rate-limitation on tibia.com"
			log.Printf("[warning] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
			return "", errors.New(LogMessage)

		case http.StatusFound:
			// Check if page is in maintenance mode
			location, _ := res.RawResponse.Location()
			if location.Host == "maintenance.tibia.com" {
				LogMessage := "maintenance mode detected on tibia.com"
				log.Printf("[info] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
				return "", errors.New(LogMessage)
			}
			fallthrough

		default:
			LogMessage = "unknown error occurred on tibia.com"
			log.Printf("[error] TibiaDataHTMLDataCollectorV3: %s!", LogMessage)
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
		log.Printf("[error] TibiaDataHTMLDataCollectorV3 (URL: %s) error: %s", TibiaDataRequest.URL, err)
	}

	// Find of this to get div with class BoxContent
	data, err := doc.Find(".Border_2 .Border_3").Html()
	if err != nil {
		log.Fatal(err)
	}

	// Return of extracted html to functions..
	return data, nil
}

// TibiaDataRequestTraceLogger func - prints out trace information to log
func TibiaDataRequestTraceLogger(res *resty.Response, err error) {
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

// healthz is a k8s liveness probe
func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusText(http.StatusOK)})
}

// readyz is a k8s readiness probe
func readyz(c *gin.Context) {
	if isReady == nil || !isReady.Load().(bool) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": http.StatusText(http.StatusServiceUnavailable)})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"status": http.StatusText(http.StatusOK)})
	TibiaDataAPIHandleResponse(c, http.StatusOK, "readyz", gin.H{"status": http.StatusText(http.StatusOK)})
}
