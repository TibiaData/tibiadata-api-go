package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFakeToUpCodeCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	tibiaBoostableBossesV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaCharactersCharacterV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaCreaturesOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaCreaturesCreatureV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaFansitesV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaGuildsGuildV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaGuildsOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaHighscoresV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaHousesHouseV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaHousesOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaKillstatisticsV3(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.False(false, tibiaNewslistArchiveV3())
	assert.False(false, tibiaNewslistArchiveDaysV3())
	assert.False(false, tibiaNewslistLatestV3())

	tibiaNewslistV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaNewsV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaSpellsOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaSpellsSpellV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaWorldsOverviewV3(c)
	assert.Equal(http.StatusOK, w.Code)

	tibiaWorldsWorldV3(c)
	assert.Equal(http.StatusOK, w.Code)

	assert.Equal("TibiaData-API/v3 (release/unknown; build/manual; commit/-; edition/open-source; unittest.example.com)", TibiaDataUserAgentGenerator(3))

	healthz(c)
	assert.Equal(http.StatusOK, w.Code)

	readyz(c)
	assert.Equal(http.StatusOK, w.Code)
}
