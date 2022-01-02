package main

import (
	"html"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaFansitesV3 func
func TibiaFansitesV3(c *gin.Context) {

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
	type JSONData struct {
		Fansites    Fansites    `json:"fansites"`
		Information Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/community/?subtopic=fansites")

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty PromotedFansitesData and SupportedFansitesData var
	var PromotedFansitesData, SupportedFansitesData []Fansite

	// list of different fansite types
	FansiteTypes := []string{"promoted", "supported"}
	// running over the FansiteTypes array
	for _, FansiteType := range FansiteTypes {

		// Running query over each tr in <FansiteType>fansitesinnertable
		ReaderHTML.Find("#" + FansiteType + "fansitesinnertable tr").First().NextAll().Each(func(index int, s *goquery.Selection) {
			// #promotedfansitesinnertable
			// #supportedfansitesinnertable

			// Storing HTML into FansiteTrHTML
			FansiteTrHTML, err := s.Html()
			if err != nil {
				log.Fatal(err)
			}

			// Removing line breaks
			FansiteTrHTML = TibiadataHTMLRemoveLinebreaksV3(FansiteTrHTML)

			// Regex to get data for fansites
			regex1 := regexp.MustCompile(`<td><a href="(.*)" target.*img .*src="(.*)" alt="(.*)"\/><\/a>.*<a href=".*">(.*)<\/a><\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*<ul><li>(.*)<\/li><\/ul><\/td><td.*top;">(.*)<\/td>`)
			subma1 := regex1.FindAllStringSubmatch(FansiteTrHTML, -1)

			if len(subma1) > 0 {

				// ContentType
				ContentTypeData := ContentType{}
				var imgRE1 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
				imgs1 := imgRE1.FindAllStringSubmatch(subma1[0][5], -1)
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
				var imgRE2 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
				imgs2 := imgRE2.FindAllStringSubmatch(subma1[0][6], -1)
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
				re := regexp.MustCompile("iti__flag.iti__(..)")
				found := re.FindAllString(subma1[0][7], -1)
				FansiteLanguagesData := make([]string, len(found))
				for i := range FansiteLanguagesData {
					FansiteLanguagesData[i] = strings.ReplaceAll(found[i], "iti__flag iti__", "")
				}

				// Specials
				subma1[0][8] = html.UnescapeString(subma1[0][8])
				FansiteSpecialsData := strings.Split(subma1[0][8], "</li><li>")

				// FansiteItem & FansiteItemURL
				var FansiteItemData bool
				var FansiteItemURLData string
				regex2 := regexp.MustCompile(`.*src="(.*)" alt=".*`)
				subma1item := regex2.FindAllStringSubmatch(subma1[0][9], -1)
				if len(subma1item) > 0 {
					FansiteItemData = true
					FansiteItemURLData = subma1item[0][1]
				} else {
					FansiteItemData = false
					FansiteItemURLData = ""
				}

				switch FansiteType {
				case "promoted":
					PromotedFansitesData = append(PromotedFansitesData, Fansite{
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
				case "supported":
					SupportedFansitesData = append(SupportedFansitesData, Fansite{
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
			}

		})

	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Fansites{
			PromotedFansites:  PromotedFansitesData,
			SupportedFansites: SupportedFansitesData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaFansitesV3", jsonData)
}
