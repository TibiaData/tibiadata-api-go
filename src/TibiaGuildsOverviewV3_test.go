package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPremia(t *testing.T) {
	data, err := os.ReadFile("../testdata/guilds/Premia.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	premiaGuildsJson := TibiaGuildsOverviewV3Impl("Premia", string(data))
	assert := assert.New(t)

	assert.Equal("Premia", premiaGuildsJson.Guilds.World)
	assert.Equal(38, len(premiaGuildsJson.Guilds.Active))
	assert.Equal(3, len(premiaGuildsJson.Guilds.Formation))

	orderOfGloryActiveGuild := premiaGuildsJson.Guilds.Active[28]
	assert.Equal("Order of Glory", orderOfGloryActiveGuild.Name)
	assert.Equal("https://static.tibia.com/images/guildlogos/Order_of_Glory.gif", orderOfGloryActiveGuild.LogoURL)
	assert.Equal("We are an English speaking guild of friends and allies from around the world who seek only peaceful questing, exploring, team hunts and a chill place to hang out. Message any of our leaders for an invitation. Contact Zyb with any problems.", orderOfGloryActiveGuild.Description)

	secondGuildInFormation := premiaGuildsJson.Guilds.Formation[1]
	assert.Equal("Konungen", secondGuildInFormation.Name)
	assert.Equal("https://static.tibia.com/images/community/default_logo.gif", secondGuildInFormation.LogoURL)
	assert.Empty(secondGuildInFormation.Description)
}
