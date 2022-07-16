package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighscoreCategoryAchievementsString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreAchievements
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("achievements", stringValue)
	assert.Equal(HighscoreCategory(1), highscoreCategory)
}

func TestHighscoreCategoryExperienceString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreExperience
	stringValue, err := highscoreCategory.String()

	assert.Nil(err)
	assert.Equal("experience", stringValue)
	assert.Equal(HighscoreCategory(6), highscoreCategory)
}

func TestHighscoreCategoryDromescoreString(t *testing.T) {
	assert := assert.New(t)

	highscoreCategory := HighScoreDromescore
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

	assert.Equal(HighScoreExperience, HighscoreCategoryFromString("experience"))
	assert.Equal(HighScoreExperience, HighscoreCategoryFromString(""))

	assert.Equal(HighScoreAchievements, HighscoreCategoryFromString("achievements"))
	assert.Equal(HighScoreAchievements, HighscoreCategoryFromString("achievement"))

	assert.Equal(HighScoreAxefighting, HighscoreCategoryFromString("axe"))
	assert.Equal(HighScoreAxefighting, HighscoreCategoryFromString("axefighting"))

	assert.Equal(HighScoreCharmpoints, HighscoreCategoryFromString("charm"))
	assert.Equal(HighScoreCharmpoints, HighscoreCategoryFromString("charms"))
	assert.Equal(HighScoreCharmpoints, HighscoreCategoryFromString("charmpoints"))
	assert.Equal(HighScoreCharmpoints, HighscoreCategoryFromString("charmspoints"))

	assert.Equal(HighScoreClubfighting, HighscoreCategoryFromString("club"))
	assert.Equal(HighScoreClubfighting, HighscoreCategoryFromString("clubfighting"))

	assert.Equal(HighScoreDistancefighting, HighscoreCategoryFromString("distance"))
	assert.Equal(HighScoreDistancefighting, HighscoreCategoryFromString("distancefighting"))

	assert.Equal(HighScoreFishing, HighscoreCategoryFromString("fishing"))

	assert.Equal(HighScoreFistfighting, HighscoreCategoryFromString("fist"))
	assert.Equal(HighScoreFistfighting, HighscoreCategoryFromString("fistfighting"))

	assert.Equal(HighScoreGoshnarstaint, HighscoreCategoryFromString("goshnar"))
	assert.Equal(HighScoreGoshnarstaint, HighscoreCategoryFromString("goshnars"))
	assert.Equal(HighScoreGoshnarstaint, HighscoreCategoryFromString("goshnarstaint"))

	assert.Equal(HighScoreLoyaltypoints, HighscoreCategoryFromString("loyalty"))
	assert.Equal(HighScoreLoyaltypoints, HighscoreCategoryFromString("loyaltypoints"))

	assert.Equal(HighScoreMagiclevel, HighscoreCategoryFromString("magic"))
	assert.Equal(HighScoreMagiclevel, HighscoreCategoryFromString("mlvl"))
	assert.Equal(HighScoreMagiclevel, HighscoreCategoryFromString("magiclevel"))

	assert.Equal(HighScoreShielding, HighscoreCategoryFromString("shielding"))
	assert.Equal(HighScoreShielding, HighscoreCategoryFromString("shield"))

	assert.Equal(HighScoreSwordfighting, HighscoreCategoryFromString("sword"))
	assert.Equal(HighScoreSwordfighting, HighscoreCategoryFromString("swordfighting"))

	assert.Equal(HighScoreDromescore, HighscoreCategoryFromString("drome"))
	assert.Equal(HighScoreDromescore, HighscoreCategoryFromString("dromescore"))
}
