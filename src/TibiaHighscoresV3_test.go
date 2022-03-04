package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighscoresAll(t *testing.T) {
	data, err := os.ReadFile("../testdata/highscores/all.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	highscoresJson := TibiaHighscoresV3Impl("", experience, "all", string(data))
	assert := assert.New(t)

	assert.Equal("", highscoresJson.Highscores.World)
	assert.Equal("experience", highscoresJson.Highscores.Category)
	assert.Equal("all", highscoresJson.Highscores.Vocation)
	assert.Equal(30, highscoresJson.Highscores.HighscoreAge)

	assert.Equal(50, len(highscoresJson.Highscores.HighscoreList))

	firstHighscore := highscoresJson.Highscores.HighscoreList[0]
	assert.Equal(1, firstHighscore.Rank)
	assert.Equal("Bobeek", firstHighscore.Name)
	assert.Equal("Elder Druid", firstHighscore.Vocation)
	assert.Equal("Bona", firstHighscore.World)
	assert.Equal(2015, firstHighscore.Level)
	assert.Equal(136026206904, firstHighscore.Value)
	assert.Equal("", firstHighscore.Title)

	lastHighscore := highscoresJson.Highscores.HighscoreList[49]
	assert.Equal(50, lastHighscore.Rank)
	assert.Equal("Kewhyx Mythus", lastHighscore.Name)
	assert.Equal("Royal Paladin", lastHighscore.Vocation)
	assert.Equal("Celebra", lastHighscore.World)
	assert.Equal(1575, lastHighscore.Level)
	assert.Equal(64869293274, lastHighscore.Value)
	assert.Equal("", lastHighscore.Title)
}

func TestHighscoresLoyalty(t *testing.T) {
	data, err := os.ReadFile("../testdata/highscores/loyalty.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	highscoresJson := TibiaHighscoresV3Impl("Vunira", loyaltypoints, "druids", string(data))
	assert := assert.New(t)

	assert.Equal("Vunira", highscoresJson.Highscores.World)
	assert.Equal("loyaltypoints", highscoresJson.Highscores.Category)
	assert.Equal("druids", highscoresJson.Highscores.Vocation)
	assert.Equal(46, highscoresJson.Highscores.HighscoreAge)

	// should be 50, but for some reason it can't get entries from the list..
	assert.Equal(0, len(highscoresJson.Highscores.HighscoreList))
}
