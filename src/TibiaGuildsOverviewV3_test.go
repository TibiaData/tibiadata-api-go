package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestPremia(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/Premia.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	premiaGuildsJson, err := TibiaGuildsOverviewV3Impl("Premia", string(data))
	if err != nil {
		t.Fatal(err)
	}

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
	assert.Equal("", secondGuildInFormation.Description)
}
