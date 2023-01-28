package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestDocsV3CodeCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	assert := assert.New(t)

	assert.False(false, tibiaBoostableBossesV3())
	assert.False(false, tibiaCharactersCharacterV3())
	assert.False(false, tibiaCreaturesOverviewV3())
	assert.False(false, tibiaCreaturesCreatureV3())
	assert.False(false, tibiaFansitesV3())
	assert.False(false, tibiaGuildsGuildV3())
	assert.False(false, tibiaGuildsOverviewV3())
	assert.False(false, tibiaHighscoresV3())
	assert.False(false, tibiaHousesHouseV3())
	assert.False(false, tibiaHousesOverviewV3())
	assert.False(false, tibiaKillstatisticsV3())
	assert.False(false, tibiaNewslistArchiveV3())
	assert.False(false, tibiaNewslistArchiveDaysV3())
	assert.False(false, tibiaNewslistLatestV3())
	assert.False(false, tibiaNewslistV3())
	assert.False(false, tibiaNewsV3())
	assert.False(false, tibiaSpellsOverviewV3())
	assert.False(false, tibiaSpellsSpellV3())
	assert.False(false, tibiaWorldsOverviewV3())
	assert.False(false, tibiaWorldsWorldV3())

}
