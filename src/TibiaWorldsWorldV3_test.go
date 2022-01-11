package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWintera(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/worlds/world/Wintera.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	worldJson := TibiaWorldsWorldV3Impl("Wintera", string(data))
	assert := assert.New(t)

	world := worldJson.Worlds.World

	assert.Equal("Wintera", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(281, world.PlayersOnline)
	assert.Equal(1023, world.RecordPlayers)
	assert.Equal("2020-05-04T01:25:30Z", world.RecordDate)
	assert.Equal("2018-04", world.CreationDate)
	assert.Equal("North America", world.Location)
	assert.Equal("Open PvP", world.PvpType)
	assert.False(world.PremiumOnly)
	assert.Empty(world.TransferType)
	assert.Equal(4, len(world.WorldsQuestTitles))
	assert.Equal("A Piece of Cake", world.WorldsQuestTitles[0])
	assert.Equal("Rise of Devovorga", world.WorldsQuestTitles[1])
	assert.Equal("Bewitched", world.WorldsQuestTitles[2])
	assert.Equal("The Colours of Magic", world.WorldsQuestTitles[3])
	assert.True(world.BattleyeProtected)
	assert.Equal("2018-04-19", world.BattleyeDate)
	assert.Equal("regular", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(281, len(world.OnlinePlayers))

	firstPlayer := world.OnlinePlayers[0]
	assert.Equal("Akiles Boy", firstPlayer.Name)
	assert.Equal(281, firstPlayer.Level)
	assert.Equal("Royal Paladin", firstPlayer.Vocation)
}
