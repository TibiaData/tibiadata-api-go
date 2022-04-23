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
type HousesAuction struct {
	AuctionBid  int    `json:"current_bid"`
	AuctionLeft string `json:"time_left"`
}

// Child of HousesHouses
type HousesHouse struct {
	Name        string        `json:"name"`
	HouseID     int           `json:"house_id"`
	Size        int           `json:"size"`
	Rent        int           `json:"rent"`
	IsRented    bool          `json:"rented"`
	IsAuctioned bool          `json:"auctioned"`
	IsFinished  bool          `json:"finished"`
	Auction     HousesAuction `json:"auction"`
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

var (
	houseOverviewDataRegex      = regexp.MustCompile(`<td.*><nobr>(.*)<\/nobr><\/td><td.*><nobr>([0-9]+).sqm<\/nobr><\/td><td.*><nobr>([0-9]+)(k+).gold<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td>.*houseid" value="([0-9]+)"\/><div.*`)
	houseOverviewAuctionedRegex = regexp.MustCompile(`auctioned.\(([0-9]+).gold;.(finished|(.*).left)\)`)
)

// TibiaHousesOverviewV3 func
func TibiaHousesOverviewV3Impl(c *gin.Context, world string, town string, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) HousesOverviewResponse {
	var (
		// Creating empty vars
		HouseData, GuildhallData []HousesHouse
	)

	// list of different fansite types
	HouseTypes := []string{"houses", "guildhalls"}
	// running over the FansiteTypes array
	for _, HouseType := range HouseTypes {
		tibiadataRequest := TibiaDataRequestStruct{
			Method: resty.MethodGet,
			URL:    "https://www.tibia.com/community/?subtopic=houses&world=" + TibiaDataQueryEscapeStringV3(world) + "&town=" + TibiaDataQueryEscapeStringV3(town) + "&type=" + TibiaDataQueryEscapeStringV3(HouseType),
		}
		BoxContentHTML, err := htmlDataCollector(tibiadataRequest)

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
			HousesDivHTML = TibiaDataHTMLRemoveLinebreaksV3(HousesDivHTML)
			HousesDivHTML = TibiaDataSanitizeStrings(HousesDivHTML)

			subma1 := houseOverviewDataRegex.FindAllStringSubmatch(HousesDivHTML, -1)

			if len(subma1) > 0 {
				// House details
				house.Name = TibiaDataSanitizeEscapedString(subma1[0][1])
				house.HouseID = TibiaDataStringToIntegerV3(subma1[0][6])
				house.Size = TibiaDataStringToIntegerV3(subma1[0][2])
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
					subma1b := houseOverviewAuctionedRegex.FindAllStringSubmatch(s, -1)
					house.Auction.AuctionBid = TibiaDataStringToIntegerV3(subma1b[0][1])
					if subma1b[0][2] == "finished" {
						house.IsFinished = true
					} else {
						house.Auction.AuctionLeft = subma1b[0][3]
					}
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
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
