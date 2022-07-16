package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of World
type OnlinePlayers struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
}

// Child of Worlds
type World struct {
	Name                string          `json:"name"`
	Status              string          `json:"status"`                // Status:
	PlayersOnline       int             `json:"players_online"`        // Players Online:
	RecordPlayers       int             `json:"record_players"`        // Online Record:
	RecordDate          string          `json:"record_date"`           // Online Record:
	CreationDate        string          `json:"creation_date"`         // Creation Date: -> convert to YYYY-MM
	Location            string          `json:"location"`              // Location:
	PvpType             string          `json:"pvp_type"`              // PvP Type:
	PremiumOnly         bool            `json:"premium_only"`          // Premium Type: premium = true / else: false
	TransferType        string          `json:"transfer_type"`         // Transfer Type: regular (if not present) / locked / blocked
	WorldsQuestTitles   []string        `json:"world_quest_titles"`    // World Quest Titles:
	BattleyeProtected   bool            `json:"battleye_protected"`    // BattlEye Status: true if protected / false if "Not protected by BattlEye."
	BattleyeDate        string          `json:"battleye_date"`         // BattlEye Status: null if since release / else show date?
	GameWorldType       string          `json:"game_world_type"`       // Game World Type: regular / experimental / tournament (if Tournament World Type exists)
	TournamentWorldType string          `json:"tournament_world_type"` // Tournament World Type: "" (default?) / regular / restricted
	OnlinePlayers       []OnlinePlayers `json:"online_players"`
}

// Child of JSONData
type Worlds struct {
	World World `json:"world"`
}

//
// The base includes two levels: World and Information
type WorldResponse struct {
	Worlds      Worlds      `json:"worlds"`
	Information Information `json:"information"`
}

var (
	WorldDataRowRegex           = regexp.MustCompile(`<td class=.*>(.*):<\/td><td>(.*)<\/td>`)
	WorldRecordInformationRegex = regexp.MustCompile(`(.*) players \(on (.*)\)`)
	BattlEyeProtectedSinceRegex = regexp.MustCompile(`Protected by BattlEye since (.*)\.`)
	OnlinePlayerRegex           = regexp.MustCompile(`<td style=.*name=.*">(.*)<\/a>.*">(.*)<\/td>.*">(.*)<\/td>`)
)

