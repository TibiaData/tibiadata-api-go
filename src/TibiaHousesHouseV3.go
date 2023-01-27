package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TibiaData/tibiadata-api-go/src/validation"
)

// Child of Status
type HouseRental struct {
	Owner            string `json:"owner"`
	OwnerSex         string `json:"owner_sex"`
	PaidUntil        string `json:"paid_until"`
	MovingDate       string `json:"moving_date"`
	TransferReceiver string `json:"transfer_receiver"`
	TransferPrice    int    `json:"transfer_price"`
	TransferAccept   bool   `json:"transfer_accept"`
}

// Child of Status
type HouseAuction struct {
	CurrentBid     int    `json:"current_bid"`
	CurrentBidder  string `json:"current_bidder"`
	AuctionOngoing bool   `json:"auction_ongoing"`
	AuctionEnd     string `json:"auction_end"`
}

// Child of House
type HouseStatus struct {
	IsAuctioned   bool         `json:"is_auctioned"`
	IsRented      bool         `json:"is_rented"`
	IsMoving      bool         `json:"is_moving"`
	IsTransfering bool         `json:"is_transfering"`
	Auction       HouseAuction `json:"auction"`
	Rental        HouseRental  `json:"rental"`
	Original      string       `json:"original"`
}

// Child of JSONData
type House struct {
	Houseid int         `json:"houseid"`
	World   string      `json:"world"`
	Town    string      `json:"town,omitempty"`
	Name    string      `json:"name"`
	Type    string      `json:"type,omitempty"`
	Beds    int         `json:"beds"`
	Size    int         `json:"size"`
	Rent    int         `json:"rent"`
	Img     string      `json:"img"`
	Status  HouseStatus `json:"status"`
}

//
// The base includes two levels: Houses and Information
type HouseResponse struct {
	House       House       `json:"house"`
	Information Information `json:"information"`
}

var (
	houseDataRegex = regexp.MustCompile(`<td.*src="(.*)" width.*<b>(.*)<\/b>.*This (house|guildhall) can.*to ([0-9]+) beds?..*<b>([0-9]+) square.*<b>([0-9]+)([k]+).gold<\/b>.*on <b>([A-Za-z]+)<\/b>.(.*)<\/td>`)
	// matching for this: and <wants to|will> pass the <HouseType> to <TransferReceiver> for <TransferPrice> gold
	housePassingRegex = regexp.MustCompile(`and (wants to|will) pass the (house|guildhall) to (.*) for ([0-9]+) gold`)
	// matching for this: <OwnerSex> will move out on <MovingDate> (
	moveOutRegex = regexp.MustCompile(`(He|She) will move out on (.*?) \(`)
	// matching for this: The <HouseType> has been rented by <Owner>. <OwnerSex> has paid the rent until <PaidUntil>.
	paidUntilRegex      = regexp.MustCompile(`The (house|guildhall) has been rented by (.*). (He|She) has paid.*until (.*?)\.`)
	houseAuctionedRegex = regexp.MustCompile(`The (house|guildhall) is currently.*The auction (will end|has ended) at (.*)\. The.*is ([0-9]+) gold.*submitted by (.*)\.`)
)

// TibiaHousesHouseV3 func
func TibiaHousesHouseV3Impl(houseid int, BoxContentHTML string) (*HouseResponse, error) {
	// Creating empty vars
	var HouseData House

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaHousesHouseV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Running query over each div
	HouseHTML, err := ReaderHTML.Find(".BoxContent table tr").First().Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaHousesHouseV3Impl failed at ReaderHTML.Find, err: %s", err)
	}

	// Regex to get data for house
	subma1 := houseDataRegex.FindAllStringSubmatch(HouseHTML, -1)

	if len(subma1) > 0 {
		HouseData.Houseid = houseid
		HouseData.World = subma1[0][8]

		rawHouse, err := validation.GetHouseRaw(HouseData.Houseid)
		if err != nil {
			return nil, err
		}

		HouseData.Town = rawHouse.Town
		HouseData.Type = rawHouse.Type

		HouseData.Name = TibiaDataSanitizeEscapedString(subma1[0][2])
		HouseData.Img = subma1[0][1]
		HouseData.Beds = TibiaDataStringToIntegerV3(subma1[0][4])
		HouseData.Size = TibiaDataStringToIntegerV3(subma1[0][5])
		HouseData.Rent = TibiaDataConvertValuesWithK(subma1[0][6] + subma1[0][7])

		HouseData.Status.Original = strings.TrimSpace(TibiaDataSanitizeStrings(TibiaDataSanitizeEscapedString(RemoveHtmlTag(subma1[0][9]))))

		switch {
		case strings.Contains(HouseData.Status.Original, "has been rented by"):
			// rented
			switch {
			case strings.Contains(HouseData.Status.Original, " pass the house to "), strings.Contains(HouseData.Status.Original, " pass the guildhall to "):
				HouseData.Status.IsTransfering = true

				subma2 := housePassingRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				if subma2[0][1] == "will" {
					HouseData.Status.Rental.TransferAccept = true
				}
				HouseData.Status.Rental.TransferReceiver = subma2[0][3]
				HouseData.Status.Rental.TransferPrice = TibiaDataStringToIntegerV3(subma2[0][4])
				fallthrough

			case strings.Contains(HouseData.Status.Original, " will move out on "):
				HouseData.Status.IsMoving = true
				subma2 := moveOutRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.MovingDate = TibiaDataDatetimeV3(subma2[0][2])
				fallthrough

			default:
				HouseData.Status.IsRented = true
				subma2 := paidUntilRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.Owner = subma2[0][2]
				HouseData.Status.Rental.PaidUntil = TibiaDataDatetimeV3(subma2[0][4])
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
				subma2 := houseAuctionedRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Auction.AuctionEnd = TibiaDataDatetimeV3(subma2[0][3])
				HouseData.Status.Auction.CurrentBid = TibiaDataStringToIntegerV3(subma2[0][4])
				HouseData.Status.Auction.CurrentBidder = TibiaDataSanitizeStrings(subma2[0][5])
				if subma2[0][2] == "will end" {
					HouseData.Status.Auction.AuctionOngoing = true
				}
			}
		}
	}

	// Build the data-blob
	return &HouseResponse{
		HouseData,
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
