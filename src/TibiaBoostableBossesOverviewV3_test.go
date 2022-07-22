package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoostableBossesOverview(t *testing.T) {
	data, err := os.ReadFile("../testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	boostableBossesJson := TibiaBoostableBossesOverviewV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal("The Pale Worm", boostableBossesJson.BoostableBosses.Boosted.Name)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/paleworm.gif", boostableBossesJson.BoostableBosses.Boosted.ImageURL)

	assert.Equal(88, len(boostableBossesJson.BoostableBosses.BoostableBosses))

	gnomevil := boostableBossesJson.BoostableBosses.BoostableBosses[18]
	assert.Equal("Gnomevil", gnomevil.Name)
	assert.Equal("https://static.tibia.com/images/library/gnomehorticulist.gif", gnomevil.ImageURL)
	assert.False(gnomevil.Featured)

	paleworm := boostableBossesJson.BoostableBosses.BoostableBosses[72]
	assert.Equal("The Pale Worm", paleworm.Name)
	assert.Equal("https://static.tibia.com/images/library/paleworm.gif", paleworm.ImageURL)
	assert.True(paleworm.Featured)
}
