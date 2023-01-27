package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of Fansite
type ContentType struct {
	Statistics bool `json:"statistics"`
	Texts      bool `json:"texts"`
	Tools      bool `json:"tools"`
	Wiki       bool `json:"wiki"`
}

// Child of Fansite
type SocialMedia struct {
	Discord   bool `json:"discord"`
	Facebook  bool `json:"facebook"`
	Instagram bool `json:"instagram"`
	Reddit    bool `json:"reddit"`
	Twitch    bool `json:"twitch"`
	Twitter   bool `json:"twitter"`
	Youtube   bool `json:"youtube"`
}

// Child of Fansites
type Fansite struct {
	Name           string      `json:"name"`
	LogoURL        string      `json:"logo_url"`
	Homepage       string      `json:"homepage"`
	Contact        string      `json:"contact"`
	ContentType    ContentType `json:"content_type"`
	SocialMedia    SocialMedia `json:"social_media"`
	Languages      []string    `json:"languages"`
	Specials       []string    `json:"specials"`
	FansiteItem    bool        `json:"fansite_item"`
	FansiteItemURL string      `json:"fansite_item_url"`
}

// Child of JSONData
type Fansites struct {
	PromotedFansites  []Fansite `json:"promoted"`
	SupportedFansites []Fansite `json:"supported"`
}

//
// The base includes two levels: Fansites and Information
type FansitesResponse struct {
	Fansites    Fansites    `json:"fansites"`
	Information Information `json:"information"`
}

var (
	FansiteInformationRegex = regexp.MustCompile(`<td><a href="(.*)" target.*img .*src="(.*)" alt="(.*)"\/><\/a>.*<a href=".*">(.*)<\/a><\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*<ul><li>(.*)<\/li><\/ul><\/td><td.*top;">(.*)<\/td>`)
	FansiteImgTagRegex      = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
	FansiteLanguagesRegex   = regexp.MustCompile(`id="Language_([a-z]{2})`)
	FansiteAnchorRegex      = regexp.MustCompile(`.*src="(.*)" alt=".*`)
)

func TibiaFansitesV3Impl(BoxContentHTML string) (*FansitesResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaFansitesV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty PromotedFansitesData and SupportedFansitesData var
	var PromotedFansitesData, SupportedFansitesData []Fansite

	// list of different fansite types
	FansiteTypes := []string{"promoted", "supported"}
	// running over the FansiteTypes array
	for _, FansiteType := range FansiteTypes {
		fansites, err := makeFansiteRequest(FansiteType, ReaderHTML)
		if err != nil {
			return nil, fmt.Errorf("[error] TibiaFansitesV3Impl failed at makeFansiteRequest, type: %s, err: %s", FansiteType, err)
		}

		switch FansiteType {
		case "promoted":
			PromotedFansitesData = fansites
		case "supported":
			SupportedFansitesData = fansites
		}
	}

	// Build the data-blob
	return &FansitesResponse{
		Fansites{
			PromotedFansites:  PromotedFansitesData,
			SupportedFansites: SupportedFansitesData,
		},
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}

func makeFansiteRequest(FansiteType string, ReaderHTML *goquery.Document) ([]Fansite, error) {
	var output []Fansite
	var insideError error

	// Running query over each tr in <FansiteType>fansitesinnertable
	ReaderHTML.Find("#" + FansiteType + "fansitesinnertable tr").First().NextAll().EachWithBreak(func(index int, s *goquery.Selection) bool {
		// #promotedfansitesinnertable
		// #supportedfansitesinnertable

		// Storing HTML into FansiteTrHTML
		FansiteTrHTML, err := s.Html()
		if err != nil {
			insideError = err
			return false
		}

		// Removing line breaks
		FansiteTrHTML = TibiaDataHTMLRemoveLinebreaksV3(FansiteTrHTML)

		// Regex to get data for fansites
		subma1 := FansiteInformationRegex.FindAllStringSubmatch(FansiteTrHTML, -1)

		if len(subma1) > 0 {
			// ContentType
			ContentTypeData := ContentType{}
			imgs1 := FansiteImgTagRegex.FindAllStringSubmatch(subma1[0][5], -1)
			out := make([]string, len(imgs1))
			for i := range out {
				s := imgs1[i][1]
				switch {
				case strings.Contains(s, "Statistics"):
					ContentTypeData.Statistics = true
				case strings.Contains(s, "ArticlesNews"):
					ContentTypeData.Texts = true
				case strings.Contains(s, "Tools"):
					ContentTypeData.Tools = true
				case strings.Contains(s, "Wiki"):
					ContentTypeData.Wiki = true
				}
			}

			// SocialMedia
			SocialMediaData := SocialMedia{}
			imgs2 := FansiteImgTagRegex.FindAllStringSubmatch(subma1[0][6], -1)
			out2 := make([]string, len(imgs2))
			for i := range out2 {
				s := imgs2[i][1]
				switch {
				case strings.Contains(s, "Discord"):
					SocialMediaData.Discord = true
				case strings.Contains(s, "Facebook"):
					SocialMediaData.Facebook = true
				case strings.Contains(s, "Instagram"):
					SocialMediaData.Instagram = true
				case strings.Contains(s, "Reddit"):
					SocialMediaData.Reddit = true
				case strings.Contains(s, "Twitch"):
					SocialMediaData.Twitch = true
				case strings.Contains(s, "Twitter"):
					SocialMediaData.Twitter = true
				case strings.Contains(s, "Youtube"):
					SocialMediaData.Youtube = true
				}
			}

			// Languages
			found := FansiteLanguagesRegex.FindAllString(subma1[0][7], -1)
			FansiteLanguagesData := make([]string, len(found))
			for i := range FansiteLanguagesData {
				FansiteLanguagesData[i] = strings.ReplaceAll(found[i], "id=\"Language_", "")
			}

			// Specials
			subma1[0][8] = TibiaDataSanitizeEscapedString(subma1[0][8])
			FansiteSpecialsData := strings.Split(subma1[0][8], "</li><li>")

			// FansiteItem & FansiteItemURL
			var FansiteItemData bool
			var FansiteItemURLData string
			subma1item := FansiteAnchorRegex.FindAllStringSubmatch(subma1[0][9], -1)
			if len(subma1item) > 0 {
				FansiteItemData = true
				FansiteItemURLData = subma1item[0][1]
			} else {
				FansiteItemData = false
				FansiteItemURLData = ""
			}

			output = append(output, Fansite{
				Name:           subma1[0][3],
				LogoURL:        subma1[0][2],
				Homepage:       subma1[0][1],
				Contact:        subma1[0][4],
				ContentType:    ContentTypeData,
				SocialMedia:    SocialMediaData,
				Languages:      FansiteLanguagesData,
				Specials:       FansiteSpecialsData,
				FansiteItem:    FansiteItemData,
				FansiteItemURL: FansiteItemURLData,
			})
		}

		return true
	})

	return output, insideError
}
