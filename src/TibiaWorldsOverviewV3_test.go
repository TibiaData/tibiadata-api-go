package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorlds(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/worlds/worlds.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	worldsJson := TibiaWorldsOverviewV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal(8756, worldsJson.Worlds.PlayersOnline)
	assert.Equal(64028, worldsJson.Worlds.RecordPlayers)
	assert.Equal("2007-11-28T18:26:00Z", worldsJson.Worlds.RecordDate)
	assert.Equal(76, len(worldsJson.Worlds.RegularWorlds))
	assert.Equal(6, len(worldsJson.Worlds.TournamentWorlds))

	adra := worldsJson.Worlds.RegularWorlds[0]
	assert.Equal("Adra", adra.Name)
	assert.Equal("online", adra.Status)
	assert.Equal(18, adra.PlayersOnline)
	assert.Equal("Europe", adra.Location)
	assert.Equal("Open PvP", adra.PvpType)
	assert.Equal(false, adra.PremiumOnly)
	assert.Equal("blocked", adra.TransferType)
	assert.Equal(true, adra.BattleyeProtected)
	assert.Equal("release", adra.BattleyeDate)
	assert.Equal("regular", adra.GameWorldType)
	assert.Equal("", adra.TournamentWorldType)

	astera := worldsJson.Worlds.RegularWorlds[4]
	assert.Equal("Astera", astera.Name)
	assert.Equal("online", astera.Status)
	assert.Equal(222, astera.PlayersOnline)
	assert.Equal("North America", astera.Location)
	assert.Equal("Optional PvP", astera.PvpType)
	assert.Equal(false, astera.PremiumOnly)
	assert.Equal("regular", astera.TransferType)
	assert.Equal(true, astera.BattleyeProtected)
	assert.Equal("2017-09-12", astera.BattleyeDate)
	assert.Equal("regular", astera.GameWorldType)
	assert.Equal("", astera.TournamentWorldType)

	endera := worldsJson.Worlds.TournamentWorlds[1]
	assert.Equal("Endera", endera.Name)
	assert.Equal("unknown", endera.Status)
	assert.Equal(0, endera.PlayersOnline)
	assert.Equal("North America", endera.Location)
	assert.Equal("Optional PvP", endera.PvpType)
	assert.Equal(true, endera.PremiumOnly)
	assert.Equal("blocked", endera.TransferType)
	assert.Equal(true, endera.BattleyeProtected)
	assert.Equal("release", endera.BattleyeDate)
	assert.Equal("tournament", endera.GameWorldType)
	assert.Equal("restricted", endera.TournamentWorldType)
}
