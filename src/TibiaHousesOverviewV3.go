package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Child of House
type HousesAction struct {
	AuctionBid  int    `json:"current_bid"`
	AuctionLeft string `json:"time_left"`
}

// Child of HousesHouses
type HousesHouse struct {
	Name        string       `json:"name"`
	HouseID     int          `json:"house_id"`
	Size        int          `json:"size"`
	Rent        int          `json:"rent"`
	IsRented    bool         `json:"rented"`
	IsAuctioned bool         `json:"auctioned"`
	Auction     HousesAction `json:"auction"`
}

// Child of JSONData
type HousesHouses struct {
	World         string        `json:"world"`
	Town          string        `json:"town"`
	HouseList     []HousesHouse `json:"house_list"`
	GuildhallList []HousesHouse `json:"guildhall_list"`
}

// The base includes two levels: HousesHouses and Information
type HousesOverviewResponse struct {
	Houses      HousesHouses `json:"houses"`
	Information Information  `json:"information"`
}

// TibiaHousesOverviewV3 func
func TibiaHousesOverviewV3Impl(c *gin.Context, world string, town string) HousesOverviewResponse {
	var (
		// Creating empty vars
		HouseData, GuildhallData []HousesHouse
	)

	// list of different fansite types
	HouseTypes := []string{"houses", "guildhalls"}
	// running over the FansiteTypes array
	for _, HouseType := range HouseTypes {
		tibiadataRequest := TibiadataRequestStruct{
			Method: resty.MethodGet,
			URL:    "https://www.tibia.com/community/?subtopic=houses&world=" + TibiadataQueryEscapeStringV3(world) + "&town=" + TibiadataQueryEscapeStringV3(town) + "&type=" + TibiadataQueryEscapeStringV3(HouseType),
		}
		BoxContentHTML, err := TibiadataHTMLDataCollectorV3(tibiadataRequest)

		// return error (e.g. for maintenance mode)
		if err != nil {
			TibiaDataAPIHandleResponse(c, http.StatusBadGateway, "TibiaHousesOverviewV3", gin.H{"error": err.Error()})

			//TODO: Need to refactor this properly
			return HousesOverviewResponse{}
		}

		// Loading HTML data into ReaderHTML for goquery with NewReader
		ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
		if err != nil {
			log.Fatal(err)
		}

		ReaderHTML.Find(".TableContentContainer .TableContent tr").Each(func(index int, s *goquery.Selection) {
			house := HousesHouse{}

			// Storing HTML into HousesDivHTML
			HousesDivHTML, err := s.Html()
			if err != nil {
				log.Fatal(err)
			}

			// Removing linebreaks from HTML
			HousesDivHTML = TibiadataHTMLRemoveLinebreaksV3(HousesDivHTML)
			HousesDivHTML = TibiaDataSanitizeNbspSpaceString(HousesDivHTML)

			// Regex to get data for record values
			regex1 := regexp.MustCompile(`<td.*><nobr>(.*)<\/nobr><\/td><td.*><nobr>([0-9]+).sqm<\/nobr><\/td><td.*><nobr>([0-9]+)(k+).gold<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td>.*houseid" value="([0-9]+)"\/><div.*`)
			subma1 := regex1.FindAllStringSubmatch(HousesDivHTML, -1)

			if len(subma1) > 0 {
				// House details
				house.Name = TibiaDataSanitizeEscapedString(subma1[0][1])
				house.HouseID = TibiadataStringToIntegerV3(subma1[0][6])
				house.Size = TibiadataStringToIntegerV3(subma1[0][2])
				house.Rent = TibiaDataConvertValuesWithK(subma1[0][3] + subma1[0][4])

				// HousesAction details
				s := subma1[0][5]
				switch {
				case strings.Contains(s, "rented"):
					house.IsRented = true
				case strings.Contains(s, "auctioned (no bid yet)"):
					house.IsAuctioned = true
				case strings.Contains(s, "auctioned"):
					house.IsAuctioned = true
					regex1b := regexp.MustCompile(`auctioned.\(([0-9]+).gold;.(.*).left\)`)
					subma1b := regex1b.FindAllStringSubmatch(s, -1)
					house.Auction.AuctionBid = TibiadataStringToIntegerV3(subma1b[0][1])
					house.Auction.AuctionLeft = subma1b[0][2]
				}

				// append house to list houses/guildhalls
				switch HouseType {
				case "houses":
					HouseData = append(HouseData, house)
				case "guildhalls":
					GuildhallData = append(GuildhallData, house)
				}
			}
		})
	}

	// Build the data-blob
	return HousesOverviewResponse{
		HousesHouses{
			World:         world,
			Town:          town,
			HouseList:     HouseData,
			GuildhallList: GuildhallData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
