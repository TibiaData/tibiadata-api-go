package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Child of Guilds
type OverviewGuild struct {
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	Description string `json:"description"`
}

// Child of JSONData
type OverviewGuilds struct {
	World     string          `json:"world"`
	Active    []OverviewGuild `json:"active"`
	Formation []OverviewGuild `json:"formation"`
}

//
// The base includes two levels: Guilds and Information
type GuildsOverviewResponse struct {
	Guilds      OverviewGuilds `json:"guilds"`
	Information Information    `json:"information"`
}

// TibiaGuildsOverviewV3 func
func TibiaGuildsOverviewV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=guilds&world=" + TibiadataQueryEscapeStringV3(world)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaGuildsOverviewV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaGuildsOverviewV3Impl(world, BoxContentHTML)

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaGuildsOverviewV3", jsonData)
}

func TibiaGuildsOverviewV3Impl(world string, BoxContentHTML string) GuildsOverviewResponse {
	// Creating empty vars
	var (
		ActiveGuilds, FormationGuilds []OverviewGuild
		GuildCategory                 string
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Running query over each div
	ReaderHTML.Find(".TableContainer").Each(func(index int, s *goquery.Selection) {

		// Figure out the guild category
		s.Find(".Text").Each(func(index int, s *goquery.Selection) {
			tableName := s.Nodes[0].FirstChild.Data
			if strings.Contains(tableName, "Active Guilds") {
				GuildCategory = "active"

			} else if strings.Contains(tableName, "Guilds in Course of Formation") {
				GuildCategory = "formation"
			}
		})

		if GuildCategory != "" {
			// Extract guilds
			s.Find(".TableContent tbody").Children().NextAll().Each(func(index int, s *goquery.Selection) {
				tableRow := s.Nodes[0]
				nameAndDescriptionNode := tableRow.FirstChild.NextSibling.NextSibling

				name := nameAndDescriptionNode.FirstChild.FirstChild.Data
				logoURL := tableRow.FirstChild.FirstChild.Attr[0].Val
				description := ""

				// Check if there's a description to fetch.
				if nameAndDescriptionNode.FirstChild.NextSibling != nil && nameAndDescriptionNode.FirstChild.NextSibling.NextSibling != nil {
					description = strings.TrimSpace(nameAndDescriptionNode.FirstChild.NextSibling.NextSibling.Data)
				}

				OneGuild := OverviewGuild{
					Name:        name,
					LogoURL:     logoURL,
					Description: description,
				}

				// Adding OneGuild to correct category
				if GuildCategory == "active" {
					ActiveGuilds = append(ActiveGuilds, OneGuild)
				} else if GuildCategory == "formation" {
					FormationGuilds = append(FormationGuilds, OneGuild)
				}
			})
		}
	})

	//
	// Build the data-blob
	return GuildsOverviewResponse{
		OverviewGuilds{
			World:     world,
			Active:    ActiveGuilds,
			Formation: FormationGuilds,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
