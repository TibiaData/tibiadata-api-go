package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber1(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Darkside Rafa.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Darkside Rafa", characterJson.Characters.Character.Name)
	assert.Nil(characterJson.Characters.Character.FormerNames)
	assert.False(characterJson.Characters.Character.Traded)
	assert.Empty(characterJson.Characters.Character.DeletionDate)
	assert.Equal("male", characterJson.Characters.Character.Sex)
	assert.Equal("Silencer", characterJson.Characters.Character.Title)
	assert.Equal(18, characterJson.Characters.Character.UnlockedTitles)
	assert.Equal("Elite Knight", characterJson.Characters.Character.Vocation)
	assert.Equal(790, characterJson.Characters.Character.Level)
	assert.Equal(596, characterJson.Characters.Character.AchievementPoints)
	assert.Equal("Gladera", characterJson.Characters.Character.World)
	assert.Nil(characterJson.Characters.Character.FormerWorlds)
	assert.Equal("Thais", characterJson.Characters.Character.Residence)
	assert.Empty(characterJson.Characters.Character.MarriedTo)
	assert.Nil(characterJson.Characters.Character.Houses)
	assert.Equal("Jokerz", characterJson.Characters.Character.Guild.GuildName)
	assert.Equal("Trial of the ", characterJson.Characters.Character.Guild.Rank)
	assert.Equal("2022-01-05T21:23:32Z", characterJson.Characters.Character.LastLogin)
	assert.Equal("Premium Account", characterJson.Characters.Character.AccountStatus)
	assert.Empty(characterJson.Characters.Character.Comment)
}

func TestNumber2(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Zugspitze Housekeeper.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Zugspitze Housekeeper", characterJson.Characters.Character.Name)
	assert.Nil(characterJson.Characters.Character.FormerNames)
	assert.False(characterJson.Characters.Character.Traded)
	assert.Empty(characterJson.Characters.Character.DeletionDate)
	assert.Equal("male", characterJson.Characters.Character.Sex)
	assert.Equal("None", characterJson.Characters.Character.Title)
	assert.Equal(13, characterJson.Characters.Character.UnlockedTitles)
	assert.Equal("Elite Knight", characterJson.Characters.Character.Vocation)
	assert.Equal(79, characterJson.Characters.Character.Level)
	assert.Equal(262, characterJson.Characters.Character.AchievementPoints)
	assert.Equal("Venebra", characterJson.Characters.Character.World)
	assert.Nil(characterJson.Characters.Character.FormerWorlds)
	assert.Equal("Darashia", characterJson.Characters.Character.Residence)
	assert.Empty(characterJson.Characters.Character.MarriedTo)
	assert.Equal(35056, characterJson.Characters.Character.Houses[0].HouseID)
	assert.Equal("Loot Lane 1 (Shop)", characterJson.Characters.Character.Houses[0].Name)
	assert.Equal("Venore", characterJson.Characters.Character.Houses[0].Town)
	assert.Equal("2022-01-16", characterJson.Characters.Character.Houses[0].Paid)
	assert.Equal("Magnus Magister of the ", characterJson.Characters.Character.Guild.Rank)
	assert.Equal("Lionheart Society", characterJson.Characters.Character.Guild.GuildName)
	assert.Equal("2022-01-06T21:38:44Z", characterJson.Characters.Character.LastLogin)
	assert.Equal("Testa de Ferro do Lejonhjartat ;)", characterJson.Characters.Character.Comment)
	assert.Equal("Premium Account", characterJson.Characters.Character.AccountStatus)

	//validate other characters
	assert.Equal(7, len(characterJson.Characters.OtherCharacters))

	onlineMainCharacter := characterJson.Characters.OtherCharacters[3]
	assert.Equal("Lejonhjartat", onlineMainCharacter.Name)
	assert.Equal("Venebra", onlineMainCharacter.World)
	assert.Equal("online", onlineMainCharacter.Status)
	assert.Equal(false, onlineMainCharacter.Deleted)
	assert.Equal(true, onlineMainCharacter.Main)
	assert.Equal(false, onlineMainCharacter.Traded)

	offlineCharacter := characterJson.Characters.OtherCharacters[5]
	assert.Equal("Oak Knight Disruptivo", offlineCharacter.Name)
	assert.Equal("Libertabra", offlineCharacter.World)
	assert.Equal("offline", offlineCharacter.Status)
	assert.Equal(false, offlineCharacter.Deleted)
	assert.Equal(false, offlineCharacter.Main)
	assert.Equal(false, offlineCharacter.Traded)
}

func TestNumber3(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Borttagna Gubben.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Borttagna Gubben", characterJson.Characters.Character.Name)
	assert.Equal("2022-03-08T00:09:13Z", characterJson.Characters.Character.DeletionDate)
	assert.Equal("", characterJson.Characters.Character.LastLogin)
	assert.Equal("Free Account", characterJson.Characters.Character.AccountStatus)
}

func TestNumber4(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Riley No Hands.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Riley No Hands", characterJson.Characters.Character.Name)

	//validate former names
	assert.Equal(4, len(characterJson.Characters.Character.FormerNames))
	assert.Equal("Dura Malandro", characterJson.Characters.Character.FormerNames[0])
	assert.Equal("Letsgo Brandon", characterJson.Characters.Character.FormerNames[1])
	assert.Equal("Letsgo Brandon", characterJson.Characters.Character.FormerNames[2]) //yes, this name is listed twice
	assert.Equal("Nataraya Soldrac", characterJson.Characters.Character.FormerNames[3])

	//validate death data
	assert.Equal(79, len(characterJson.Characters.Deaths))

	firstDeath := characterJson.Characters.Deaths[0]
	assert.Equal(28, len(firstDeath.Killers))
}

func TestNumber5(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Torbjörn.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Torbjörn", characterJson.Characters.Character.Name)
	assert.Equal("___$$$$$$$$_______$$$$$$$$\n_$$$$$$$$$$$$__$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n_$$$$$$$$$$-Snulliz-$$$$$$$$$$$\n__$$$$$$$$$$$$$$$$$$$$$$$$$$\n____$$$$$$$$$$$$$$$$$$$$$$\n______$$$$$$$$$$$$$$$$$$\n________$$$$$$$$$$$$$$\n___________$$$$$$$$$\n____________$$$$$$\n_____________$$", characterJson.Characters.Character.Comment)
}

func BenchmarkNumber1(b *testing.B) {
	data, err := ioutil.ReadFile("../testdata/characters/Darkside Rafa.html")
	if err != nil {
		b.Errorf("File reading error: %s", err)
		return
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson := TibiaCharactersCharacterV3Impl(string(data))

		assert.Equal(b, "Darkside Rafa", characterJson.Characters.Character.Name)
	}
}

func BenchmarkNumber2(b *testing.B) {
	data, err := ioutil.ReadFile("../testdata/characters/Riley No Hands.html")
	if err != nil {
		b.Errorf("File reading error: %s", err)
		return
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson := TibiaCharactersCharacterV3Impl(string(data))

		assert.Equal(b, "Riley No Hands", characterJson.Characters.Character.Name)
	}
}
