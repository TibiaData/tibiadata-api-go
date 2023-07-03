package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTibiaCETDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T08:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 CET"))
}

func TestTibiaCESTDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T07:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 CEST"))
}

func TestTibiaUTCDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T09:52:16Z", TibiaDataDatetime("Dec 24 2021, 09:52:16 UTC"))
}

func TestTibiaForumDateFormat(t *testing.T) {
	assert.Equal(t, "2023-06-20T13:55:36Z", TibiaDataDatetime("20.06.2023 15:55:36"))
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

func TestTibiaDataGetNewsCategory(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("cipsoft", TibiaDataGetNewsCategory("newsicon_cipsoft"))
	assert.Equal("community", TibiaDataGetNewsCategory("newsicon_community"))
	assert.Equal("development", TibiaDataGetNewsCategory("newsicon_development"))
	assert.Equal("support", TibiaDataGetNewsCategory("newsicon_support"))
	assert.Equal("technical", TibiaDataGetNewsCategory("newsicon_technical"))
	assert.Equal("unknown", TibiaDataGetNewsCategory("newsicon_tibiadata"))
}

func TestTibiaDataGetNewsType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("ticker", TibiaDataGetNewsType("News Ticker"))
	assert.Equal("article", TibiaDataGetNewsType("Featured Article"))
	assert.Equal("news", TibiaDataGetNewsType("News"))
	assert.Equal("unknown", TibiaDataGetNewsType("TibiaData"))
}

func TestTibiaDataForumNameValidator(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("worldboards", TibiaDataForumNameValidator("world boards"))
	assert.Equal("worldboards", TibiaDataForumNameValidator("world"))
	assert.Equal("worldboards", TibiaDataForumNameValidator("worldboards"))

	assert.Equal("tradeboards", TibiaDataForumNameValidator("trade boards"))
	assert.Equal("tradeboards", TibiaDataForumNameValidator("trade"))
	assert.Equal("tradeboards", TibiaDataForumNameValidator("tradeboards"))

	assert.Equal("communityboards", TibiaDataForumNameValidator("community boards"))
	assert.Equal("communityboards", TibiaDataForumNameValidator("community"))
	assert.Equal("communityboards", TibiaDataForumNameValidator("communityboards"))

	assert.Equal("supportboards", TibiaDataForumNameValidator("support boards"))
	assert.Equal("supportboards", TibiaDataForumNameValidator("support"))
	assert.Equal("supportboards", TibiaDataForumNameValidator("supportboards"))

	assert.Equal("worldboards", TibiaDataForumNameValidator("def"))
}

func TestHTMLLineBreakRemover(t *testing.T) {
	const str = "a\nb\nc\nd\ne\nf\ng\nh\n"

	sanitizedStr := TibiaDataHTMLRemoveLinebreaks(str)
	assert.Equal(t, sanitizedStr, "abcdefgh")
}

func TestURLsRemover(t *testing.T) {
	const str = `<a href="https://www.tibia.com/community/?subtopic=characters&amp;name=Bobeek">Bobeek</a>`

	sanitizedStr := TibiaDataRemoveURLs(str)
	assert.Equal(t, sanitizedStr, "Bobeek")
}

func TestWorldFormater(t *testing.T) {
	const str = "hEsThDIáÛõ"

	sanitizedStr := TibiaDataStringWorldFormatToTitle(str)
	assert.Equal(t, sanitizedStr, "Hesthdiáûõ")
}

func TestEscaper(t *testing.T) {
	const (
		strOne   = "god durin"
		strTwo   = "god+durin"
		strThree = "gód"
	)

	sanitizedStrOne := TibiaDataQueryEscapeString(strOne)
	sanitizedStrTwo := TibiaDataQueryEscapeString(strTwo)
	sanitizedStrThree := TibiaDataQueryEscapeString(strThree)

	assert := assert.New(t)
	assert.Equal(sanitizedStrOne, "god+durin")
	assert.Equal(sanitizedStrTwo, "god+durin")
	assert.Equal(sanitizedStrThree, "g%F3d")
}

func TestDateParser(t *testing.T) {
	const str = "Mar 09 2022"

	sanitizedString := TibiaDataDate(str)
	assert.Equal(t, sanitizedString, "2022-03-09")
}

func TestStringToInt(t *testing.T) {
	const str = "1"

	convertedStr := TibiaDataStringToInteger(str)
	assert.Equal(t, 1, convertedStr)
}

func TestHTMLRemover(t *testing.T) {
	const str = `<div id="DeactivationContainer" onclick="ActivateWebsiteFrame();$('.LightBoxContentToHide').css('display', 'none');">abc</div>`

	sanitizedString := RemoveHtmlTag(str)
	assert.Equal(t, sanitizedString, "abc")
}

func TestFake(t *testing.T) {
	assert := assert.New(t)
	const str = "&lt;"

	htmlString := TibiaDataSanitizeEscapedString(str)
	assert.Equal(htmlString, "<")

	const strTwo = `"`

	doubleString := TibiaDataSanitizeDoubleQuoteString(strTwo)
	assert.Equal(doubleString, "'")

	const strThree = "\u00A0"

	nbspString := TibiaDataSanitizeStrings(strThree)
	assert.Equal(nbspString, " ")

	const strFour = "\u0026#39;"

	string0026 := TibiaDataSanitize0026String(strFour)
	assert.Equal(string0026, "'")

	const strFive = "1kk"

	kkInt := TibiaDataConvertValuesWithK(strFive)
	assert.Equal(kkInt, 1000000)
}
