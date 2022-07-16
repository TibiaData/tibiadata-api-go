package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TibiaData/tibiadata-api-go/src/validation"
)

// Child of Guild
type Guildhall struct {
	Name  string `json:"name"`
	World string `json:"world"` // Maybe duplicate info? Guild can only be on one world..
	/*
		Town      string `json:"town"`       // We can collect that from cached info?
		Status    string `json:"status"`     // rented (but maybe also auctioned)
		Owner     string `json:"owner"`      // We can collect that from cached info?
		HouseID   int    `json:"houseid"`    // We can collect that from cached info?
	*/
	PaidUntil string `json:"paid_until"` // Paid until date
}

// Child of Guild
type GuildMember struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Rank     string `json:"rank"`
	Vocation string `json:"vocation"`
	Level    int    `json:"level"`
	Joined   string `json:"joined"`
	Status   string `json:"status"`
}

// Child of Guild
type InvitedGuildMember struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// Child of Guilds
type Guild struct {
	Name               string               `json:"name"`
	World              string               `json:"world"`
	LogoURL            string               `json:"logo_url"`
	Description        string               `json:"description"`
	Guildhalls         []Guildhall          `json:"guildhalls"`
	Active             bool                 `json:"active"`
	Founded            string               `json:"founded"`
	Applications       bool                 `json:"open_applications"`
	Homepage           string               `json:"homepage"`
	InWar              bool                 `json:"in_war"`
	DisbandedDate      string               `json:"disband_date"`
	DisbandedCondition string               `json:"disband_condition"`
	PlayersOnline      int                  `json:"players_online"`
	PlayersOffline     int                  `json:"players_offline"`
	MembersTotal       int                  `json:"members_total"`
	MembersInvited     int                  `json:"members_invited"`
	Members            []GuildMember        `json:"members"`
	Invited            []InvitedGuildMember `json:"invites"`
}

// Child of JSONData
type Guilds struct {
	Guild Guild `json:"guild"`
}

//
// The base includes two levels: Guild and Information
type GuildResponse struct {
	Guilds      Guilds      `json:"guilds"`
	Information Information `json:"information"`
}

var (
	GuildLogoRegex                     = regexp.MustCompile(`.*img src="(.*)" width=.*`)
	GuildWorldAndFoundationRegex       = regexp.MustCompile(`The guild was founded on (.*) on (.*).<br/>`)
	GuildHomepageRegex                 = regexp.MustCompile(`<a href="(.*)" target=.*>`)
	GuildhallRegex                     = regexp.MustCompile(` is (.*). The rent is paid until (.*).<br/>`)
	GuildDisbaneRegex                  = regexp.MustCompile(`<b>It will be disbanded on (.*.[0-9]+.[0-9]+) (.*)\.<\/b>.*`)
	GuildMemberInformationRegex        = regexp.MustCompile(`<td>(.*)<\/td><td><a.*">(.*)<\/a>(.*)<\/td><td>(.*)<\/td><td>([0-9]+)<\/td><td>(.*)<\/td><td class.*class.*">(.*)<\/span><\/td>`)
	GuildMemberInvitesInformationRegex = regexp.MustCompile(`<td><a.*">(.*)<\/a><\/td><td>(.*)<\/td>`)
)

