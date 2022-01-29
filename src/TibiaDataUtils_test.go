package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTibiaCETDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T08:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 CET"))
}

func TestTibiaCESTDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T07:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 CEST"))
}

func TestTibiaUTCDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T09:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 UTC"))
}

func TestEnvFunctions(t *testing.T) {
	assert := assert.New(t)

	assert.False(false, isEnvExist("test"))

	assert.Equal("default", getEnv("TIBIADATA_ENV", "default"))

	assert.False(false, getEnvAsBool("TIBIADATA_ENV", true))
}

func TestTibiaDataVocationValidator(t *testing.T) {
	assert := assert.New(t)

	var (
		x, y string
	)

	x, y = TibiaDataVocationValidator("none")
	assert.Equal(x, "none")
	assert.Equal(y, "1")
	x, y = TibiaDataVocationValidator("knight")
	assert.Equal(x, "knights")
	assert.Equal(y, "2")
	x, y = TibiaDataVocationValidator("paladin")
	assert.Equal(x, "paladins")
	assert.Equal(y, "3")
	x, y = TibiaDataVocationValidator("sorcerer")
	assert.Equal(x, "sorcerers")
	assert.Equal(y, "4")
	x, y = TibiaDataVocationValidator("druid")
	assert.Equal(x, "druids")
	assert.Equal(y, "5")
}

func TestTibiadataGetNewsCategory(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("cipsoft", TibiadataGetNewsCategory("newsicon_cipsoft"))
	assert.Equal("community", TibiadataGetNewsCategory("newsicon_community"))
	assert.Equal("development", TibiadataGetNewsCategory("newsicon_development"))
	assert.Equal("support", TibiadataGetNewsCategory("newsicon_support"))
	assert.Equal("technical", TibiadataGetNewsCategory("newsicon_technical"))
	assert.Equal("unknown", TibiadataGetNewsCategory("newsicon_tibiadata"))
}

func TestTibiadataGetNewsType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("ticker", TibiadataGetNewsType("News Ticker"))
	assert.Equal("article", TibiadataGetNewsType("Featured Article"))
	assert.Equal("news", TibiadataGetNewsType("News"))
	assert.Equal("unknown", TibiadataGetNewsType("TibiaData"))
}
