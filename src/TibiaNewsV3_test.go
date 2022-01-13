package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewsById(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/news/archive/6529.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	newsArticleJson := TibiaNewsV3Impl(6529, "https://www.tibia.com/news/?subtopic=newsarchive&id=6529", string(data))
	assert := assert.New(t)

	assert.Equal(6529, newsArticleJson.News.ID)
	assert.Equal("2022-01-12", newsArticleJson.News.Date)
	assert.Equal("", newsArticleJson.News.Title)
	assert.Equal("development", newsArticleJson.News.Category)
	assert.Equal("ticker", newsArticleJson.News.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6529", newsArticleJson.News.TibiaURL)
	assert.Equal("A number of issues related to the 25 years activities have been fixed, including the following: Dragon pinatas that were stuck in the inbox have been changed into dragon pinata kits. The wind-up loco can now be taken even after defeating Lord Retro. The baby seals can now be painted even when the quest The Ice Islands is active. The weight of wallpapers and fairy lights has been increased and their market category has been changed to decoration. A number of typos and map issues have been fixed as well.", newsArticleJson.News.Content)
	assert.Equal("A number of issues related to the 25 years activities have been fixed, including the following: Dragon pinatas that were stuck in the inbox have been changed into dragon pinata kits. The wind-up loco can now be taken even after defeating Lord Retro. The baby seals can now be painted even when the quest The Ice Islands is active. The weight of wallpapers and fairy lights has been increased and their market category has been changed to decoration. A number of typos and map issues have been fixed as well.", newsArticleJson.News.ContentHTML)
}
