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

	highscoresJson := TibiaHighscoresV3Impl("", experience, "all", 1, string(data))
	assert := assert.New(t)

	assert.Equal("", highscoresJson.Highscores.World)
	assert.Equal("experience", highscoresJson.Highscores.Category)
	assert.Equal("all", highscoresJson.Highscores.Vocation)
	assert.Equal(12, highscoresJson.Highscores.HighscoreAge)

	assert.Equal(50, len(highscoresJson.Highscores.HighscoreList))

	assert.Equal(1, highscoresJson.Highscores.HighscorePage.CurrentPage)
	assert.Equal(20, highscoresJson.Highscores.HighscorePage.TotalPages)
	assert.Equal(1000, highscoresJson.Highscores.HighscorePage.TotalHighscores)

	firstHighscore := highscoresJson.Highscores.HighscoreList[0]
	assert.Equal(1, firstHighscore.Rank)
	assert.Equal("Goraca", firstHighscore.Name)
	assert.Equal("Master Sorcerer", firstHighscore.Vocation)
	assert.Equal("Bona", firstHighscore.World)
	assert.Equal(2197, firstHighscore.Level)
	assert.Equal(176271164607, firstHighscore.Value)
	assert.Equal("", firstHighscore.Title)

	lastHighscore := highscoresJson.Highscores.HighscoreList[49]
	assert.Equal(50, lastHighscore.Rank)
	assert.Equal("Wujo Daro", lastHighscore.Name)
	assert.Equal("Elite Knight", lastHighscore.Vocation)
	assert.Equal("Refugia", lastHighscore.World)
	assert.Equal(1701, lastHighscore.Level)
	assert.Equal(81816135617, lastHighscore.Value)
	assert.Equal("", lastHighscore.Title)
}

func TestHighscoresLoyalty(t *testing.T) {
	data, err := os.ReadFile("../testdata/highscores/loyalty.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	highscoresJson := TibiaHighscoresV3Impl("Vunira", loyaltypoints, "druids", 4, string(data))
	assert := assert.New(t)

	assert.Equal("Vunira", highscoresJson.Highscores.World)
	assert.Equal("loyaltypoints", highscoresJson.Highscores.Category)
	assert.Equal("druids", highscoresJson.Highscores.Vocation)
	assert.Equal(12, highscoresJson.Highscores.HighscoreAge)

	assert.Equal(50, len(highscoresJson.Highscores.HighscoreList))
}
