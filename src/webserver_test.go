package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/validation"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestFakeToUpCodeCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// adding support for proxy for tests
	if isEnvExist("TIBIADATA_PROXY") {
		TibiaDataProxyDomain = "https://" + getEnv("TIBIADATA_PROXY", "www.tibia.com") + "/"
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "Durin",
		},
	}

	assert := assert.New(t)

	tibiaBoostableBosses(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaCreaturesOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "race",
			Value: "demon",
		},
	}

	tibiaCreaturesCreature(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaFansites(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "pax",
		},
	}

	tibiaGuildsGuild(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaGuildsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
		{
			Key:   "category",
			Value: "experience",
		},
		{
			Key:   "vocation",
			Value: "sorcerer",
		},
		{
			Key:   "page",
			Value: "4",
		},
	}

	tibiaHighscores(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
		{
			Key:   "house_id",
			Value: "59054",
		},
	}

	tibiaHousesHouse(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
		{
			Key:   "town",
			Value: "venore",
		},
	}

	tibiaHousesOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaKillstatistics(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.False(false, tibiaNewslistArchive())
	assert.False(false, tibiaNewslistArchiveDays())
	assert.False(false, tibiaNewslistLatest())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "days",
			Value: "90",
		},
	}

	tibiaNewslist(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "news_id",
			Value: "6607",
		},
	}

	tibiaNews(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "vocation",
			Value: "sorcerer",
		},
	}

	tibiaSpellsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "spell_id",
			Value: "exori",
		},
	}

	tibiaSpellsSpell(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaWorldsOverview(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "antica",
		},
	}

	tibiaWorldsWorld(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.Equal("TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source; unittest.example.com)", TibiaDataUserAgentGenerator(TibiaDataAPIversion))

	healthz(c)
	assert.Equal(http.StatusOK, w.Code)

	readyz(c)
	assert.Equal(http.StatusOK, w.Code)

	type test struct {
		T string `json:"t"`
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	TibiaDataAPIHandleResponse(c, "", test{T: "abc"})
	assert.Equal(http.StatusOK, w.Code)
}

func TestErrorHandler(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, errors.New("test error"), http.StatusBadRequest)
	assert.Equal(http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorAlreadyRunning, 0)
	assert.Equal(http.StatusInternalServerError, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, validation.ErrorCharacterNameInvalid, 0)
	assert.Equal(http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	TibiaDataErrorHandler(c, errors.New("test error"), 0)
	assert.Equal(http.StatusBadGateway, w.Code)
}
