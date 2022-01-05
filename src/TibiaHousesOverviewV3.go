package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaHousesOverviewV3 func
func TibiaHousesOverviewV3(c *gin.Context) {

	// getting params from URL
	world := c.Param("world")
	town := c.Param("town")

	// Child of House
	type Auction struct {
		AuctionBid  int    `json:"current_bid"`
		AuctionLeft string `json:"time_left"`
	}

	// Child of Houses
	type House struct {
		Name        string  `json:"name"`
		HouseID     int     `json:"house_id"`
		Size        int     `json:"size"`
		Rent        int     `json:"rent"`
		IsRented    bool    `json:"rented"`
		IsAuctioned bool    `json:"auctioned"`
		Auction     Auction `json:"auction"`
	}

	// Child of JSONData
	type Houses struct {
		World         string  `json:"world"`
		Town          string  `json:"town"`
		HouseList     []House `json:"house_list"`
		GuildhallList []House `json:"guildhall_list"`
	}

	//
	// The base includes two levels: Houses and Information
	type JSONData struct {
		Houses      Houses      `json:"houses"`
		Information Information `json:"information"`
	}

	// Creating empty vars
	var HouseData, GuildhallData []House
	HouseData = []House{}
	GuildhallData = []House{}

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)
	town = TibiadataStringWorldFormatToTitleV3(town)

	// list of different house types
	HouseTypes := []string{"houses", "guildhalls"}
	// running over the HouseTypes array
	for _, HouseType := range HouseTypes {

		// Getting data with TibiadataHTMLDataCollectorV3
		BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/community/?subtopic=houses&world=" + TibiadataQueryEscapeStringV3(world) + "&town=" + TibiadataQueryEscapeStringV3(town) + "&type=" + TibiadataQueryEscapeStringV3(HouseType))

		// Loading HTML data into ReaderHTML for goquery with NewReader
		ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
		if err != nil {
			log.Fatal(err)
		}

		// Running query over each div
		ReaderHTML.Find(".BoxContent table tr").Each(func(index int, s *goquery.Selection) {

			// Storing HTML into HousesDivHTML
			HousesDivHTML, err := s.Html()
			if err != nil {
				log.Fatal(err)
			}

			// Removing linebreaks from HTML
			HousesDivHTML = TibiadataHTMLRemoveLinebreaksV3(HousesDivHTML)

			// Regex to get data for record values
			regex1 := regexp.MustCompile(`<td.*><nobr>(.*)<\/nobr><\/td><td.*><nobr>([0-9]+).sqm<\/nobr><\/td><td.*><nobr>([0-9]+)(k+).gold<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td>.*houseid" value="([0-9]+)"\/><div.*`)
			subma1 := regex1.FindAllStringSubmatch(HousesDivHTML, -1)

			if len(subma1) > 0 {

				// Default of vars
				IsRented := false
				IsAuctioned := false
				AuctionBid := 0
				AuctionLeft := ""

				switch {
				case strings.Contains(subma1[0][5], "rented"):
					IsRented = true
				case strings.Contains(subma1[0][5], "no bid yet"):
					// nothing to set?
				case strings.Contains(subma1[0][5], "auctioned"):
					IsAuctioned = true
					regex1b := regexp.MustCompile(`auctioned..([0-9]+).gold..(.*).`)
					subma1b := regex1b.FindAllStringSubmatch(subma1[0][5], -1)
					AuctionBid = TibiadataStringToIntegerV3(subma1b[0][1])
					AuctionLeft = subma1b[0][2]
				}

				// Name
				Name := TibiaDataSanitizeEscapedString(subma1[0][1])
				// HouseID
				HouseID := TibiadataStringToIntegerV3(subma1[0][6])
				// Size
				Size := TibiadataStringToIntegerV3(subma1[0][2])
				// Rent
				Rent := TibiadataStringToIntegerV3(subma1[0][3] + strings.Repeat("000", strings.Count(subma1[0][4], "k")))

				switch HouseType {
				case "houses":
					HouseData = append(HouseData, House{
						Name:        Name,
						HouseID:     HouseID,
						Size:        Size,
						Rent:        Rent,
						IsRented:    IsRented,
						IsAuctioned: IsAuctioned,
						Auction: Auction{
							AuctionBid:  AuctionBid,
							AuctionLeft: AuctionLeft,
						},
					})
				case "guildhalls":
					GuildhallData = append(GuildhallData, House{
						Name:        Name,
						HouseID:     HouseID,
						Size:        Size,
						Rent:        Rent,
						IsRented:    IsRented,
						IsAuctioned: IsAuctioned,
						Auction: Auction{
							AuctionBid:  AuctionBid,
							AuctionLeft: AuctionLeft,
						},
					})
				}
			}
		})

	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Houses{
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

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaHousesOverviewV3", jsonData)
}
