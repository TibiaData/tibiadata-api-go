package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaWorldsWorldV3 func
func TibiaWorldsWorldV3(c *gin.Context) {

	// getting params from URL
	world := c.Param("world")

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
	type JSONData struct {
		Worlds      Worlds      `json:"worlds"`
		Information Information `json:"information"`
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=worlds&world=" + TibiadataQueryEscapeStringV3(world)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaWorldsWorldV3", gin.H{"error": err.Error()})
		return
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty vars
	var WorldsStatus, WorldsRecordDate, WorldsCreationDate, WorldsLocation, WorldsPvpType, WorldsTransferType, WorldsBattleyeDate, WorldsGameWorldType, WorldsTournamentWorldType string
	var WorldsQuestTitles []string
	var WorldsPlayersOnline, WorldsRecordPlayers int
	var WorldsPremiumOnly, WorldsBattleyeProtected bool
	var WorldsOnlinePlayers []OnlinePlayers

	// Running query over each div
	ReaderHTML.Find(".Table1 .InnerTableContainer table tr").Each(func(index int, s *goquery.Selection) {

		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for record values
		regex1 := regexp.MustCompile(`<td class=.*>(.*):<\/td><td>(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

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
				WorldsPlayersOnline = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			}
			if WorldsInformationLeftColumn == "Online Record" {
				// Regex to get data for record values
				regex2 := regexp.MustCompile(`(.*) players \(on (.*)\)`)
				subma2 := regex2.FindAllStringSubmatch(WorldsInformationRightColumn, -1)

				if len(subma2) > 0 {
					// setting record values
					WorldsRecordPlayers = TibiadataStringToIntegerV3(subma2[0][1])
					WorldsRecordDate = TibiadataDatetimeV3(subma2[0][2])
				}
			}
			if WorldsInformationLeftColumn == "Creation Date" {
				WorldsCreationDate = WorldsInformationRightColumn
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
							WorldsQuestTitles = append(WorldsQuestTitles, TibiadataRemoveURLsV3(str))
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
						regex21 := regexp.MustCompile(`Protected by BattlEye since (.*)\.`)
						subma21 := regex21.FindAllStringSubmatch(WorldsInformationRightColumn, -1)
						WorldsBattleyeDate = subma21[0][1]
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

	})

	// Running query over each div
	ReaderHTML.Find(".Table2 .InnerTableContainer table tr").First().NextAll().Each(func(index int, s *goquery.Selection) {

		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for record values
		regex1 := regexp.MustCompile(`<td style=.*name=.*">(.*)<\/a>.*">(.*)<\/td>.*">(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

		if len(subma1) > 0 {

			WorldsOnlinePlayers = append(WorldsOnlinePlayers, OnlinePlayers{
				Name:     subma1[0][1],
				Level:    TibiadataStringToIntegerV3(subma1[0][2]),
				Vocation: subma1[0][3],
			})
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
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
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaWorldsWorldV3", jsonData)
}
