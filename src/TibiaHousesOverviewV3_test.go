package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnticaThaisHousesOverview(t *testing.T) {
	houseData, err := os.ReadFile("../testdata/houses/overview/AnticaThaisHouses.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	guildData, err := os.ReadFile("../testdata/houses/overview/AnticaThaisGuilds.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	housesJson := TibiaHousesOverviewV3Impl(
		nil,
		"Antica",
		"Thais",
		func(request TibiaDataRequestStruct) (string, error) {
			if strings.Contains(request.URL, "guildhalls") {
				return string(guildData), nil
			}

			return string(houseData), nil
		})
	assert := assert.New(t)

	assert.Equal("Antica", housesJson.Houses.World)
	assert.Equal("Thais", housesJson.Houses.Town)

	assert.Equal(118, len(housesJson.Houses.HouseList))

	firstHouse := housesJson.Houses.HouseList[0]
	assert.Equal("Alai Flats, Flat 01", firstHouse.Name)
	assert.Equal(10301, firstHouse.HouseID)
	assert.Equal(17, firstHouse.Size)
	assert.Equal(50000, firstHouse.Rent)
	assert.True(firstHouse.IsRented)
	assert.False(firstHouse.IsAuctioned)
	assert.False(firstHouse.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 0, AuctionLeft: ""}, firstHouse.Auction)

	lastHouse := housesJson.Houses.HouseList[117]
	assert.Equal("Upper Swamp Lane 8", lastHouse.Name)
	assert.Equal(10405, lastHouse.HouseID)
	assert.Equal(132, lastHouse.Size)
	assert.Equal(600000, lastHouse.Rent)
	assert.True(lastHouse.IsRented)
	assert.False(lastHouse.IsAuctioned)
	assert.False(lastHouse.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 0, AuctionLeft: ""}, lastHouse.Auction)

	assert.Equal(14, len(housesJson.Houses.GuildhallList))

	firstGuild := housesJson.Houses.GuildhallList[0]
	assert.Equal("Bloodhall", firstGuild.Name)
	assert.Equal(10005, firstGuild.HouseID)
	assert.Equal(306, firstGuild.Size)
	assert.Equal(500000, firstGuild.Rent)
	assert.True(firstGuild.IsRented)
	assert.False(firstGuild.IsAuctioned)
	assert.False(firstGuild.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 0, AuctionLeft: ""}, firstGuild.Auction)

	lastGuild := housesJson.Houses.GuildhallList[13]
	assert.Equal("Warriors' Guildhall", lastGuild.Name)
	assert.Equal(10801, lastGuild.HouseID)
	assert.Equal(306, lastGuild.Size)
	assert.Equal(5000000, lastGuild.Rent)
	assert.True(lastGuild.IsRented)
	assert.False(lastGuild.IsAuctioned)
	assert.False(lastGuild.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 0, AuctionLeft: ""}, lastGuild.Auction)
}

func TestPremiaFarmineHousesOverview(t *testing.T) {
	houseData, err := os.ReadFile("../testdata/houses/overview/PremiaFarmineHouses.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	guildData, err := os.ReadFile("../testdata/houses/overview/PremiaFarmineGuilds.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	housesJson := TibiaHousesOverviewV3Impl(
		nil,
		"Premia",
		"Farmine",
		func(request TibiaDataRequestStruct) (string, error) {
			if strings.Contains(request.URL, "guildhalls") {
				return string(guildData), nil
			}

			return string(houseData), nil
		})
	assert := assert.New(t)

	assert.Equal("Premia", housesJson.Houses.World)
	assert.Equal("Farmine", housesJson.Houses.Town)

	assert.Equal(2, len(housesJson.Houses.HouseList))

	firstHouse := housesJson.Houses.HouseList[0]
	assert.Equal("Caveman Shelter", firstHouse.Name)
	assert.Equal(15001, firstHouse.HouseID)
	assert.Equal(87, firstHouse.Size)
	assert.Equal(150000, firstHouse.Rent)
	assert.False(firstHouse.IsRented)
	assert.True(firstHouse.IsAuctioned)
	assert.False(firstHouse.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 0, AuctionLeft: ""}, firstHouse.Auction)

	assert.Equal(0, len(housesJson.Houses.GuildhallList))
}

func TestPremiaEdronHousesOverview(t *testing.T) {
	houseData, err := os.ReadFile("../testdata/houses/overview/PremiaEdronHouses.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	guildData, err := os.ReadFile("../testdata/houses/overview/PremiaEdronGuilds.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	housesJson := TibiaHousesOverviewV3Impl(
		nil,
		"Premia",
		"Edron",
		func(request TibiaDataRequestStruct) (string, error) {
			if strings.Contains(request.URL, "guildhalls") {
				return string(guildData), nil
			}

			return string(houseData), nil
		})
	assert := assert.New(t)

	assert.Equal("Premia", housesJson.Houses.World)
	assert.Equal("Edron", housesJson.Houses.Town)

	assert.Equal(129, len(housesJson.Houses.HouseList))

	auctionedHouse := housesJson.Houses.HouseList[46]
	assert.Equal("Cormaya 11", auctionedHouse.Name)
	assert.Equal(54026, auctionedHouse.HouseID)
	assert.Equal(43, auctionedHouse.Size)
	assert.Equal(150000, auctionedHouse.Rent)
	assert.False(auctionedHouse.IsRented)
	assert.True(auctionedHouse.IsAuctioned)
	assert.False(auctionedHouse.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 200000, AuctionLeft: "9 hours"}, auctionedHouse.Auction)

	secondAuctionedHouse := housesJson.Houses.HouseList[56]
	assert.Equal("Cormaya 9c", secondAuctionedHouse.Name)
	assert.Equal(54023, secondAuctionedHouse.HouseID)
	assert.Equal(25, secondAuctionedHouse.Size)
	assert.Equal(80000, secondAuctionedHouse.Rent)
	assert.False(secondAuctionedHouse.IsRented)
	assert.True(secondAuctionedHouse.IsAuctioned)
	assert.True(secondAuctionedHouse.IsFinished)
	assert.Equal(HousesAuction{AuctionBid: 12345, AuctionLeft: ""}, secondAuctionedHouse.Auction)

	assert.Equal(6, len(housesJson.Houses.GuildhallList))
}
