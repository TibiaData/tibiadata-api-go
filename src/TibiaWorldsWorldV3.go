package main

import (
	"log"
	"regexp"
	"strings"
	"tibiadata-api-go/src/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaWorldsWorldV3 func
func TibiaWorldsWorldV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// The base includes two levels: World and Information
	type JSONData struct {
		World       structs.World       `json:"worlds"`
		Information structs.Information `json:"information"`
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=worlds&world=" + TibiadataQueryEscapeStringV3(world)
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty vars
	var (
		WorldsStatus, WorldsRecordDate, WorldsCreationDate, WorldsLocation, WorldsPvpType, WorldsTransferType, WorldsBattleyeDate, WorldsGameWorldType, WorldsTournamentWorldType string
		WorldsQuestTitles                                                                                                                                                         []string
		WorldsPlayersOnline, WorldsRecordPlayers                                                                                                                                  int
		WorldsPremiumOnly, WorldsBattleyeProtected                                                                                                                                bool
		WorldsOnlinePlayers                                                                                                                                                       []structs.OnlinePlayer
	)

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

			switch WorldsInformationLeftColumn {
			case "Status":
				switch {
				case strings.Contains(WorldsInformationRightColumn, "</div>Online"):
					WorldsStatus = "online"
				case strings.Contains(WorldsInformationRightColumn, "</div>Offline"):
					WorldsStatus = "offline"
				default:
					WorldsStatus = "unknown"
				}
			case "Players Online":
				WorldsPlayersOnline = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			case "Online Record":
				// Regex to get data for record values
				regex2 := regexp.MustCompile(`(.*) players \(on (.*)\)`)
				subma2 := regex2.FindAllStringSubmatch(WorldsInformationRightColumn, -1)

				if len(subma2) > 0 {
					// setting record values
					WorldsRecordPlayers = TibiadataStringToIntegerV3(subma2[0][1])
					WorldsRecordDate = TibiadataDatetimeV3(subma2[0][2])
				}
			case "Creation Date":
				WorldsCreationDate = WorldsInformationRightColumn
			case "Location":
				WorldsLocation = WorldsInformationRightColumn
			case "PvP Type":
				WorldsPvpType = WorldsInformationRightColumn
			case "Premium Type":
				WorldsPremiumOnly = true
			case "Transfer Type":
				WorldsTransferType = WorldsInformationRightColumn
			case "World Quest Titles":
				if WorldsInformationRightColumn != "This game world currently has no title." {
					WorldsQuestTitlesTmp := strings.Split(WorldsInformationRightColumn, ", ")
					for _, str := range WorldsQuestTitlesTmp {
						if str != "" {
							WorldsQuestTitles = append(WorldsQuestTitles, TibiadataRemoveURLsV3(str))
						}
					}
				}
			case "BattlEye Status":
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
			case "Game World Type":
				WorldsGameWorldType = strings.ToLower(WorldsInformationRightColumn)
			case "Tournament World Type":
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
			WorldsOnlinePlayers = append(WorldsOnlinePlayers, structs.OnlinePlayer{
				Name:     subma1[0][1],
				Level:    TibiadataStringToIntegerV3(subma1[0][2]),
				Vocation: subma1[0][3],
			})
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		World: structs.World{
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
		Information: structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaWorldsWorldV3", jsonData)
}
