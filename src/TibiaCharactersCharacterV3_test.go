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

	assert.Equal(t, "Darkside Rafa", characterJson.Characters.Character.Name)
	assert.Nil(t, characterJson.Characters.Character.FormerNames)
	assert.False(t, characterJson.Characters.Character.Traded)
	assert.Empty(t, characterJson.Characters.Character.DeletionDate)
	assert.Equal(t, "male", characterJson.Characters.Character.Sex)
	assert.Equal(t, "Silencer", characterJson.Characters.Character.Title)
	assert.Equal(t, 18, characterJson.Characters.Character.UnlockedTitles)
	assert.Equal(t, "Elite Knight", characterJson.Characters.Character.Vocation)
	assert.Equal(t, 790, characterJson.Characters.Character.Level)
	assert.Equal(t, 596, characterJson.Characters.Character.AchievementPoints)
	assert.Equal(t, "Gladera", characterJson.Characters.Character.World)
	assert.Nil(t, characterJson.Characters.Character.FormerWorlds)
	assert.Equal(t, "Thais", characterJson.Characters.Character.Residence)
	assert.Empty(t, characterJson.Characters.Character.MarriedTo)
	assert.Nil(t, characterJson.Characters.Character.Houses)
	assert.Equal(t, "Jokerz", characterJson.Characters.Character.Guild.GuildName)
	assert.Equal(t, "Trial of the ", characterJson.Characters.Character.Guild.Rank)
	assert.Equal(t, "2022-01-05T21:23:32Z", characterJson.Characters.Character.LastLogin)
	assert.Equal(t, "Premium Account", characterJson.Characters.Character.AccountStatus)
	assert.Empty(t, characterJson.Characters.Character.Comment)
}

func TestNumber2(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Zugspitze Housekeeper.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))

	assert.Equal(t, "Zugspitze Housekeeper", characterJson.Characters.Character.Name)
	assert.Nil(t, characterJson.Characters.Character.FormerNames)
	assert.False(t, characterJson.Characters.Character.Traded)
	assert.Empty(t, characterJson.Characters.Character.DeletionDate)
	assert.Equal(t, "male", characterJson.Characters.Character.Sex)
	assert.Equal(t, "None", characterJson.Characters.Character.Title)
	assert.Equal(t, 13, characterJson.Characters.Character.UnlockedTitles)
	assert.Equal(t, "Elite Knight", characterJson.Characters.Character.Vocation)
	assert.Equal(t, 79, characterJson.Characters.Character.Level)
	assert.Equal(t, 262, characterJson.Characters.Character.AchievementPoints)
	assert.Equal(t, "Venebra", characterJson.Characters.Character.World)
	assert.Nil(t, characterJson.Characters.Character.FormerWorlds)
	assert.Equal(t, "Darashia", characterJson.Characters.Character.Residence)
	assert.Empty(t, characterJson.Characters.Character.MarriedTo)
	assert.Equal(t, 35056, characterJson.Characters.Character.Houses[0].HouseID)
	assert.Equal(t, "Loot Lane 1 (Shop)", characterJson.Characters.Character.Houses[0].Name)
	assert.Equal(t, "Venore", characterJson.Characters.Character.Houses[0].Town)
	assert.Equal(t, "2022-01-16", characterJson.Characters.Character.Houses[0].Paid)
	assert.Equal(t, "Magnus Magister of the ", characterJson.Characters.Character.Guild.Rank)
	assert.Equal(t, "Lionheart Society", characterJson.Characters.Character.Guild.GuildName)
	assert.Equal(t, "2022-01-06T21:38:44Z", characterJson.Characters.Character.LastLogin)
	assert.Equal(t, "Testa de Ferro do Lejonhjartat ;)", characterJson.Characters.Character.Comment)
	assert.Equal(t, "Premium Account", characterJson.Characters.Character.AccountStatus)

	//validate other characters
	assert.Equal(t, 7, len(characterJson.Characters.OtherCharacters))

	onlineMainCharacter := characterJson.Characters.OtherCharacters[3]
	assert.Equal(t, "Lejonhjartat", onlineMainCharacter.Name)
	assert.Equal(t, "Venebra", onlineMainCharacter.World)
	assert.Equal(t, "online", onlineMainCharacter.Status)
	assert.Equal(t, false, onlineMainCharacter.Deleted)
	assert.Equal(t, true, onlineMainCharacter.Main)
	assert.Equal(t, false, onlineMainCharacter.Traded)

	offlineCharacter := characterJson.Characters.OtherCharacters[5]
	assert.Equal(t, "Oak Knight Disruptivo", offlineCharacter.Name)
	assert.Equal(t, "Libertabra", offlineCharacter.World)
	assert.Equal(t, "offline", offlineCharacter.Status)
	assert.Equal(t, false, offlineCharacter.Deleted)
	assert.Equal(t, false, offlineCharacter.Main)
	assert.Equal(t, false, offlineCharacter.Traded)
}

func TestNumber3(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Borttagna Gubben.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))

	assert.Equal(t, "Borttagna Gubben", characterJson.Characters.Character.Name)
	assert.Equal(t, "2022-03-08T00:09:13Z", characterJson.Characters.Character.DeletionDate)
	assert.Equal(t, "", characterJson.Characters.Character.LastLogin)
	assert.Equal(t, "Free Account", characterJson.Characters.Character.AccountStatus)
}

func TestNumber4(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Riley No Hands.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))

	assert.Equal(t, "Riley No Hands", characterJson.Characters.Character.Name)

	//validate former names
	assert.Equal(t, 4, len(characterJson.Characters.Character.FormerNames))
	assert.Equal(t, "Dura Malandro", characterJson.Characters.Character.FormerNames[0])
	assert.Equal(t, "Letsgo Brandon", characterJson.Characters.Character.FormerNames[1])
	assert.Equal(t, "Letsgo Brandon", characterJson.Characters.Character.FormerNames[2]) //yes, this name is listed twice
	assert.Equal(t, "Nataraya Soldrac", characterJson.Characters.Character.FormerNames[3])

	//validate death data
	assert.Equal(t, 79, len(characterJson.Characters.Deaths.DeathEntries))

	firstDeath := characterJson.Characters.Deaths.DeathEntries[0]
	assert.Equal(t, 28, len(firstDeath.Killers))
}

func TestNumber5(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/characters/Torbjörn.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	characterJson := TibiaCharactersCharacterV3Impl(string(data))

	assert.Equal(t, "Torbjörn", characterJson.Characters.Character.Name)
	assert.Equal(t, "___$$$$$$$$_______$$$$$$$$\n_$$$$$$$$$$$$__$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n_$$$$$$$$$$-Snulliz-$$$$$$$$$$$\n__$$$$$$$$$$$$$$$$$$$$$$$$$$\n____$$$$$$$$$$$$$$$$$$$$$$\n______$$$$$$$$$$$$$$$$$$\n________$$$$$$$$$$$$$$\n___________$$$$$$$$$\n____________$$$$$$\n_____________$$", characterJson.Characters.Character.Comment)
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
