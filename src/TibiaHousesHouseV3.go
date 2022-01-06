package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaHousesHouseV3 func
func TibiaHousesHouseV3(c *gin.Context) {

	// getting params from URL
	world := c.Param("world")
	houseid := c.Param("houseid")

	// Child of Status
	type Rental struct {
		Owner            string `json:"owner"`
		OwnerSex         string `json:"owner_sex"`
		PaidUntil        string `json:"paid_until"`
		MovingDate       string `json:"moving_date"`
		TransferReceiver string `json:"transfer_receiver"`
		TransferPrice    int    `json:"transfer_price"`
		TransferAccept   bool   `json:"transfer_accept"`
	}

	// Child of Status
	type Auction struct {
		CurrentBid     int    `json:"current_bid"`
		CurrentBidder  string `json:"current_bidder"`
		AuctionOngoing bool   `json:"auction_ongoing"`
		AuctionEnd     string `json:"auction_end"`
	}

	// Child of House
	type Status struct {
		IsAuctioned   bool    `json:"is_auctioned"`
		IsRented      bool    `json:"is_rented"`
		IsMoving      bool    `json:"is_moving"`
		IsTransfering bool    `json:"is_transfering"`
		Auction       Auction `json:"auction"`
		Rental        Rental  `json:"rental"`
		Original      string  `json:"original"`
	}

	// Child of JSONData
	type House struct {
		Houseid int    `json:"houseid"`
		World   string `json:"world"`
		Town    string `json:"town"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Beds    int    `json:"beds"`
		Size    int    `json:"size"`
		Rent    int    `json:"rent"`
		Img     string `json:"img"`
		Status  Status `json:"status"`
	}

	//
	// The base includes two levels: Houses and Information
	type JSONData struct {
		House       House       `json:"house"`
		Information Information `json:"information"`
	}

	// Creating empty vars
	var HouseData House

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=houses&page=view&world=" + TibiadataQueryEscapeStringV3(world) + "&houseid=" + TibiadataQueryEscapeStringV3(houseid)
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Running query over each div
	HouseHTML, err := ReaderHTML.Find(".BoxContent table tr").First().Html()

	if err != nil {
		log.Fatal(err)
	}

	// Regex to get data for house
	regex1 := regexp.MustCompile(`<td.*src="(.*)" width.*<b>(.*)<\/b>.*This (house|guildhall) can.*to ([0-9]+) beds..*<b>([0-9]+) square.*<b>([0-9]+)([k]+).gold<\/b>.*on <b>([A-Za-z]+)<\/b>.(.*)<\/td>`)
	subma1 := regex1.FindAllStringSubmatch(HouseHTML, -1)

	if len(subma1) > 0 {
		HouseData.Houseid = TibiadataStringToIntegerV3(houseid)
		HouseData.World = subma1[0][8]

		HouseData.Name = TibiaDataSanitizeEscapedString(subma1[0][2])
		HouseData.Img = subma1[0][1]
		HouseData.Type = subma1[0][3]
		HouseData.Beds = TibiadataStringToIntegerV3(subma1[0][4])
		HouseData.Size = TibiadataStringToIntegerV3(subma1[0][5])
		HouseData.Rent = TibiaDataConvertValuesWithK(subma1[0][6] + subma1[0][7])

		HouseData.Status.Original = TibiaDataSanitizeEscapedString(RemoveHtmlTag(subma1[0][9]))

		switch {
		case strings.Contains(HouseData.Status.Original, "has been rented by"):
			// rented

			switch {
			case strings.Contains(HouseData.Status.Original, " pass the "+HouseData.Type+" to "):
				HouseData.Status.IsTransfering = true
				// matching for this: and <wants to|will> pass the <HouseType> to <TransferReceiver> for <TransferPrice> gold
				regex2 := regexp.MustCompile(`and (wants to|will) pass the (house|guildhall) to (.*) for ([0-9]+) gold`)
				subma2 := regex2.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				if subma2[0][1] == "will" {
					HouseData.Status.Rental.TransferAccept = true
				}
				HouseData.Status.Rental.TransferReceiver = subma2[0][3]
				HouseData.Status.Rental.TransferPrice = TibiadataStringToIntegerV3(subma2[0][4])
				fallthrough

			case strings.Contains(HouseData.Status.Original, " will move out on "):
				HouseData.Status.IsMoving = true
				// matching for this: <OwnerSex> will move out on <MovingDate> (
				regex2 := regexp.MustCompile(`(He|She) will move out on (.*?) \(`)
				subma2 := regex2.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.MovingDate = TibiadataDatetimeV3(subma2[0][2])
				fallthrough

			default:
				HouseData.Status.IsRented = true
				// matching for this: The <HouseType> has been rented by <Owner>. <OwnerSex> has paid the rent until <PaidUntil>.
				regex2 := regexp.MustCompile(`The (house|guildhall) has been rented by (.*). (He|She) has paid.*until (.*?)\.`)
				subma2 := regex2.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.Owner = subma2[0][2]
				HouseData.Status.Rental.PaidUntil = TibiadataDatetimeV3(subma2[0][4])
				switch subma2[0][3] {
				case "She":
					HouseData.Status.Rental.OwnerSex = "female"
				case "He":
					HouseData.Status.Rental.OwnerSex = "male"
				}
			}

		case strings.Contains(HouseData.Status.Original, "is currently being auctioned"):
			// auctioned
			HouseData.Status.IsAuctioned = true

			// check if bid is going on
			if !strings.Contains(HouseData.Status.Original, "No bid has been submitted so far.") {
				regex2 := regexp.MustCompile(`The (house|guildhall) is currently.*The auction (will end|has ended) at (.*)\. The.*is ([0-9]+) gold.*submitted by (.*)\.`)
				subma2 := regex2.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Auction.AuctionEnd = TibiadataDatetimeV3(subma2[0][3])
				HouseData.Status.Auction.CurrentBid = TibiadataStringToIntegerV3(subma2[0][4])
				HouseData.Status.Auction.CurrentBidder = subma2[0][5]
				if subma2[0][2] == "will end" {
					HouseData.Status.Auction.AuctionOngoing = true
				}
			}
		}

	}

	//
	// Build the data-blob
	jsonData := JSONData{
		HouseData,
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaHousesHouseV3", jsonData)
}