func TibiaGuildsGuildV3Impl(guild string, BoxContentHTML string) (*GuildResponse, error) {
	// Creating empty vars
	var (
		MembersData                                                                                                                                                    []GuildMember
		InvitedData                                                                                                                                                    []InvitedGuildMember
		GuildGuildhallData                                                                                                                                             []Guildhall
		MembersRank, MembersTitle, MembersStatus, GuildDescription, GuildDisbandedDate, GuildDisbandedCondition, GuildHomepage, GuildWorld, GuildLogoURL, GuildFounded string
		GuildActive, GuildApplications, GuildInWar, GuildDescriptionFinished                                                                                           bool
		MembersCountOnline, MembersCountOffline, MembersCountInvited                                                                                                   int
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaGuildsGuildV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	guildName, err := ReaderHTML.Find("h1").First().Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaGuildsGuildV3Impl failed at ReaderHTML.Find, err: %s", err)
	}

	if guildName == "" {
		return nil, validation.ErrorGuildNotFound
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPA, err := ReaderHTML.Find(".BoxContent table").Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaGuildsGuildV3Impl failed at InnerTableContainerTMPA ReaderHTML.Find, err: %s", err)
	}

	subma1b := GuildLogoRegex.FindAllStringSubmatch(InnerTableContainerTMPA, -1)

	if len(subma1b) > 0 {
		GuildLogoURL = subma1b[0][1]
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find("#GuildInformationContainer").Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaGuildsGuildV3Impl failed at InnerTableContainerTMPB ReaderHTML.Find, err: %s", err)
	}

	for _, line := range strings.Split(strings.TrimSuffix(InnerTableContainerTMPB, "\n"), "\n") {
		// Guild information
		if !GuildDescriptionFinished {
			// First line is the description..
			GuildDescription += strings.ReplaceAll(line+"\n", "<br/><br/>\n", "")

			// Abort loop and continue wiht next section
			if strings.Contains(line, "<br/><br/>") {
				GuildDescription = strings.TrimSpace(TibiaDataSanitizeEscapedString(GuildDescription))
				GuildDescriptionFinished = true
			}

		}
		
		if GuildDescriptionFinished || strings.HasPrefix(line, "The guild was founded on ") {
			// The rest of the Guild information

			if strings.HasPrefix(GuildDescription, "The guild was founded on ") {
				GuildDescription = ""
				GuildDescriptionFinished = true
			}

			if strings.Contains(line, "The guild was founded on") {
				// Regex to get GuildWorld and GuildFounded
				subma1b := GuildWorldAndFoundationRegex.FindAllStringSubmatch(line, -1)
				GuildWorld = subma1b[0][1]
				GuildFounded = TibiaDataDateV3(subma1b[0][2])
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
				subma1c := GuildHomepageRegex.FindAllStringSubmatch(line, -1)
				GuildHomepage = subma1c[0][1]
			}

			// If guildhall
			if strings.Contains(line, "Their home on "+GuildWorld) {
				subma1b := GuildhallRegex.FindAllStringSubmatch(line, -1)

				GuildGuildhallData = append(GuildGuildhallData, Guildhall{
					Name:      TibiaDataSanitizeEscapedString(subma1b[0][1]),
					World:     GuildWorld,
					PaidUntil: TibiaDataDateV3(subma1b[0][2]),
				})
			}

			// If disbanded
			if strings.Contains(line, "<b>It will be disbanded on ") {
				subma1c := GuildDisbaneRegex.FindAllStringSubmatch(line, -1)
				if len(subma1c) > 0 {
					GuildDisbandedDate = subma1c[0][1]
					GuildDisbandedCondition = subma1c[0][2]
				}
			}
		}
	}

	var insideError error

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into GuildsDivHTML
		GuildsDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error]  TibiaGuildsGuildV3Impl failed at GuildsDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Removing linebreaks from HTML
		GuildsDivHTML = TibiaDataHTMLRemoveLinebreaksV3(GuildsDivHTML)

		// Regex to get data for record values
		subma1 := GuildMemberInformationRegex.FindAllStringSubmatch(GuildsDivHTML, -1)

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

			MembersData = append(MembersData, GuildMember{
				Name:     TibiaDataSanitizeStrings(subma1[0][2]),
				Title:    MembersTitle,
				Rank:     MembersRank,
				Vocation: subma1[0][4],
				Level:    TibiaDataStringToIntegerV3(subma1[0][5]),
				Joined:   TibiaDataDateV3(subma1[0][6]),
				Status:   MembersStatus,
			})
		} else {
			// Regex to get data for record values
			subma2 := GuildMemberInvitesInformationRegex.FindAllStringSubmatch(GuildsDivHTML, -1)

			if len(subma2) > 0 {
				MembersCountInvited++
				InvitedData = append(InvitedData, InvitedGuildMember{
					Name: TibiaDataSanitizeStrings(subma2[0][1]),
					Date: subma2[0][2],
				})
			}
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &GuildResponse{
		Guilds{
			Guild{
				Name:               guildName,
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
				PlayersOnline:      MembersCountOnline,
				PlayersOffline:     MembersCountOffline,
				MembersTotal:       (MembersCountOnline + MembersCountOffline),
				MembersInvited:     MembersCountInvited,
				Members:            MembersData,
				Invited:            InvitedData,
			},
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
