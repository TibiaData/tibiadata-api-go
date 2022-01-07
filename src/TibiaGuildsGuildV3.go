package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"tibiadata-api-go/src/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaGuildsGuildV3 func
func TibiaGuildsGuildV3(c *gin.Context) {
	// getting params from URL
	guild := c.Param("guild")

	// The base includes two levels: Guild and Information
	type JSONData struct {
		Guild       structs.Guild       `json:"guild"`
		Information structs.Information `json:"information"`
	}

	// Creating empty vars
	var (
		MembersData                                                                                                                                                    []structs.GuildMembers
		InvitedData                                                                                                                                                    []structs.GuildInvited
		GuildGuildhallData                                                                                                                                             []structs.Guildhall
		MembersRank, MembersTitle, MembersStatus, GuildDescription, GuildDisbandedDate, GuildDisbandedCondition, GuildHomepage, GuildWorld, GuildLogoURL, GuildFounded string
		GuildActive, GuildApplications, GuildInWar, GuildNameDetected, GuildDescriptionFinished                                                                        bool
		MembersCountOnline, MembersCountOffline, MembersCountInvited                                                                                                   int
	)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=" + TibiadataQueryEscapeStringV3(guild)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaGuildsGuildV3", gin.H{"error": err.Error()})
		return
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPA, err := ReaderHTML.Find(".BoxContent table").Html()
	if err != nil {
		log.Fatal(err)
	}
	regex1b := regexp.MustCompile(`.*img src="(.*)" width=.*`)
	subma1b := regex1b.FindAllStringSubmatch(InnerTableContainerTMPA, -1)
	GuildLogoURL = subma1b[0][1]

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find("#GuildInformationContainer").Html()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range strings.Split(strings.TrimSuffix(InnerTableContainerTMPB, "\n"), "\n") {
		// setting guild name based on html
		if !GuildNameDetected {
			guild = strings.TrimSpace(RemoveHtmlTag(line))
			GuildNameDetected = true
		}

		// Guild information
		if !GuildDescriptionFinished {
			// First line is the description..
			GuildDescription += strings.ReplaceAll(line+"\n", "<br/><br/>\n", "")

			// Abort loop and continue wiht next section
			if strings.Contains(line, "<br/><br/>") {
				GuildDescription = strings.TrimSpace(TibiaDataSanitizeEscapedString(GuildDescription))
				GuildDescriptionFinished = true
			}
		} else if GuildDescriptionFinished {
			// The rest of the Guild information

			if strings.Contains(line, "The guild was founded on") {
				// Regex to get GuildWorld and GuildFounded
				regex1b := regexp.MustCompile(`The guild was founded on (.*) on (.*).<br/>`)
				subma1b := regex1b.FindAllStringSubmatch(line, -1)
				GuildWorld = subma1b[0][1]
				GuildFounded = TibiadataDateV3(subma1b[0][2])
			}

			// If to get GuildActive
			if strings.Contains(line, "It is currently active") {
				GuildActive = true
			}

			// If open for applications
			if strings.Contains(line, "Guild is opened for applications.") {
				GuildApplications = true
			} else if strings.Contains(line, "Guild is closed for applications during war.") {
				GuildInWar = true
			}

			if strings.Contains(line, "The official homepage is") {
				regex1c := regexp.MustCompile(`<a href="(.*)" target=.*>`)
				subma1c := regex1c.FindAllStringSubmatch(line, -1)
				GuildHomepage = subma1c[0][1]
			}

			// If guildhall
			if strings.Contains(line, "Their home on "+GuildWorld) {
				// Regex to get GuildWorld and GuildFounded
				regex1b := regexp.MustCompile(`Their home on ` + GuildWorld + ` is (.*). The rent is paid until (.*).<br/>`)
				subma1b := regex1b.FindAllStringSubmatch(line, -1)

				GuildGuildhallData = append(GuildGuildhallData, structs.Guildhall{
					Name:      TibiaDataSanitizeEscapedString(subma1b[0][1]),
					World:     GuildWorld,
					PaidUntil: TibiadataDateV3(subma1b[0][2]),
				})
			}

			// If disbanded
			if strings.Contains(line, "<b>It will be disbanded on ") {
				regex1c := regexp.MustCompile(`<b>It will be disbanded on (.*.[0-9]+.[0-9]+) (.*)\.<\/b>.*`)
				subma1c := regex1c.FindAllStringSubmatch(line, -1)
				if len(subma1c) > 0 {
					GuildDisbandedDate = subma1c[0][1]
					GuildDisbandedCondition = subma1c[0][2]
				}
			}
		}
	}

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into GuildsDivHTML
		GuildsDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Removing linebreaks from HTML
		GuildsDivHTML = TibiadataHTMLRemoveLinebreaksV3(GuildsDivHTML)

		// Regex to get data for record values
		regex1 := regexp.MustCompile(`<td>(.*)<\/td><td><a.*">(.*)<\/a>(.*)<\/td><td>(.*)<\/td><td>([0-9]+)<\/td><td>(.*)<\/td><td class.*class.*">(.*)<\/span><\/td>`)
		subma1 := regex1.FindAllStringSubmatch(GuildsDivHTML, -1)

		if len(subma1) > 0 {
			// Rank name
			if len(subma1[0][1]) > 2 {
				MembersRank = subma1[0][1]
			}

			// Title
			MembersTitle = strings.ReplaceAll(strings.ReplaceAll(subma1[0][3], " (", ""), ")", "")

			// Status
			if strings.Contains(subma1[0][7], "online") {
				MembersStatus = "online"
				MembersCountOnline++
			} else {
				MembersStatus = "offline"
				MembersCountOffline++
			}

			MembersData = append(MembersData, structs.GuildMembers{
				Name:     TibiaDataSanitizeEscapedString(subma1[0][2]),
				Title:    MembersTitle,
				Rank:     MembersRank,
				Vocation: subma1[0][4],
				Level:    TibiadataStringToIntegerV3(subma1[0][5]),
				Joined:   TibiadataDateV3(subma1[0][6]),
				Status:   MembersStatus,
			})
		} else {
			// Regex to get data for record values
			regex2 := regexp.MustCompile(`<td><a.*">(.*)<\/a><\/td><td>(.*)<\/td>`)
			subma2 := regex2.FindAllStringSubmatch(GuildsDivHTML, -1)

			if len(subma2) > 0 {
				MembersCountInvited++
				InvitedData = append(InvitedData, structs.GuildInvited{
					Name: subma2[0][1],
					Date: subma2[0][2],
				})
			}
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		structs.Guild{
			Name:               guild,
			World:              GuildWorld,
			LogoURL:            GuildLogoURL,
			Description:        GuildDescription,
			Guildhalls:         GuildGuildhallData,
			Active:             GuildActive,
			Founded:            GuildFounded,
			Applications:       GuildApplications,
			Homepage:           GuildHomepage,
			InWar:              GuildInWar,
			DisbandedDate:      GuildDisbandedDate,
			DisbandedCondition: GuildDisbandedCondition,

			PlayersOnline:  MembersCountOnline,
			PlayersOffline: MembersCountOffline,
			MembersTotal:   (MembersCountOnline + MembersCountOffline),
			MembersInvited: MembersCountInvited,
			Members:        MembersData,
			Invited:        InvitedData,
		},
		structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaGuildsGuildV3", jsonData)
}
