package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFansites(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/fansites/all.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	fansitesJson := TibiaFansitesV3Impl(string(data))
	assert := assert.New(t)

	assert.Equal(17, len(fansitesJson.Fansites.PromotedFansites))

	assert.Equal(22, len(fansitesJson.Fansites.SupportedFansites))

	tibiaDataFansite := fansitesJson.Fansites.SupportedFansites[6]
	assert.Equal("TibiaData.com", tibiaDataFansite.Name)
	assert.Equal("https://static.tibia.com/images/community/fansitelogos/TibiaData.com.gif", tibiaDataFansite.LogoURL)
	assert.Equal("https://tibiadata.com", tibiaDataFansite.Homepage)
	assert.Equal("Trollefar", tibiaDataFansite.Contact)
	assert.False(tibiaDataFansite.ContentType.Statistics)
	assert.False(tibiaDataFansite.ContentType.Texts)
	assert.True(tibiaDataFansite.ContentType.Tools)
	assert.False(tibiaDataFansite.ContentType.Wiki)
	assert.False(tibiaDataFansite.SocialMedia.Discord)
	assert.False(tibiaDataFansite.SocialMedia.Facebook)
	assert.False(tibiaDataFansite.SocialMedia.Instagram)
	assert.False(tibiaDataFansite.SocialMedia.Reddit)
	assert.False(tibiaDataFansite.SocialMedia.Twitch)
	assert.False(tibiaDataFansite.SocialMedia.Twitter)
	assert.False(tibiaDataFansite.SocialMedia.Youtube)
	assert.Equal(1, len(tibiaDataFansite.Languages))
	assert.Equal("us", tibiaDataFansite.Languages[0])
	assert.Equal(3, len(tibiaDataFansite.Specials))
	assert.Equal("API for Tibia data in JSON.", tibiaDataFansite.Specials[0])
	assert.Equal("Discord server.", tibiaDataFansite.Specials[1])
	assert.Equal("GitHub participant.", tibiaDataFansite.Specials[2])
	assert.False(tibiaDataFansite.FansiteItem)
	assert.Equal("", tibiaDataFansite.FansiteItemURL)

	tibiaGalleryFansite := fansitesJson.Fansites.SupportedFansites[9]
	assert.Equal("TibiaGallery.com", tibiaGalleryFansite.Name)
	assert.Equal("https://static.tibia.com/images/community/fansitelogos/TibiaGallery.com.gif", tibiaGalleryFansite.LogoURL)
	assert.Equal("https://tibiagallery.com/", tibiaGalleryFansite.Homepage)
	assert.Equal("Ewrr", tibiaGalleryFansite.Contact)
	assert.False(tibiaGalleryFansite.ContentType.Statistics)
	assert.False(tibiaGalleryFansite.ContentType.Texts)
	assert.True(tibiaGalleryFansite.ContentType.Tools)
	assert.False(tibiaGalleryFansite.ContentType.Wiki)
	assert.False(tibiaGalleryFansite.SocialMedia.Discord)
	assert.False(tibiaGalleryFansite.SocialMedia.Facebook)
	assert.True(tibiaGalleryFansite.SocialMedia.Instagram)
	assert.False(tibiaGalleryFansite.SocialMedia.Reddit)
	assert.False(tibiaGalleryFansite.SocialMedia.Twitch)
	assert.False(tibiaGalleryFansite.SocialMedia.Twitter)
	assert.False(tibiaGalleryFansite.SocialMedia.Youtube)
	assert.Equal(9, len(tibiaGalleryFansite.Languages))
	assert.Equal("br", tibiaGalleryFansite.Languages[0])
	assert.Equal("pl", tibiaGalleryFansite.Languages[1])
	assert.Equal("mx", tibiaGalleryFansite.Languages[2])
	assert.Equal("us", tibiaGalleryFansite.Languages[3])
	assert.Equal("se", tibiaGalleryFansite.Languages[4])
	assert.Equal("de", tibiaGalleryFansite.Languages[5])
	assert.Equal("fi", tibiaGalleryFansite.Languages[6])
	assert.Equal("fr", tibiaGalleryFansite.Languages[7])
	assert.Equal("tr", tibiaGalleryFansite.Languages[8])
	assert.Equal(1, len(tibiaGalleryFansite.Specials))
	assert.Equal("Upload, browse, like and share pictures.", tibiaGalleryFansite.Specials[0])
	assert.True(tibiaGalleryFansite.FansiteItem)
	assert.Equal("https://static.tibia.com/images/community/fansiteitems/TibiaGallery.com.gif", tibiaGalleryFansite.FansiteItemURL)
}
