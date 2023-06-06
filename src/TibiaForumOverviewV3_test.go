package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestForums(t *testing.T) {
	data, err := os.ReadFile("../testdata/forums/worldboards.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	boardsJson := TibiaForumOverviewV3Impl(string(data))
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
