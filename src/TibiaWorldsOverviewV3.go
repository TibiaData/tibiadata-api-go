package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of Worlds
type OverviewWorld struct {
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
type OverviewWorlds struct {
	PlayersOnline    int             `json:"players_online"` // Calculated value
	RecordPlayers    int             `json:"record_players"` // Overall Maximum:
	RecordDate       string          `json:"record_date"`    // Overall Maximum:
	RegularWorlds    []OverviewWorld `json:"regular_worlds"`
	TournamentWorlds []OverviewWorld `json:"tournament_worlds"`
}

//
// The base includes two levels: Worlds and Information
type WorldsOverviewResponse struct {
	Worlds      OverviewWorlds `json:"worlds"`
	Information Information    `json:"information"`
}

var (
	worldPlayerRecordRegex           = regexp.MustCompile(`.*<\/b>...(.*) players \(on (.*)\)`)
	worldInformationRegex            = regexp.MustCompile(`.*world=.*">(.*)<\/a><\/td>.*right;">(.*)<\/td><td>(.*)<\/td><td>(.*)<\/td><td align="center" valign="middle">(.*)<\/td><td>(.*)<\/td>`)
	worldBattlEyeProtectedSinceRegex = regexp.MustCompile(`.*game world has been protected by BattlEye since (.*).&lt;\/p.*`)
)

// TibiaWorldsOverviewV3 func
func TibiaWorldsOverviewV3Impl(BoxContentHTML string) (*WorldsOverviewResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaWorldsOverviewV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty vars
	var (
		RegularWorldsData, TournamentWorldsData                                                                                                     []OverviewWorld
		WorldsRecordDate, WorldsWorldCategory, WorldsBattleyeDate, WorldsTransferType, WorldsTournamentWorldType, WorldsGameWorldType, WorldsStatus string
		WorldsRecordPlayers, WorldsAllOnlinePlayers                                                                                                 int
		WorldsPremiumOnly, WorldsBattleyeProtected                                                                                                  bool

		insideError error
	)

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		WorldsDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaWorldsOverviewV3Impl failed at WorldsDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Regex to get data for record values
		subma1 := worldPlayerRecordRegex.FindAllStringSubmatch(WorldsDivHTML, -1)

		if len(subma1) > 0 {
			// setting record values
			WorldsRecordPlayers = TibiaDataStringToIntegerV3(subma1[0][1])
			WorldsRecordDate = TibiaDataDatetimeV3(subma1[0][2])
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
				WorldsPlayersOnline = TibiaDataStringToIntegerV3(subma2[0][2])

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
					WorldsBattleyeDate = TibiaDataDateV3(subma21[0][1])
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
			OneWorld := OverviewWorld{
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

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &WorldsOverviewResponse{
		OverviewWorlds{
			PlayersOnline:    WorldsAllOnlinePlayers,
			RecordPlayers:    WorldsRecordPlayers,
			RecordDate:       WorldsRecordDate,
			RegularWorlds:    RegularWorldsData,
			TournamentWorlds: TournamentWorldsData,
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
