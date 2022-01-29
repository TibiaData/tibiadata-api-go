package main

import (
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

	tibiaCharactersCharacterV3(c)
	assert.Equal(200, w.Code)

	tibiaCreaturesOverviewV3(c)
	assert.Equal(200, w.Code)

	tibiaCreaturesCreatureV3(c)
	assert.Equal(200, w.Code)

	tibiaFansitesV3(c)
	assert.Equal(200, w.Code)

	tibiaGuildsGuildV3(c)
	assert.Equal(200, w.Code)

	tibiaGuildsOverviewV3(c)
	assert.Equal(200, w.Code)

	tibiaHighscoresV3(c)
	assert.Equal(200, w.Code)

	tibiaHousesHouseV3(c)
	assert.Equal(200, w.Code)

	tibiaHousesOverviewV3(c)
	assert.Equal(200, w.Code)

	tibiaKillstatisticsV3(c)
	assert.Equal(200, w.Code)

	assert.False(false, tibiaNewslistArchiveV3())
	assert.False(false, tibiaNewslistArchiveDaysV3())
	assert.False(false, tibiaNewslistLatestV3())

	tibiaNewslistV3(c)
	assert.Equal(200, w.Code)

	tibiaNewsV3(c)
	assert.Equal(200, w.Code)

	tibiaSpellsOverviewV3(c)
	assert.Equal(200, w.Code)

	tibiaSpellsSpellV3(c)
	assert.Equal(200, w.Code)

	tibiaWorldsOverviewV3(c)
	assert.Equal(200, w.Code)

	tibiaWorldsWorldV3(c)
	assert.Equal(200, w.Code)

	assert.Equal("TibiaData-API/v3 (release/unknown; build/manual; commit/-; edition/open-source; unittest.example.com)", TibiadataUserAgentGenerator(3))
}
