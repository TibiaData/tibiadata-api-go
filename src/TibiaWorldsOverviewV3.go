package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Child of Worlds
type World struct {
	Name                string `json:"name"`
	Status              string `json:"status"`                // Online:
	PlayersOnline       int    `json:"players_online"`        // Online:
	Location            string `json:"location"`              // Location:
	PvpType             string `json:"pvp_type"`              // PvP Type:
	PremiumOnly         bool   `json:"premium_only"`          // Additional Information: premium = true / else: false
	TransferType        string `json:"transfer_type"`         // Additional Information: regular (if not present) / locked / blocked
	BattleyeProtected   bool   `json:"battleye_protected"`    // BattlEye Status: true if protected / false if "Not protected by BattlEye."
	BattleyeDate        string `json:"battleye_date"`         // BattlEye Status: null if since release / else show date?
	GameWorldType       string `json:"game_world_type"`       // BattlEye Status: regular / experimental / tournament (if Tournament World Type exists)
	TournamentWorldType string `json:"tournament_world_type"` // BattlEye Status: null (default?) / regular / restricted
}

// Child of JSONData
type Worlds struct {
	PlayersOnline    int     `json:"players_online"` // Calculated value
	RecordPlayers    int     `json:"record_players"` // Overall Maximum:
	RecordDate       string  `json:"record_date"`    // Overall Maximum:
	RegularWorlds    []World `json:"regular_worlds"`
	TournamentWorlds []World `json:"tournament_worlds"`
}

//
// The base includes two levels: Worlds and Information
type WorldsOverviewResponse struct {
	Worlds      Worlds      `json:"worlds"`
	Information Information `json:"information"`
}

var (
	worldPlayerRecordRegex           = regexp.MustCompile(`.*<\/b>...(.*) players \(on (.*)\)`)
	worldInformationRegex            = regexp.MustCompile(`.*world=.*">(.*)<\/a><\/td>.*right;">(.*)<\/td><td>(.*)<\/td><td>(.*)<\/td><td align="center" valign="middle">(.*)<\/td><td>(.*)<\/td>`)
	worldBattlEyeProtectedSinceRegex = regexp.MustCompile(`.*game world has been protected by BattlEye since (.*).&lt;\/p.*`)
)

func TibiaWorldsOverviewV3(c *gin.Context) {
	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=worlds"
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaWorldsOverviewV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaWorldsOverviewV3Impl(BoxContentHTML)

	TibiaDataAPIHandleSuccessResponse(c, "TibiaWorldsOverviewV3", jsonData)
}

// TibiaWorldsOverviewV3 func
func TibiaWorldsOverviewV3Impl(BoxContentHTML string) WorldsOverviewResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty vars
	var (
		RegularWorldsData, TournamentWorldsData                                                                                                     []World
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
		subma1 := worldPlayerRecordRegex.FindAllStringSubmatch(WorldsDivHTML, -1)

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

		subma2 := worldInformationRegex.FindAllStringSubmatch(WorldsDivHTML, -1)

		// check if regex return length is over 0
		if len(subma2) > 0 {
			WorldsPlayersOnline := 0

			if subma2[0][2] == "-" {
				WorldsStatus = "unknown"
			} else {
				WorldsPlayersOnline = TibiadataStringToIntegerV3(subma2[0][2])

				// Setting the players_online & overall players_online
				WorldsAllOnlinePlayers += WorldsPlayersOnline

				if WorldsPlayersOnline > 0 {
					WorldsStatus = "online"
				} else {
					WorldsStatus = "offline"
				}
			}

			// Creating better to use vars
			WorldsBattlEye := subma2[0][5]
			WorldsAdditionalInfo := subma2[0][6]

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
					subma21 := worldBattlEyeProtectedSinceRegex.FindAllStringSubmatch(WorldsBattlEye, -1)
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
			OneWorld := World{
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
	return WorldsOverviewResponse{
		Worlds{
			PlayersOnline:    WorldsAllOnlinePlayers,
			RecordPlayers:    WorldsRecordPlayers,
			RecordDate:       WorldsRecordDate,
			RegularWorlds:    RegularWorldsData,
			TournamentWorlds: TournamentWorldsData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
