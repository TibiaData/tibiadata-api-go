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
	assert.Empty(orderOfGloryJson.Guilds.Guild.Homepage)
	assert.False(orderOfGloryJson.Guilds.Guild.InWar)
	assert.Empty(orderOfGloryJson.Guilds.Guild.DisbandedDate)
	assert.Empty(orderOfGloryJson.Guilds.Guild.DisbandedCondition)
	assert.Equal(1, orderOfGloryJson.Guilds.Guild.PlayersOnline)
	assert.Equal(32, orderOfGloryJson.Guilds.Guild.PlayersOffline)
	assert.Equal(33, orderOfGloryJson.Guilds.Guild.MembersTotal)
	assert.Equal(0, orderOfGloryJson.Guilds.Guild.MembersInvited)
	assert.Equal(33, len(orderOfGloryJson.Guilds.Guild.Members))

	guildLeader := orderOfGloryJson.Guilds.Guild.Members[0]
	assert.Equal("Zyb the Warrior", guildLeader.Name)
	assert.Empty(guildLeader.Title)
	assert.Equal("Leader", guildLeader.Rank)
	assert.Equal("Elite Knight", guildLeader.Vocation)
	assert.Equal(385, guildLeader.Level)
	assert.Equal("2020-10-13", guildLeader.Joined)
	assert.Equal("online", guildLeader.Status)

	assert.Nil(orderOfGloryJson.Guilds.Guild.Invited)
}

func TestElysium(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Elysium.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	elysiumJson, _ := TibiaGuildsGuildV3Impl("Elysium", string(data))
	assert := assert.New(t)

	assert.Equal("Elysium", elysiumJson.Guilds.Guild.Name)
	assert.Equal("Vunira", elysiumJson.Guilds.Guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Elysium.gif", elysiumJson.Guilds.Guild.LogoURL)
	assert.Equal("The place you want to be...\nIt is the land of peace and harmony, the home of the immortal, the blessed, home of the passed away legends... Hail all defenders of righteousness and the old virtues which shall never be forgotten!\nIf you would like to join us, feel free to contact one of our leaders.", elysiumJson.Guilds.Guild.Description)
	assert.NotNil(elysiumJson.Guilds.Guild.Guildhalls)
	assert.Equal("Ab'Dendriel Clanhall", elysiumJson.Guilds.Guild.Guildhalls[0].Name)
	assert.Equal("2023-02-18", elysiumJson.Guilds.Guild.Guildhalls[0].PaidUntil)
	assert.Equal("Vunira", elysiumJson.Guilds.Guild.Guildhalls[0].World)
	assert.True(elysiumJson.Guilds.Guild.Active)
	assert.Equal("2004-05-26", elysiumJson.Guilds.Guild.Founded)
	assert.True(elysiumJson.Guilds.Guild.Applications)
	assert.Empty(elysiumJson.Guilds.Guild.Homepage)
	assert.False(elysiumJson.Guilds.Guild.InWar)
	assert.Empty(elysiumJson.Guilds.Guild.DisbandedDate)
	assert.Empty(elysiumJson.Guilds.Guild.DisbandedCondition)
	assert.Equal(4, elysiumJson.Guilds.Guild.PlayersOnline)
	assert.Equal(154, elysiumJson.Guilds.Guild.PlayersOffline)
	assert.Equal(158, elysiumJson.Guilds.Guild.MembersTotal)
	assert.Equal(1, elysiumJson.Guilds.Guild.MembersInvited)
	assert.Equal(158, len(elysiumJson.Guilds.Guild.Members))

	guildFollower := elysiumJson.Guilds.Guild.Members[101]
	assert.Equal("Trollefar", guildFollower.Name)
	assert.Equal("Troll Giant", guildFollower.Title)
	assert.Equal("Follower", guildFollower.Rank)
	assert.Equal("Elite Knight", guildFollower.Vocation)
	assert.Equal(202, guildFollower.Level)
	assert.Equal("2013-10-20", guildFollower.Joined)
	assert.Equal("offline", guildFollower.Status)

	assert.NotNil(elysiumJson.Guilds.Guild.Invited)
	evelynInvite := elysiumJson.Guilds.Guild.Invited[0]
	assert.Equal("Evelyn Earlong", evelynInvite.Name)
	assert.Equal("2023-01-20", evelynInvite.Date)
}

func TestMercenarys(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Mercenarys.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	mercenarysJson, _ := TibiaGuildsGuildV3Impl("Mercenarys", string(data))
	assert := assert.New(t)

	assert.Equal("Mercenarys", mercenarysJson.Guilds.Guild.Name)
	assert.Equal("Antica", mercenarysJson.Guilds.Guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Mercenarys.gif", mercenarysJson.Guilds.Guild.LogoURL)
	assert.NotNil(mercenarysJson.Guilds.Guild.Guildhalls)
	assert.Equal("Mercenary Tower", mercenarysJson.Guilds.Guild.Guildhalls[0].Name)
	assert.Equal("2023-01-28", mercenarysJson.Guilds.Guild.Guildhalls[0].PaidUntil)
	assert.Equal("Antica", mercenarysJson.Guilds.Guild.Guildhalls[0].World)
	assert.True(mercenarysJson.Guilds.Guild.Active)
	assert.Equal("2002-02-18", mercenarysJson.Guilds.Guild.Founded)
	assert.True(mercenarysJson.Guilds.Guild.Applications)
	assert.Equal("http://www.mercenarys.net", mercenarysJson.Guilds.Guild.Homepage)
	assert.False(mercenarysJson.Guilds.Guild.InWar)
	assert.Equal("2023-02-07", mercenarysJson.Guilds.Guild.DisbandedDate)
	assert.Equal("if there are still less than four vice leaders or an insufficient amount of premium accounts in the leading ranks by then", mercenarysJson.Guilds.Guild.DisbandedCondition)
}

func TestKotkiAntica(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Kotki Antica.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	kotkianticaJson, _ := TibiaGuildsGuildV3Impl("Kotki Antica", string(data))
	assert := assert.New(t)

	assert.Equal("Kotki Antica", kotkianticaJson.Guilds.Guild.Name)
	assert.Equal("Antica", kotkianticaJson.Guilds.Guild.World)
	assert.Empty(kotkianticaJson.Guilds.Guild.Description)
	assert.True(kotkianticaJson.Guilds.Guild.Active)
	assert.Equal("2021-09-22", kotkianticaJson.Guilds.Guild.Founded)
	assert.False(kotkianticaJson.Guilds.Guild.Applications)
}
