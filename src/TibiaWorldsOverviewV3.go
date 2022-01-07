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

// TibiaWorldsOverviewV3 func
func TibiaWorldsOverviewV3(c *gin.Context) {
	// The base includes two levels: Worlds and Information
	type JSONData struct {
		Worlds      structs.Worlds      `json:"worlds"`
		Information structs.Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=worlds"
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaWorldsOverviewV3", gin.H{"error": err.Error()})
		return
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty vars
	var (
		RegularWorldsData, TournamentWorldsData                                                                                                     []structs.World
		WorldsRecordDate, WorldsWorldCategory, WorldsBattleyeDate, WorldsTransferType, WorldsTournamentWorldType, WorldsGameWorldType, WorldsStatus string
		WorldsRecordPlayers, WorldsAllOnlinePlayers                                                                                                 int
		WorldsPremiumOnly, WorldsBattleyeProtected                                                                                                  bool
	)

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CreatureDivHTML
		WorldsDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for record values
		regex1 := regexp.MustCompile(`.*<\/b>...(.*) players \(on (.*)\)`)
		subma1 := regex1.FindAllStringSubmatch(WorldsDivHTML, -1)

		if len(subma1) > 0 {
			// setting record values
			WorldsRecordPlayers = TibiadataStringToIntegerV3(subma1[0][1])
			WorldsRecordDate = TibiadataDatetimeV3(subma1[0][2])
		}

		if strings.Contains(WorldsDivHTML, ">Regular Worlds<") {
			WorldsWorldCategory = "regular"
		} else if strings.Contains(WorldsDivHTML, ">Tournament Worlds<") {
			WorldsWorldCategory = "tournament"
		}

		// Regex to get data for name, race and img src param for creature
		regex2 := regexp.MustCompile(`.*world=.*">(.*)<\/a><\/td>.*right;">(.*)<\/td><td>(.*)<\/td><td>(.*)<\/td><td align="center" valign="middle">(.*)<\/td><td>(.*)<\/td>`)
		subma2 := regex2.FindAllStringSubmatch(WorldsDivHTML, -1)

		// check if regex return length is over 0
		if len(subma2) > 0 {
			// Creating better to use vars
			WorldsPlayersOnline := TibiadataStringToIntegerV3(subma2[0][2])
			WorldsBattlEye := subma2[0][5]
			WorldsAdditionalInfo := subma2[0][6]

			// Setting the players_online & overall players_online
			WorldsAllOnlinePlayers += WorldsPlayersOnline
			switch {
			case WorldsPlayersOnline > 0:
				WorldsStatus = "online"
			case subma2[0][2] == "-":
				WorldsStatus = "unknown"
			default:
				WorldsStatus = "offline"
			}

			// Setting the premium_only
			if strings.Contains(WorldsAdditionalInfo, "premium") {
				WorldsPremiumOnly = true
			} else {
				WorldsPremiumOnly = false
			}

			// Setting the transfer_type
			WorldsTransferType = "regular"
			if strings.Contains(WorldsAdditionalInfo, "blocked") {
				WorldsTransferType = "blocked"
			} else if strings.Contains(WorldsAdditionalInfo, "locked") {
				WorldsTransferType = "locked"
			}

			// Setting the game_world_type
			WorldsGameWorldType = "regular"
			if strings.Contains(WorldsAdditionalInfo, "experimental") {
				WorldsGameWorldType = "experimental"
			} else if WorldsWorldCategory == "tournament" {
				WorldsGameWorldType = "tournament"
			}

			// Determine Battleye Protection
			if len(WorldsBattlEye) > 0 {
				WorldsBattleyeProtected = true
				if strings.Contains(WorldsBattlEye, "BattlEye since its release") {
					WorldsBattleyeDate = "release"
				} else {
					regex21 := regexp.MustCompile(`.*game world has been protected by BattlEye since (.*).&lt;\/p.*`)
					subma21 := regex21.FindAllStringSubmatch(WorldsBattlEye, -1)
					WorldsBattleyeDate = subma21[0][1]
				}
			} else {
				// This world is without protection..
				WorldsBattleyeProtected = false
				WorldsBattleyeDate = ""
			}

			// Setting the tournament_world_type param
			switch WorldsWorldCategory {
			case "regular":
				WorldsTournamentWorldType = ""
			case "tournament":
				WorldsGameWorldType = "tournament"
				WorldsTournamentWorldType = "regular"
				if strings.Contains(WorldsAdditionalInfo, "restricted") {
					WorldsTournamentWorldType = "restricted"
				}
			}

			// Creating data block to return
			OneWorld := structs.World{
				Name:                subma2[0][1],
				Status:              WorldsStatus,
				PlayersOnline:       WorldsPlayersOnline,
				Location:            subma2[0][3],
				PvpType:             subma2[0][4],
				PremiumOnly:         WorldsPremiumOnly,
				TransferType:        WorldsTransferType,
				BattleyeProtected:   WorldsBattleyeProtected,
				BattleyeDate:        WorldsBattleyeDate,
				GameWorldType:       WorldsGameWorldType,
				TournamentWorldType: WorldsTournamentWorldType,
			}

			// Adding OneWorld to correct category
			switch WorldsWorldCategory {
			case "regular":
				RegularWorldsData = append(RegularWorldsData, OneWorld)
			case "tournament":
				TournamentWorldsData = append(TournamentWorldsData, OneWorld)
			}
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		structs.Worlds{
			PlayersOnline:    WorldsAllOnlinePlayers,
			RecordPlayers:    WorldsRecordPlayers,
			RecordDate:       WorldsRecordDate,
			RegularWorlds:    RegularWorldsData,
			TournamentWorlds: TournamentWorldsData,
		},
		structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaWorldsOverviewV3", jsonData)
}
