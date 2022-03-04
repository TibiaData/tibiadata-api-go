package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewsList(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/news/newslist.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	TibiaDataHost = "unittest.example.com"

	newsListJson := TibiaNewslistV3Impl(90, string(data))
	assert := assert.New(t)

	assert.Equal(50, len(newsListJson.News))

	firstArticle := newsListJson.News[0]
	assert.Equal(6529, firstArticle.ID)
	assert.Equal("2022-01-12", firstArticle.Date)
	assert.Equal("A number of issues related to the 25 years activities have been fixed,...", firstArticle.News)
	assert.Equal("development", firstArticle.Category)
	assert.Equal("ticker", firstArticle.Type)
	assert.Equal("https://www.tibia.com/news/?subtopic=newsarchive&id=6529", firstArticle.TibiaURL)
	assert.Equal("https://unittest.example.com/v3/news/id/6529", firstArticle.ApiURL)
}
