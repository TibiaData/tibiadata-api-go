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

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Darkside Rafa", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("male", character.Sex)
	assert.Equal("Silencer", character.Title)
	assert.Equal(18, character.UnlockedTitles)
	assert.Equal("Elite Knight", character.Vocation)
	assert.Equal(790, character.Level)
	assert.Equal(596, character.AchievementPoints)
	assert.Equal("Gladera", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Thais", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Nil(character.Houses)
	assert.Equal("Jokerz", character.Guild.GuildName)
	assert.Equal("Trial", character.Guild.Rank)
	assert.Equal("2022-01-05T21:23:32Z", character.LastLogin)
	assert.Equal("Premium Account", character.AccountStatus)
	assert.Empty(character.Comment)
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

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Zugspitze Housekeeper", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("male", character.Sex)
	assert.Equal("None", character.Title)
	assert.Equal(13, character.UnlockedTitles)
	assert.Equal("Elite Knight", character.Vocation)
	assert.Equal(79, character.Level)
	assert.Equal(262, character.AchievementPoints)
	assert.Equal("Venebra", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Darashia", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Equal(35056, character.Houses[0].HouseID)
	assert.Equal("Loot Lane 1 (Shop)", character.Houses[0].Name)
	assert.Equal("Venore", character.Houses[0].Town)
	assert.Equal("2022-01-16", character.Houses[0].Paid)
	assert.Equal("Magnus Magister", character.Guild.Rank)
	assert.Equal("Lionheart Society", character.Guild.GuildName)
	assert.Equal("2022-01-06T21:38:44Z", character.LastLogin)
	assert.Empty(character.Position)
	assert.Equal("Testa de Ferro do Lejonhjartat ;)", character.Comment)
	assert.Equal("Premium Account", character.AccountStatus)

	//validate other characters
	assert.Equal(7, len(characterJson.Character.OtherCharacters))

	onlineMainCharacter := characterJson.Character.OtherCharacters[3]
	assert.Equal("Lejonhjartat", onlineMainCharacter.Name)
	assert.Equal("Venebra", onlineMainCharacter.World)
	assert.Equal("online", onlineMainCharacter.Status)
	assert.False(onlineMainCharacter.Deleted)
	assert.True(onlineMainCharacter.Main)
	assert.False(onlineMainCharacter.Traded)
	assert.Empty(onlineMainCharacter.Position)

	offlineCharacter := characterJson.Character.OtherCharacters[5]
	assert.Equal("Oak Knight Disruptivo", offlineCharacter.Name)
	assert.Equal("Libertabra", offlineCharacter.World)
	assert.Equal("offline", offlineCharacter.Status)
	assert.False(offlineCharacter.Deleted)
	assert.False(offlineCharacter.Main)
	assert.False(offlineCharacter.Traded)
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

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Borttagna Gubben", character.Name)
	assert.True(character.Traded)
	assert.Equal(2, len(character.FormerWorlds))
	assert.Equal("Zuna", character.FormerWorlds[0])
	assert.Equal("Zunera", character.FormerWorlds[1])
	assert.Equal("Bubble", character.MarriedTo)
	assert.Equal("2022-03-08T00:09:13Z", character.DeletionDate)
	assert.Empty(character.LastLogin)
	assert.Equal("Free Account", character.AccountStatus)
	assert.Equal("Fansite Admin", characterJson.Character.AccountInformation.Position)
	assert.Empty(characterJson.Character.AccountInformation.LoyaltyTitle)
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

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Riley No Hands", character.Name)

	//validate former names
	assert.Equal(4, len(character.FormerNames))
	assert.Equal("Dura Malandro", character.FormerNames[0])
	assert.Equal("Letsgo Brandon", character.FormerNames[1])
	assert.Equal("Letsgo Brandon", character.FormerNames[2]) //yes, this name is listed twice
	assert.Equal("Nataraya Soldrac", character.FormerNames[3])

	//validate death data
	assert.Equal(79, len(characterJson.Character.Deaths))

	firstDeath := characterJson.Character.Deaths[0]
	assert.Equal(28, len(firstDeath.Killers))

	creatureWithOfDeath := characterJson.Character.Deaths[16]
	assert.Equal(2, len(creatureWithOfDeath.Killers))
	assert.Equal(260, creatureWithOfDeath.Level)
	assert.Equal("an undead elite gladiator", creatureWithOfDeath.Killers[0].Name)
	assert.False(creatureWithOfDeath.Killers[0].Player)
	assert.False(creatureWithOfDeath.Killers[0].Traded)
	assert.Empty(creatureWithOfDeath.Killers[0].Summon)
	assert.Equal("a priestess of the wild sun", creatureWithOfDeath.Killers[1].Name)

	tradedInDeath := characterJson.Character.Deaths[18]
	assert.Equal(3, len(tradedInDeath.Assists))
	assert.Equal(261, tradedInDeath.Level)
	assert.Equal("Vithrann", tradedInDeath.Assists[1].Name)
	assert.True(tradedInDeath.Assists[1].Player)
	assert.Equal("Adam No Hands", tradedInDeath.Assists[2].Name)
	assert.True(tradedInDeath.Assists[2].Traded)

	longDeath := characterJson.Character.Deaths[78]
	assert.Equal(5, len(longDeath.Assists))
	assert.Equal(231, longDeath.Level)
	assert.Equal("Adam No Hands", longDeath.Assists[4].Name)
	assert.Equal("a paladin familiar", longDeath.Assists[4].Summon)
	assert.True(longDeath.Assists[4].Traded)
}

func TestNumber5(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Rejana on Fera.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Rejana on Fera", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("female", character.Sex)
	assert.Equal("None", character.Title)
	assert.Equal(12, character.UnlockedTitles)
	assert.Equal("Knight", character.Vocation)
	assert.Equal(8, character.Level)
	assert.Equal(0, character.AchievementPoints)
	assert.Equal("Fera", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Isle of Solitude", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Equal("2021-10-25T04:37:46Z", character.LastLogin)
	assert.Equal("CipSoft Member", character.Position)
	assert.Equal("Premium Account", character.AccountStatus)

	assert.Empty(characterJson.Character.Achievements)
	assert.Empty(characterJson.Character.AccountInformation.LoyaltyTitle)

	//validate other characters
	assert.Equal(4, len(characterJson.Character.OtherCharacters))

	positionCharacter := characterJson.Character.OtherCharacters[1]
	assert.Equal("Rejana on Fera", positionCharacter.Name)
	assert.Equal("Fera", positionCharacter.World)
	assert.Equal("offline", positionCharacter.Status)
	assert.False(positionCharacter.Deleted)
	assert.False(positionCharacter.Main)
	assert.False(positionCharacter.Traded)
	assert.Equal("CipSoft Member", positionCharacter.Position)
}

func TestNumber6(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Luminals.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Luminals", character.Name)

	assert.Equal(8, len(characterJson.Character.AccountBadges))
	globalPlayerBadge := characterJson.Character.AccountBadges[4]
	assert.Equal("Global Player (Grade 3)", globalPlayerBadge.Name)
	assert.Equal("https://static.tibia.com/images//badges/badge_globalplayer3.png", globalPlayerBadge.IconURL)
	assert.Equal("Summing up the levels of all characters on the account amounts to at least 2000.", globalPlayerBadge.Description)
	masterClassBadge := characterJson.Character.AccountBadges[7]
	assert.Equal("Master Class (Grade 1)", masterClassBadge.Name)
	assert.Equal("https://static.tibia.com/images//badges/badge_masterclass1.png", masterClassBadge.IconURL)
	assert.Equal("The account has reached at least level 100 with all four vocations.", masterClassBadge.Description)

	assert.Len(characterJson.Character.Achievements, 5)
	assert.Equal(characterJson.Character.Achievements[0].Name, "Alumni")
	assert.Equal(characterJson.Character.Achievements[0].Grade, 2)
	assert.Equal(characterJson.Character.Achievements[0].Secret, false)
	assert.Equal(characterJson.Character.Achievements[1].Name, "Forbidden Fruit")
	assert.Equal(characterJson.Character.Achievements[1].Grade, 1)
	assert.Equal(characterJson.Character.Achievements[1].Secret, true)
	assert.Equal(characterJson.Character.Achievements[2].Name, "Goldhunter")
	assert.Equal(characterJson.Character.Achievements[2].Grade, 1)
	assert.Equal(characterJson.Character.Achievements[2].Secret, true)
	assert.Equal(characterJson.Character.Achievements[3].Name, "Pyromaniac")
	assert.Equal(characterJson.Character.Achievements[3].Grade, 2)
	assert.Equal(characterJson.Character.Achievements[3].Secret, true)
	assert.Equal(characterJson.Character.Achievements[4].Name, "Razing!")
	assert.Equal(characterJson.Character.Achievements[4].Grade, 3)
	assert.Equal(characterJson.Character.Achievements[4].Secret, true)
}

func TestNumber7(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Torbjörn.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Len(characterJson.Character.Achievements, 0)
	assert.Equal("Torbjörn", character.Name)
	assert.Equal("___$$$$$$$$_______$$$$$$$$\n_$$$$$$$$$$$$__$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n_$$$$$$$$$$-Snulliz-$$$$$$$$$$$\n__$$$$$$$$$$$$$$$$$$$$$$$$$$\n____$$$$$$$$$$$$$$$$$$$$$$\n______$$$$$$$$$$$$$$$$$$\n________$$$$$$$$$$$$$$\n___________$$$$$$$$$\n____________$$$$$$\n_____________$$", character.Comment)
}

func TestNumber8(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Jowjow Invencivel.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	character := characterJson.Character
	assert.Len(character.Achievements, 5)
	assert.Equal(character.Achievements[0].Name, "Alumni")
	assert.Equal(character.Achievements[0].Grade, 2)
	assert.Equal(character.Achievements[0].Secret, false)
	assert.Equal(character.Achievements[1].Name, "Bad Timing")
	assert.Equal(character.Achievements[1].Grade, 1)
	assert.Equal(character.Achievements[1].Secret, true)
	assert.Equal(character.Achievements[2].Name, "Cake Conqueror")
	assert.Equal(character.Achievements[2].Grade, 1)
	assert.Equal(character.Achievements[2].Secret, true)
	assert.Equal(character.Achievements[3].Name, "Hat Hunter")
	assert.Equal(character.Achievements[3].Grade, 2)
	assert.Equal(character.Achievements[3].Secret, false)
	assert.Equal(character.Achievements[4].Name, "Number of the Beast")
	assert.Equal(character.Achievements[4].Grade, 1)
	assert.Equal(character.Achievements[4].Secret, false)
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
		characterJson, _ := TibiaCharactersCharacterImpl(string(data))

		assert.Equal(b, "Darkside Rafa", characterJson.Character.CharacterInfo.Name)
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
		characterJson, _ := TibiaCharactersCharacterImpl(string(data))

		assert.Equal(b, "Riley No Hands", characterJson.Character.CharacterInfo.Name)
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

	tibiaCharactersCharacter(c)
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

	tibiaCharactersCharacter(c)
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

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(validation.ErrorCharacterNameInvalid.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameInvalid.Error(), jerr.Information.Status.Message)
}
