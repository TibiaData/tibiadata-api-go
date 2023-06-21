package main

import (
	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestForumsSectionWorldboard(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/forums/section/worldboard.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boardsJson, err := TibiaForumSectionImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(90, len(boardsJson.Boards))

	adra := boardsJson.Boards[0]
	assert.Equal(146482, adra.ID)
	assert.Equal("Adra", adra.Name)
	assert.Equal("This board is for general discussions related to the game world Adra.", adra.Description)
	assert.Equal(388, adra.Threads)
	assert.Equal(1158, adra.Posts)
	assert.Equal(39395612, adra.LastPost.ID)
	assert.Equal("Rcdohl", adra.LastPost.CharacterName)
	assert.Equal("2023-06-02T23:19:00Z", adra.LastPost.PostedAt)

	alumbra := boardsJson.Boards[1]
	assert.Equal(147016, alumbra.ID)
	assert.Equal("Alumbra", alumbra.Name)
	assert.Equal("This board is for general discussions related to the game world Alumbra.", alumbra.Description)
	assert.Equal(563, alumbra.Threads)
	assert.Equal(1011, alumbra.Posts)
	assert.Equal(39395777, alumbra.LastPost.ID)
	assert.Equal("Mad Mustazza", alumbra.LastPost.CharacterName)
	assert.Equal("2023-06-04T15:51:13Z", alumbra.LastPost.PostedAt)
}
