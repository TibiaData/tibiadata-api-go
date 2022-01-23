package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCormaya10(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/houses/Premia/Edron/Cormaya10.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	houseJson := TibiaHousesHouseV3Impl("54025", string(data))
	assert := assert.New(t)

	assert.Equal(54025, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("", houseJson.House.Town) //depends on TibiaDataHousesMapResolver
	assert.Equal("Cormaya 10", houseJson.House.Name)
	assert.Equal("", houseJson.House.Type) //depends on TibiaDataHousesMapResolver
	assert.Equal(3, houseJson.House.Beds)
	assert.Equal(80, houseJson.House.Size)
	assert.Equal(300000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_54025.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.False(houseStatus.IsAuctioned)
	assert.True(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseAuction{CurrentBid: 0, CurrentBidder: "", AuctionOngoing: false, AuctionEnd: ""}, houseStatus.Auction)
	assert.Equal("The house has been rented by Xendor of Askara. He has paid the rent until Feb 02 2022, 10:05:26 CET.", houseStatus.Original)

	houseRental := houseJson.House.Status.Rental
	assert.NotNil(houseRental)
	assert.Equal("Xendor of Askara", houseRental.Owner)
	assert.Equal("male", houseRental.OwnerSex)
	assert.Equal("2022-02-02T09:05:26Z", houseRental.PaidUntil)
	assert.Equal("", houseRental.MovingDate)
	assert.Equal("", houseRental.TransferReceiver)
	assert.Equal(0, houseRental.TransferPrice)
	assert.False(houseRental.TransferAccept)
}

func TestCormaya11(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/houses/Premia/Edron/Cormaya11.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	houseJson := TibiaHousesHouseV3Impl("54026", string(data))
	assert := assert.New(t)

	assert.Equal(54026, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("", houseJson.House.Town) //depends on TibiaDataHousesMapResolver
	assert.Equal("Cormaya 11", houseJson.House.Name)
	assert.Equal("", houseJson.House.Type) //depends on TibiaDataHousesMapResolver
	assert.Equal(2, houseJson.House.Beds)
	assert.Equal(43, houseJson.House.Size)
	assert.Equal(150000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_54026.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.True(houseStatus.IsAuctioned)
	assert.False(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseRental{Owner: "", OwnerSex: "", PaidUntil: "", MovingDate: "", TransferReceiver: "", TransferPrice: 0, TransferAccept: false}, houseStatus.Rental)
	assert.Equal("The house is currently being auctioned. The auction will end at Jan 21 2022, 10:00:00 CET. The highest bid so far is 200000 gold and has been submitted by Ciuchy Szajba.", houseStatus.Original)

	houseAuction := houseJson.House.Status.Auction
	assert.NotNil(houseAuction)
	assert.Equal(200000, houseAuction.CurrentBid)
	assert.Equal("Ciuchy Szajba", houseAuction.CurrentBidder)
	assert.True(houseAuction.AuctionOngoing)
	assert.Equal("2022-01-21T09:00:00Z", houseAuction.AuctionEnd)
}
