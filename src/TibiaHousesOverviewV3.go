package main

import (
	"fmt"
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
	IsFinished  bool   `json:"finished"`
}

// Child of HousesHouses
type HousesHouse struct {
	Name        string        `json:"name"`
	HouseID     int           `json:"house_id"`
	Size        int           `json:"size"`
	Rent        int           `json:"rent"`
	IsRented    bool          `json:"rented"`
	IsAuctioned bool          `json:"auctioned"`
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
func TibiaHousesOverviewV3Impl(c *gin.Context, world string, town string, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) (*HousesOverviewResponse, error) {
	var (
		// Creating empty vars
		HouseData, GuildhallData []HousesHouse
	)

	// list of different fansite types
	HouseTypes := []string{"houses", "guildhalls"}

	// running over the FansiteTypes array
	for _, HouseType := range HouseTypes {
		houses, err := makeHouseRequest(HouseType, world, town, htmlDataCollector)
		if err != nil {
			return nil, fmt.Errorf("[error] TibiaHousesOverviewV3Impl failed at makeHouseRequest, type: %s, err: %s", HouseType, err)
		}

		switch HouseType {
		case "houses":
			HouseData = houses
		case "guildhalls":
			GuildhallData = houses
		}
	}

	// Build the data-blob
	return &HousesOverviewResponse{
		HousesHouses{
			World:         world,
			Town:          town,
			HouseList:     HouseData,
			GuildhallList: GuildhallData,
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

func makeHouseRequest(HouseType, world, town string, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) ([]HousesHouse, error) {
	// Creating an empty var
	var output []HousesHouse

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=houses&world=" + TibiaDataQueryEscapeStringV3(world) + "&town=" + TibiaDataQueryEscapeStringV3(town) + "&type=" + TibiaDataQueryEscapeStringV3(HouseType),
	}

	BoxContentHTML, err := htmlDataCollector(tibiadataRequest)
	// return error (e.g. for maintenance mode)
	if err != nil {
		return nil, err
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, err
	}

	var insideError error

	ReaderHTML.Find(".TableContentContainer .TableContent tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		house := HousesHouse{}

		// Storing HTML into HousesDivHTML
		HousesDivHTML, err := s.Html()
		if err != nil {
			insideError = err
			return false
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
					house.Auction.IsFinished = true
				} else {
					house.Auction.AuctionLeft = subma1b[0][3]
				}
			}

			output = append(output, house)
		}

		return true
	})

	return output, insideError
}
