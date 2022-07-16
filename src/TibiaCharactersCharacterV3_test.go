package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/TibiaData/tibiadata-api-go/src/validation"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNumber1(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Darkside Rafa.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))
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
	assert.Equal("Trial", characterJson.Characters.Character.Guild.Rank)
	assert.Equal("2022-01-05T21:23:32Z", characterJson.Characters.Character.LastLogin)
	assert.Equal("Premium Account", characterJson.Characters.Character.AccountStatus)
	assert.Empty(characterJson.Characters.Character.Comment)
}

func TestNumber2(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Zugspitze Housekeeper.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))
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
	assert.Equal("Magnus Magister", characterJson.Characters.Character.Guild.Rank)
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
	file, err := static.TestFiles.Open("testdata/characters/Borttagna Gubben.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Borttagna Gubben", characterJson.Characters.Character.Name)
	assert.True(characterJson.Characters.Character.Traded)
	assert.Equal(2, len(characterJson.Characters.Character.FormerWorlds))
	assert.Equal("Zuna", characterJson.Characters.Character.FormerWorlds[0])
	assert.Equal("Zunera", characterJson.Characters.Character.FormerWorlds[1])
	assert.Equal("Bubble", characterJson.Characters.Character.MarriedTo)
	assert.Equal("2022-03-08T00:09:13Z", characterJson.Characters.Character.DeletionDate)
	assert.Equal("", characterJson.Characters.Character.LastLogin)
	assert.Equal("Free Account", characterJson.Characters.Character.AccountStatus)
	assert.Equal("Fansite Admin", characterJson.Characters.AccountInformation.Position)
	assert.Empty(characterJson.Characters.AccountInformation.LoyaltyTitle)
}

func TestNumber4(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Riley No Hands.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))
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

	creatureWithOfDeath := characterJson.Characters.Deaths[16]
	assert.Equal(2, len(creatureWithOfDeath.Killers))
	assert.Equal(260, creatureWithOfDeath.Level)
	assert.Equal("an undead elite gladiator", creatureWithOfDeath.Killers[0].Name)
	assert.False(creatureWithOfDeath.Killers[0].Player)
	assert.False(creatureWithOfDeath.Killers[0].Traded)
	assert.Empty(creatureWithOfDeath.Killers[0].Summon)
	assert.Equal("a priestess of the wild sun", creatureWithOfDeath.Killers[1].Name)

	tradedInDeath := characterJson.Characters.Deaths[18]
	assert.Equal(3, len(tradedInDeath.Assists))
	assert.Equal(261, tradedInDeath.Level)
	assert.Equal("Vithrann", tradedInDeath.Assists[1].Name)
	assert.True(tradedInDeath.Assists[1].Player)
	assert.Equal("Adam No Hands", tradedInDeath.Assists[2].Name)
	assert.True(tradedInDeath.Assists[2].Traded)

	longDeath := characterJson.Characters.Deaths[78]
	assert.Equal(5, len(longDeath.Assists))
	assert.Equal(231, longDeath.Level)
	assert.Equal("Adam No Hands", longDeath.Assists[4].Name)
	assert.Equal("a paladin familiar", longDeath.Assists[4].Summon)
	assert.True(longDeath.Assists[4].Traded)
}

func TestNumber5(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Torbjörn.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("Torbjörn", characterJson.Characters.Character.Name)
	assert.Equal("___$$$$$$$$_______$$$$$$$$\n_$$$$$$$$$$$$__$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n_$$$$$$$$$$-Snulliz-$$$$$$$$$$$\n__$$$$$$$$$$$$$$$$$$$$$$$$$$\n____$$$$$$$$$$$$$$$$$$$$$$\n______$$$$$$$$$$$$$$$$$$\n________$$$$$$$$$$$$$$\n___________$$$$$$$$$\n____________$$$$$$\n_____________$$", characterJson.Characters.Character.Comment)
}

func BenchmarkNumber1(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/characters/Darkside Rafa.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))

		assert.Equal(b, "Darkside Rafa", characterJson.Characters.Character.Name)
	}
}

func BenchmarkNumber2(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/characters/Riley No Hands.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson, _ := TibiaCharactersCharacterV3Impl(string(data))

		assert.Equal(b, "Riley No Hands", characterJson.Characters.Character.Name)
	}
}

func TestEmptyName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "",
		},
	}

	tibiaCharactersCharacterV3(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(jerr)

	assert.EqualValues(validation.ErrorCharacterNameEmpty.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameEmpty.Error(), jerr.Information.Status.Message)
}

func TestSmallName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "o",
		},
	}

	tibiaCharactersCharacterV3(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(validation.ErrorCharacterNameTooSmall.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameTooSmall.Error(), jerr.Information.Status.Message)
}

func TestInvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "12",
		},
	}

	tibiaCharactersCharacterV3(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(validation.ErrorCharacterNameInvalid.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameInvalid.Error(), jerr.Information.Status.Message)
}
