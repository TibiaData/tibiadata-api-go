package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestBoostableBossesOverview(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boostableBossesJson, _ := TibiaBoostableBossesOverviewV3Impl(string(data))
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