// TibiaWorldsWorldV3 func
func TibiaWorldsWorldV3Impl(world string, BoxContentHTML string) (*WorldResponse, error) {
	//TODO: We need to read the world name from the response rather than pass it into this func

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaWorldsWorldV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty vars
	var (
		WorldsStatus, WorldsRecordDate, WorldsCreationDate, WorldsLocation, WorldsPvpType, WorldsTransferType, WorldsBattleyeDate, WorldsGameWorldType, WorldsTournamentWorldType string
		WorldsQuestTitles                                                                                                                                                         []string
		WorldsPlayersOnline, WorldsRecordPlayers                                                                                                                                  int
		WorldsPremiumOnly, WorldsBattleyeProtected                                                                                                                                bool
		WorldsOnlinePlayers                                                                                                                                                       []OnlinePlayers

		insideError error
	)

	// Running query over each div
	ReaderHTML.Find(".Table1 .InnerTableContainer table tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaWorldsWorldV3Impl failed at WorldsInformationDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		subma1 := WorldDataRowRegex.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

		if len(subma1) > 0 {
			// Creating easy to use vars (and unescape hmtl right string)
			WorldsInformationLeftColumn := subma1[0][1]
			WorldsInformationRightColumn := TibiaDataSanitizeEscapedString(subma1[0][2])

			if WorldsInformationLeftColumn == "Status" {
				switch {
				case strings.Contains(WorldsInformationRightColumn, "</div>Online"):
					WorldsStatus = "online"
				case strings.Contains(WorldsInformationRightColumn, "</div>Offline"):
					WorldsStatus = "offline"
				default:
					WorldsStatus = "unknown"
				}
			}
			if WorldsInformationLeftColumn == "Players Online" {
				WorldsPlayersOnline = TibiaDataStringToIntegerV3(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Online Record" {
				// Regex to get data for record values
				subma2 := WorldRecordInformationRegex.FindAllStringSubmatch(WorldsInformationRightColumn, -1)

				if len(subma2) > 0 {
					// setting record values
					WorldsRecordPlayers = TibiaDataStringToIntegerV3(subma2[0][1])
					WorldsRecordDate = TibiaDataDatetimeV3(subma2[0][2])
				}
			}

			if WorldsInformationLeftColumn == "Creation Date" {
				WorldsCreationDate = TibiaDataDateV3(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Location" {
				WorldsLocation = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "PvP Type" {
				WorldsPvpType = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "Premium Type" {
				WorldsPremiumOnly = true
			}

			if WorldsInformationLeftColumn == "Transfer Type" {
				WorldsTransferType = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "World Quest Titles" {
				if WorldsInformationRightColumn != "This game world currently has no title." {
					WorldsQuestTitlesTmp := strings.Split(WorldsInformationRightColumn, ", ")
					for _, str := range WorldsQuestTitlesTmp {
						if str != "" {
							WorldsQuestTitles = append(WorldsQuestTitles, TibiaDataRemoveURLsV3(str))
						}
					}
				}
			}

			if WorldsInformationLeftColumn == "BattlEye Status" {
				if WorldsInformationRightColumn == "Not protected by BattlEye." {
					WorldsBattleyeProtected = false
				} else {
					WorldsBattleyeProtected = true
					if strings.Contains(WorldsInformationRightColumn, "BattlEye since its release") {
						WorldsBattleyeDate = "release"
					} else {
						subma21 := BattlEyeProtectedSinceRegex.FindAllStringSubmatch(WorldsInformationRightColumn, -1)
						WorldsBattleyeDate = TibiaDataDateV3(subma21[0][1])
					}
				}
			}

			if WorldsInformationLeftColumn == "Game World Type" {
				WorldsGameWorldType = strings.ToLower(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Tournament World Type" {
				WorldsGameWorldType = "tournament"
				if WorldsInformationRightColumn == "Restricted Store" {
					WorldsTournamentWorldType = "restricted"
				} else {
					WorldsTournamentWorldType = strings.ToLower(WorldsInformationRightColumn)
				}
			}
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	// Running query over each div
	ReaderHTML.Find(".Table2 .InnerTableContainer table tr").First().NextAll().EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaWorldsWorldV3Impl failed at WorldsInformationDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		subma1 := OnlinePlayerRegex.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

		if len(subma1) > 0 {
			WorldsOnlinePlayers = append(WorldsOnlinePlayers, OnlinePlayers{
				Name:     TibiaDataSanitizeStrings(subma1[0][1]),
				Level:    TibiaDataStringToIntegerV3(subma1[0][2]),
				Vocation: TibiaDataSanitizeStrings(subma1[0][3]),
			})
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &WorldResponse{
		Worlds: Worlds{
			World{
				Name:                world,
				Status:              WorldsStatus,
				PlayersOnline:       WorldsPlayersOnline,
				RecordPlayers:       WorldsRecordPlayers,
				RecordDate:          WorldsRecordDate,
				CreationDate:        WorldsCreationDate,
				Location:            WorldsLocation,
				PvpType:             WorldsPvpType,
				PremiumOnly:         WorldsPremiumOnly,
				TransferType:        WorldsTransferType,
				WorldsQuestTitles:   WorldsQuestTitles,
				BattleyeProtected:   WorldsBattleyeProtected,
				BattleyeDate:        WorldsBattleyeDate,
				GameWorldType:       WorldsGameWorldType,
				TournamentWorldType: WorldsTournamentWorldType,
				OnlinePlayers:       WorldsOnlinePlayers,
			},
		},
		Information: Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
