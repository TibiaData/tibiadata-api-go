package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighscoreCategoryAchievementsString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := achievements
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("achievements", stringValue)
	assert.Equal(HighscoreCategory(1), highscoreCategory)
}

func TestHighscoreCategoryExperienceString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := experience
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("experience", stringValue)
	assert.Equal(HighscoreCategory(6), highscoreCategory)
}

func TestHighscoreCategoryDromescoreString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := dromescore
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("dromescore", stringValue)
	assert.Equal(HighscoreCategory(14), highscoreCategory)
}

func TestHighscoreCategoryInvalidValueString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighscoreCategory(99)
	_, err := highscoreCategory.String()

	assert.NotNil(err)
}

func TestHighscoreCategoryFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(experience, HighscoreCategoryFromString("experience"))
	assert.Equal(experience, HighscoreCategoryFromString(""))

	assert.Equal(achievements, HighscoreCategoryFromString("achievements"))
	assert.Equal(achievements, HighscoreCategoryFromString("achievement"))

	assert.Equal(axefighting, HighscoreCategoryFromString("axe"))
	assert.Equal(axefighting, HighscoreCategoryFromString("axefighting"))

	assert.Equal(charmpoints, HighscoreCategoryFromString("charm"))
	assert.Equal(charmpoints, HighscoreCategoryFromString("charms"))
	assert.Equal(charmpoints, HighscoreCategoryFromString("charmpoints"))

	assert.Equal(clubfighting, HighscoreCategoryFromString("club"))
	assert.Equal(clubfighting, HighscoreCategoryFromString("clubfighting"))

	assert.Equal(distancefighting, HighscoreCategoryFromString("distance"))
	assert.Equal(distancefighting, HighscoreCategoryFromString("distancefighting"))

	assert.Equal(fishing, HighscoreCategoryFromString("fishing"))

	assert.Equal(fistfighting, HighscoreCategoryFromString("fist"))
	assert.Equal(fistfighting, HighscoreCategoryFromString("fistfighting"))

	assert.Equal(goshnarstaint, HighscoreCategoryFromString("goshnar"))
	assert.Equal(goshnarstaint, HighscoreCategoryFromString("goshnars"))
	assert.Equal(goshnarstaint, HighscoreCategoryFromString("goshnarstaint"))

	assert.Equal(loyaltypoints, HighscoreCategoryFromString("loyalty"))
	assert.Equal(loyaltypoints, HighscoreCategoryFromString("loyaltypoints"))

	assert.Equal(magiclevel, HighscoreCategoryFromString("magic"))
	assert.Equal(magiclevel, HighscoreCategoryFromString("mlvl"))
	assert.Equal(magiclevel, HighscoreCategoryFromString("magiclevel"))

	assert.Equal(shielding, HighscoreCategoryFromString("shielding"))
	assert.Equal(shielding, HighscoreCategoryFromString("shield"))

	assert.Equal(swordfighting, HighscoreCategoryFromString("sword"))
	assert.Equal(swordfighting, HighscoreCategoryFromString("swordfighting"))

	assert.Equal(dromescore, HighscoreCategoryFromString("drome"))
	assert.Equal(dromescore, HighscoreCategoryFromString("dromescore"))
}
