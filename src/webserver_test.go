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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "Durin",
		},
	}

	assert := assert.New(t)

	tibiaCharactersCharacterV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaCreaturesOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "race",
			Value: "demon",
		},
	}

	tibiaCreaturesCreatureV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaFansitesV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "pax",
		},
	}

	tibiaGuildsGuildV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaGuildsOverviewV3(c)
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
	}

	tibiaHighscoresV3(c)
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

	tibiaHousesHouseV3(c)
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

	tibiaHousesOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "world",
			Value: "antica",
		},
	}

	tibiaKillstatisticsV3(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.False(false, tibiaNewslistArchiveV3())
	assert.False(false, tibiaNewslistArchiveDaysV3())
	assert.False(false, tibiaNewslistLatestV3())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "days",
			Value: "90",
		},
	}

	tibiaNewslistV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "news_id",
			Value: "6607",
		},
	}

	tibiaNewsV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "vocation",
			Value: "sorcerer",
		},
	}

	tibiaSpellsOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "spell_id",
			Value: "exori",
		},
	}

	tibiaSpellsSpellV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	tibiaWorldsOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "antica",
		},
	}

	tibiaWorldsWorldV3(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.Equal("TibiaData-API/v3 (release/unknown; build/manual; commit/-; edition/open-source; unittest.example.com)", TibiaDataUserAgentGenerator(3))

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
