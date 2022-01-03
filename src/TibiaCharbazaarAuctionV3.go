package main

import (
	"fmt"
	"log"
	//"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaCharbazaarAuctionV3 func
func TibiaCharbazaarAuctionV3(c *gin.Context) {

	id := c.Param("id")

	// Child of Auction
	type Details struct {
		CharacterName string `json:"characterName"`
		Level         int    `json:"level"`
		Vocation      string `json:"vocation"`
		Gender        string `json:"gender"`
		World         string `json:"world"`

		AuctionStart string `json:"auctionStart"`
		AuctionEnd   string `json:"auctionEnd"`
		MinimumBid   int    `json:"minimumBid"`
		WinningBid   int    `json:"winningBid"`
	}

	// Child of Auction
	type General struct {
		HitPoints int    `json:"hitPoints"`
		Mana      int    `json:"mana"`
		Capacity  int    `json:"capacity"`
		Speed     int    `json:"speed"`
		Blessings string `json:"blessings"`
		Mounts    int    `json:"mounts"`
		Outfits   int    `json:"outfits"`
		Titles    int    `json:"titles"`

		AxeFighting      int `json:"axeFighting"`
		ClubFighting     int `json:"clubFighting"`
		DistanceFighting int `json:"distanceFighting"`
		Fishing          int `json:"fishing"`
		FistFighting     int `json:"fistFighting"`
		Shielding        int `json:"shielding"`
		SwordFighting    int `json:"swordFighting"`

		CreationDate      string `json:"creationDate"`
		Experience        int    `json:"experience"`
		Gold              int    `json:"gold"`
		AchievementPoints int    `json:"achievementPoints"`

		RegularWorldTransfer string `json:"regularWorldTransfer"`

		CharmExpansion       bool `json:"charmExpansion"`
		AvailableCharmPoints int  `json:"availableCharmPoints"`
		SpentCharmPoints     int  `json:"spentCharmPoints"`

		DailyRewardStreak int `json:"dailyRewardStreak"`

		HuntingTaskPoints         int `json:"huntingTaskPoints"`
		PermanentHuntingTaskSlots int `json:"permanentHuntingTaskSlots"`
		PermanentPreySlots        int `json:"permanentPreySlots"`
		PreyWildCards             int `json:"preyWildCards"`

		Hirelings       int `json:"hirelings"`
		HirelingJobs    int `json:"hirelingJobs"`
		HirelingOutfits int `json:"hirelingOutfits"`

		ExaltedDust int `json:"exaltedDust"`
	}

	// Child of Auction
	type Item struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	// Child of Auction
	type Outfit struct {
		Name   string `json:"name"`
		Addon1 bool   `json:"addon1"`
		Addon2 bool   `json:"addon2"`
	}

	// Child of Auction
	type Blessings struct {
		AdventurersBlessing int `json:"adventurersBlessing"`
		BloodOfTheMountain  int `json:"bloodOfTheMountain"`
		EmbraceOfTibia      int `json:"embraceOfTibia"`
		FireOfTheSuns       int `json:"fireOfTheSuns"`
		HeartOfTheMountain  int `json:"heartOfTheMountain"`
		SparkOfThePhoenix   int `json:"sparkOfThePhoenix"`
		SpiritualShielding  int `json:"spiritualShielding"`
		TwistOfFate         int `json:"twistOfFate"`
		WisdomOfSolitude    int `json:"wisdomOfSolitude"`
	}

	// Child of Auction
	type Charm struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	// Child of Auction
	type BestiaryEntry struct {
		Name  string `json:"name"`
		Kills int    `json:"kills"`
		Step  int    `json:"step"`
	}

	// Child of JSONData
	type Auction struct {
		Id                          string          `json:"id"`
		Details                     Details         `json:"details"`
		General                     General         `json:"general"`
		ItemSummary                 []Item          `json:"itemSummary"`
		StoreItemSummary            []Item          `json:"storeItemSummary"`
		Mounts                      []string        `json:"mounts"`
		StoreMounts                 []string        `json:"storeMounts"`
		Outfits                     []Outfit        `json:"outfits"`
		StoreOutfits                []Outfit        `json:"storeOutfits"`
		Familiars                   []string        `json:"familiars"`
		Blessings                   Blessings       `json:"blessings"`
		Imbuements                  []string        `json:"imbuements"`
		Charms                      []Charm         `json:"charms"`
		CompletedCyclopediaMapAreas []string        `json:"completedCyclopediaMapAreas"`
		CompletedQuestLines         []string        `json:"completedQuestLines"`
		Titles                      []string        `json:"titles"`
		Achievements                []string        `json:"achievements"`
		BestiaryProgress            []BestiaryEntry `json:"bestiaryProgress"`
	}

	// The base includes two levels: Auction and Information
	type JSONData struct {
		Auction     Auction     `json:"auction"`
		Information Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/charactertrade/?page=details&auctionid=" + id)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Extract AuctionDetails
	var AuctionDetails Details
	ReaderHTML.Find("#General div.TableContentContainer table.TableContent td div").NextAll().Each(func(index int, s *goquery.Selection) {
		AuctionSection, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		LocalReader, err := goquery.NewDocumentFromReader(strings.NewReader(AuctionSection))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(AuctionSection)
	})

	ReaderHTML.Find(".Auction").First().Each(func(index int, s *goquery.Selection) {
		AuctionCharacterName, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(AuctionCharacterName)
	})

	jsonData := JSONData{
		Auction: Auction{
			Id: id,
			Details: Details{
				CharacterName: "characterName",
				Level:         170,
				Vocation:      "",
				Gender:        "",
				World:         "",
				AuctionStart:  "",
				AuctionEnd:    "",
				MinimumBid:    99,
				WinningBid:    999,
			},
			General:          General{},
			ItemSummary:      []Item{},
			StoreItemSummary: []Item{},
		},
		Information: Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaCharbazaarAuctionV3", jsonData)
}
