package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaGuildsOverviewV3 func
func TibiaGuildsOverviewV3(c *gin.Context) {

	// getting params from URL
	world := c.Param("world")

	// Child of Guilds
	type Guild struct {
		Name        string `json:"name"`
		LogoURL     string `json:"logo_url"`
		Description string `json:"description"`
	}

	// Child of JSONData
	type Guilds struct {
		World     string  `json:"world"`
		Active    []Guild `json:"active"`
		Formation []Guild `json:"formation"`
	}

	//
	// The base includes two levels: Guilds and Information
	type JSONData struct {
		Guilds      Guilds      `json:"guilds"`
		Information Information `json:"information"`
	}

	// Creating empty vars
	var ActiveGuilds, FormationGuilds []Guild
	var GuildCategory string

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/community/?subtopic=guilds&world=" + TibiadataQueryEscapeStringV3(world))

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Running query over each div
	ReaderHTML.Find(".Table3 .TableContent tbody tr").Each(func(index int, s *goquery.Selection) {

		// Storing HTML into GuildsDivHTML
		GuildsDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Removing linebreaks from HTML
		GuildsDivHTML = TibiadataHTMLRemoveLinebreaksV3(GuildsDivHTML)

		if GuildCategory == "" && strings.Contains(GuildsDivHTML, "<td width=\"64\"><b>Logo</b></td>") {
			GuildCategory = "active"
		} else if GuildCategory == "active" && strings.Contains(GuildsDivHTML, "<td width=\"64\"><b>Logo</b></td>") {
			GuildCategory = "formation"
		}

		// Regex to get data for record values
		regex1 := regexp.MustCompile(`<td><img src="(.*)" width=.*\/><\/td><td><b>(.*)<\/b>(<br\/>)?(.*)<\/td><td>.*`)
		subma1 := regex1.FindAllStringSubmatch(GuildsDivHTML, -1)

		if len(subma1) > 0 {
			OneGuild := Guild{
				Name:        TibiaDataSanitizeEscapedString(subma1[0][2]),
				LogoURL:     subma1[0][1],
				Description: TibiaDataSanitizeEscapedString(strings.TrimSpace(subma1[0][4])),
			}

			// Adding OneWorld to correct category
			if GuildCategory == "active" {
				ActiveGuilds = append(ActiveGuilds, OneGuild)
			} else if GuildCategory == "formation" {
				FormationGuilds = append(FormationGuilds, OneGuild)
			}
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		Guilds{
			World:     world,
			Active:    ActiveGuilds,
			Formation: FormationGuilds,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaGuildsOverviewV3", jsonData)
}
