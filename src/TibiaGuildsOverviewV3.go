package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func TibiaGuildsOverviewV3Impl(world string, BoxContentHTML string) (*GuildsOverviewResponse, error) {
	// Creating empty vars
	var (
		ActiveGuilds, FormationGuilds []OverviewGuild
		GuildCategory                 string
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaGuildsOverviewV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
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
	return &GuildsOverviewResponse{
		OverviewGuilds{
			World:     world,
			Active:    ActiveGuilds,
			Formation: FormationGuilds,
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
