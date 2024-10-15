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

	creaturesJson, err := TibiaCreaturesOverviewImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := creaturesJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=creatures", information.Link)
	assert.Equal("Minotaur Amazon", creaturesJson.Creatures.Boosted.Name)
	assert.Equal("minotauramazon", creaturesJson.Creatures.Boosted.Race)
	assert.Equal("https://static.tibia.com/images/global/header/monsters/minotauramazon.gif", creaturesJson.Creatures.Boosted.ImageURL)

	assert.Equal(638, len(creaturesJson.Creatures.Creatures))

	acidblob := creaturesJson.Creatures.Creatures[0]
	assert.Equal("Acid Blobs", acidblob.Name)
	assert.Equal("acidblob", acidblob.Race)
	assert.Equal("https://static.tibia.com/images/library/acidblob.gif", acidblob.ImageURL)
	assert.False(acidblob.Featured)

	minotauramazon := creaturesJson.Creatures.Creatures[360]
	assert.Equal("Minotaur Amazons", minotauramazon.Name)
	assert.Equal("minotauramazon", minotauramazon.Race)
	assert.Equal("https://static.tibia.com/images/library/minotauramazon.gif", minotauramazon.ImageURL)
	assert.True(minotauramazon.Featured)

	quarapredator := creaturesJson.Creatures.Creatures[465]
	assert.Equal("Quara Predators", quarapredator.Name)
	assert.Equal("quarapredator", quarapredator.Race)
	assert.Equal("https://static.tibia.com/images/library/quarapredator.gif", quarapredator.ImageURL)
	assert.False(quarapredator.Featured)

	slimes := creaturesJson.Creatures.Creatures[510]
	assert.Equal("Slimes", slimes.Name)
	assert.Equal("slime", slimes.Race)
	assert.Equal("https://static.tibia.com/images/library/slime.gif", slimes.ImageURL)
	assert.False(slimes.Featured)
}
