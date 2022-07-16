package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestOverview(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creatures.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	creaturesJson, err := TibiaCreaturesOverviewV3Impl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Quara Predator", creaturesJson.Creatures.Boosted.Name)
	assert.Equal("quarapredator", creaturesJson.Creatures.Boosted.Race)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/quarapredator.gif", creaturesJson.Creatures.Boosted.ImageURL)

	assert.Equal(553, len(creaturesJson.Creatures.Creatures))

	slimes := creaturesJson.Creatures.Creatures[444]
	assert.Equal("Slimes", slimes.Name)
	assert.Equal("slime", slimes.Race)
	assert.Equal("https://static.tibia.com/images/library/slime.gif", slimes.ImageURL)
	assert.False(slimes.Featured)

	quarapredator := creaturesJson.Creatures.Creatures[407]
	assert.Equal("Quara Predators", quarapredator.Name)
	assert.Equal("quarapredator", quarapredator.Race)
	assert.Equal("https://static.tibia.com/images/library/quarapredator.gif", quarapredator.ImageURL)
	assert.True(quarapredator.Featured)
}
