package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestOrderofGlory(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Order of Glory.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	orderOfGloryJson, err := TibiaGuildsGuildV3Impl("Order of Glory", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Order of Glory", orderOfGloryJson.Guilds.Guild.Name)
	assert.Equal("Premia", orderOfGloryJson.Guilds.Guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Order_of_Glory.gif", orderOfGloryJson.Guilds.Guild.LogoURL)
	assert.Equal("We are an English speaking guild of friends and allies from around the world who seek only peaceful questing, exploring, team hunts and a chill place to hang out. Message any of our leaders for an invitation. Contact Zyb with any problems.", orderOfGloryJson.Guilds.Guild.Description)
	assert.Nil(orderOfGloryJson.Guilds.Guild.Guildhalls)
	assert.True(orderOfGloryJson.Guilds.Guild.Active)
	assert.Equal("2020-06-27", orderOfGloryJson.Guilds.Guild.Founded)
	assert.True(orderOfGloryJson.Guilds.Guild.Applications)
	assert.Equal("", orderOfGloryJson.Guilds.Guild.Homepage)
	assert.False(orderOfGloryJson.Guilds.Guild.InWar)
	assert.Equal("", orderOfGloryJson.Guilds.Guild.DisbandedDate)
	assert.Equal("", orderOfGloryJson.Guilds.Guild.DisbandedCondition)
	assert.Equal(1, orderOfGloryJson.Guilds.Guild.PlayersOnline)
	assert.Equal(32, orderOfGloryJson.Guilds.Guild.PlayersOffline)
	assert.Equal(33, orderOfGloryJson.Guilds.Guild.MembersTotal)
	assert.Equal(0, orderOfGloryJson.Guilds.Guild.MembersInvited)
	assert.Equal(33, len(orderOfGloryJson.Guilds.Guild.Members))

	guildLeader := orderOfGloryJson.Guilds.Guild.Members[0]
	assert.Equal("Zyb the Warrior", guildLeader.Name)
	assert.Equal("", guildLeader.Title)
	assert.Equal("Leader", guildLeader.Rank)
	assert.Equal("Elite Knight", guildLeader.Vocation)
	assert.Equal(385, guildLeader.Level)
	assert.Equal("2020-10-13", guildLeader.Joined)
	assert.Equal("online", guildLeader.Status)

	assert.Nil(orderOfGloryJson.Guilds.Guild.Invited)
}
