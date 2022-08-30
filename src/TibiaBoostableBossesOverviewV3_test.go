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

	assert.Equal("Goshnar's Malice", boostableBossesJson.BoostableBosses.Boosted.Name)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/goshnarsmalice.gif", boostableBossesJson.BoostableBosses.Boosted.ImageURL)

	assert.Equal(89, len(boostableBossesJson.BoostableBosses.BoostableBosses))

	gnomevil := boostableBossesJson.BoostableBosses.BoostableBosses[18]
	assert.Equal("Gnomevil", gnomevil.Name)
	assert.Equal("https://static.tibia.com/images/library/gnomehorticulist.gif", gnomevil.ImageURL)
	assert.False(gnomevil.Featured)

	goshnarsmalice := boostableBossesJson.BoostableBosses.BoostableBosses[23]
	assert.Equal("Goshnar's Malice", goshnarsmalice.Name)
	assert.Equal("https://static.tibia.com/images/library/goshnarsmalice.gif", goshnarsmalice.ImageURL)
	assert.True(goshnarsmalice.Featured)

	paleworm := boostableBossesJson.BoostableBosses.BoostableBosses[73]
	assert.Equal("The Pale Worm", paleworm.Name)
	assert.Equal("https://static.tibia.com/images/library/paleworm.gif", paleworm.ImageURL)
	assert.False(paleworm.Featured)
}
