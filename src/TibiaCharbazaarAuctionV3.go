package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaCharbazaarAuctionV3 func
func TibiaCharbazaarAuctionV3(c *gin.Context) {

	id := TibiadataStringToIntegerV3(c.Param("id"))

	// Child of Details
	type Bid struct {
		Type   string `json:"type"`
		Amount int    `json:"amount"`
	}

	// Child of Auction
	type Details struct {
		CharacterName string `json:"characterName"`
		Level         int    `json:"level"`
		Vocation      string `json:"vocation"`
		Gender        string `json:"gender"`
		World         string `json:"world"`
		AuctionStart  string `json:"auctionStart"`
		AuctionEnd    string `json:"auctionEnd"`
		Bid           Bid    `json:"bid"`
	}

	// Child of Auction
	type General struct {
		HitPoints                 int    `json:"hitPoints"`
		Mana                      int    `json:"mana"`
		Capacity                  int    `json:"capacity"`
		Speed                     int    `json:"speed"`
		Blessings                 int    `json:"blessings"`
		Mounts                    int    `json:"mounts"`
		Outfits                   int    `json:"outfits"`
		Titles                    int    `json:"titles"`
		AxeFighting               int    `json:"axeFighting"`
		ClubFighting              int    `json:"clubFighting"`
		DistanceFighting          int    `json:"distanceFighting"`
		Fishing                   int    `json:"fishing"`
		FistFighting              int    `json:"fistFighting"`
		MagicLevel                int    `json:"magicLevel"`
		Shielding                 int    `json:"shielding"`
		SwordFighting             int    `json:"swordFighting"`
		CreationDate              string `json:"creationDate"`
		Experience                int    `json:"experience"`
		Gold                      int    `json:"gold"`
		AchievementPoints         int    `json:"achievementPoints"`
		RegularWorldTransfer      string `json:"regularWorldTransfer"`
		CharmExpansion            bool   `json:"charmExpansion"`
		AvailableCharmPoints      int    `json:"availableCharmPoints"`
		SpentCharmPoints          int    `json:"spentCharmPoints"`
		DailyRewardStreak         int    `json:"dailyRewardStreak"`
		HuntingTaskPoints         int    `json:"huntingTaskPoints"`
		PermanentHuntingTaskSlots int    `json:"permanentHuntingTaskSlots"`
		PermanentPreySlots        int    `json:"permanentPreySlots"`
		PreyWildCards             int    `json:"preyWildCards"`
		Hirelings                 int    `json:"hirelings"`
		HirelingJobs              int    `json:"hirelingJobs"`
		HirelingOutfits           int    `json:"hirelingOutfits"`
		ExaltedDust               int    `json:"exaltedDust"`
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
		Name string `json:"name"`
		Cost int    `json:"cost"`
	}

	// Child of Auction
	type BestiaryEntry struct {
		Name  string `json:"name"`
		Kills int    `json:"kills"`
		Step  int    `json:"step"`
	}

	// Child of JSONData
	type Auction struct {
		Id                          int             `json:"id"`
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
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/charactertrade/?page=details&auctionid=" +
		strconv.Itoa(id))

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Extract details section
	var details Details
	ReaderHTML.Find(".Auction").Each(func(index int, s *goquery.Selection) {
		detailsHeader := strings.Split(s.Find(".AuctionHeader").Text(), "Level: ")
		details.CharacterName = detailsHeader[0]

		detailsHeader = strings.Split(detailsHeader[1], "|")

		level := TibiadataStringToIntegerV3(detailsHeader[0])

		details.Level = level
		details.Vocation = strings.TrimSpace(strings.Split(detailsHeader[1], "Vocation: ")[1])
		details.Gender = strings.TrimSpace(detailsHeader[2])
		details.World = strings.TrimSpace(strings.Split(detailsHeader[3], "World: ")[1])

		s.Find(".ShortAuctionData").Each(func(index int, s *goquery.Selection) {
			nodes := s.Children().Nodes

			lookupIndex := 0
			hasTimer := strings.EqualFold(nodes[0].Attr[0].Val, "AuctionTimer")

			// In case the auction hasn't ended yet, we need to increase the lookup index.
			if hasTimer {
				lookupIndex = 1
			}

			auctionStartDate := TibiaDataSanitizeNbspSpaceString(nodes[1+lookupIndex].FirstChild.Data)
			auctionStartDate = strings.Split(auctionStartDate, " CET")[0] + ":00 CET"

			auctionEndDate := TibiaDataSanitizeNbspSpaceString(nodes[3+lookupIndex].FirstChild.Data)
			auctionEndDate = strings.Split(auctionEndDate, " CET")[0] + ":00 CET"

			details.AuctionStart = TibiadataDatetimeV3(auctionStartDate)
			details.AuctionEnd = TibiadataDatetimeV3(auctionEndDate)

			bidType := strings.Split(nodes[4+lookupIndex].FirstChild.FirstChild.Data, " Bid:")[0]
			bidAmount := TibiadataStringToIntegerV3(nodes[4+lookupIndex].LastChild.FirstChild.FirstChild.Data)

			details.Bid = Bid{
				Type:   bidType,
				Amount: bidAmount,
			}
		})
	})

	// Extract general section
	var general General
	ReaderHTML.Find("#General").Each(func(index int, s *goquery.Selection) {

		// General
		generalMap := make(map[string]string)
		s.Find(".LabelV").Each(func(index int, s *goquery.Selection) {
			generalMap[strings.Split(s.Nodes[0].FirstChild.Data, ":")[0]] = s.Nodes[0].NextSibling.FirstChild.Data
		})

		// Skills
		skillsMap := make(map[string]int)
		s.Find(".LabelColumn").Each(func(index int, s *goquery.Selection) {
			skillsMap[strings.Split(s.Nodes[0].FirstChild.FirstChild.Data, ":")[0]] =
				TibiadataStringToIntegerV3(s.Nodes[0].NextSibling.FirstChild.Data)
		})

		general.HitPoints = TibiadataStringToIntegerV3(generalMap["Hit Points"])
		general.Mana = TibiadataStringToIntegerV3(generalMap["Mana"])
		general.Capacity = TibiadataStringToIntegerV3(generalMap["Capacity"])
		general.Speed = TibiadataStringToIntegerV3(generalMap["Speed"])
		general.Blessings = TibiadataStringToIntegerV3(strings.Split(generalMap["Blessings"], "/")[0])
		general.Mounts = TibiadataStringToIntegerV3(generalMap["Mounts"])
		general.Outfits = TibiadataStringToIntegerV3(generalMap["Outfits"])
		general.Titles = TibiadataStringToIntegerV3(generalMap["Titles"])
		general.AxeFighting = skillsMap["Axe Fighting"]
		general.ClubFighting = skillsMap["Club Fighting"]
		general.DistanceFighting = skillsMap["Distance Fighting"]
		general.Fishing = skillsMap["Fishing"]
		general.FistFighting = skillsMap["Fist Fighting"]
		general.MagicLevel = skillsMap["Magic Level"]
		general.Shielding = skillsMap["Shielding"]
		general.SwordFighting = skillsMap["Sword Fighting"]
		general.CreationDate = TibiadataDatetimeV3(generalMap["Creation Date"])
		general.Experience = TibiadataStringToIntegerV3(generalMap["Experience"])
		general.Gold = TibiadataStringToIntegerV3(generalMap["Gold"])
		general.AchievementPoints = TibiadataStringToIntegerV3(generalMap["Achievement Points"])
		general.RegularWorldTransfer = generalMap["Regular World Transfer"]
		general.CharmExpansion = strings.EqualFold(generalMap["Charm Expansion"], "yes")
		general.AvailableCharmPoints = TibiadataStringToIntegerV3(generalMap["Available Charm Points"])
		general.SpentCharmPoints = TibiadataStringToIntegerV3(generalMap["Spent Charm Points"])
		general.DailyRewardStreak = TibiadataStringToIntegerV3(generalMap["Daily Reward Streak"])
		general.HuntingTaskPoints = TibiadataStringToIntegerV3(generalMap["Hunting Task Points"])
		general.PermanentHuntingTaskSlots = TibiadataStringToIntegerV3(generalMap["Permanent Hunting Task Slots"])
		general.PermanentPreySlots = TibiadataStringToIntegerV3(generalMap["Permanent Prey Slots"])
		general.PreyWildCards = TibiadataStringToIntegerV3(generalMap["Prey Wildcards"])
		general.Hirelings = TibiadataStringToIntegerV3(generalMap["Hirelings"])
		general.HirelingJobs = TibiadataStringToIntegerV3(generalMap["Hireling Jobs"])
		general.HirelingOutfits = TibiadataStringToIntegerV3(generalMap["Hireling Outfits"])
		general.ExaltedDust = TibiadataStringToIntegerV3(strings.Split(generalMap["Exalted Dust"], "/")[0])
	})

	// Extract items summary
	var itemSummary []Item
	ReaderHTML.Find("#ItemSummary").Each(func(index int, s *goquery.Selection) {

		for k, v := range ParseItems(s) {
			itemSummary = append(itemSummary, Item{Name: k, Amount: v})
		}

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				itemPage := AjaxJSONDataCollectorV3(id, ItemSummarySection, pageIndex)
				ItemPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(itemPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseItems(ItemPageReaderHTML.Contents()) {
					itemSummary = append(itemSummary, Item{Name: k, Amount: v})
				}
			}
		}
	})

	// Extract store items summary
	var storeItemSummary []Item
	ReaderHTML.Find("#StoreItemSummary").Each(func(index int, s *goquery.Selection) {

		for k, v := range ParseItems(s) {
			storeItemSummary = append(storeItemSummary, Item{Name: k, Amount: v})
		}

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				itemPage := AjaxJSONDataCollectorV3(id, StoreItemSummarySection, pageIndex)
				ItemPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(itemPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseItems(ItemPageReaderHTML.Contents()) {
					storeItemSummary = append(storeItemSummary, Item{Name: k, Amount: v})
				}
			}
		}
	})

	// Extract mounts
	var mounts []string
	ReaderHTML.Find("#Mounts").Each(func(index int, s *goquery.Selection) {

		mounts = append(mounts, ParseMounts(s)...)

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				mountsPage := AjaxJSONDataCollectorV3(id, MountsSection, pageIndex)
				MountsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(mountsPage))
				if err != nil {
					log.Fatal(err)
				}
				mounts = append(mounts, ParseMounts(MountsPageReaderHTML.Contents())...)
			}
		}
	})

	// Extract store mounts
	var storeMounts []string
	ReaderHTML.Find("#StoreMounts").Each(func(index int, s *goquery.Selection) {

		storeMounts = append(storeMounts, ParseMounts(s)...)

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				mountsPage := AjaxJSONDataCollectorV3(id, StoreMountsSection, pageIndex)
				MountsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(mountsPage))
				if err != nil {
					log.Fatal(err)
				}
				storeMounts = append(storeMounts, ParseMounts(MountsPageReaderHTML.Contents())...)
			}
		}
	})

	// Extract outfits
	var outfits []Outfit
	ReaderHTML.Find("#Outfits").Each(func(index int, s *goquery.Selection) {
		for k, v := range ParseOutfits(s) {
			outfits = append(outfits, Outfit{
				Name:   k,
				Addon1: v[0],
				Addon2: v[1],
			})
		}
		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				outfitsPage := AjaxJSONDataCollectorV3(id, OutfitsSection, pageIndex)
				OutfitsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(outfitsPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseOutfits(OutfitsPageReaderHTML.Contents()) {
					outfits = append(outfits, Outfit{
						Name:   k,
						Addon1: v[0],
						Addon2: v[1],
					})
				}
			}
		}
	})

	// Extract store outfits
	var storeOutfits []Outfit
	ReaderHTML.Find("#StoreOutfits").Each(func(index int, s *goquery.Selection) {
		for k, v := range ParseOutfits(s) {
			outfits = append(outfits, Outfit{
				Name:   k,
				Addon1: v[0],
				Addon2: v[1],
			})
		}
		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				outfitsPage := AjaxJSONDataCollectorV3(id, StoreOutfitsSection, pageIndex)
				OutfitsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(outfitsPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseOutfits(OutfitsPageReaderHTML.Contents()) {
					outfits = append(outfits, Outfit{
						Name:   k,
						Addon1: v[0],
						Addon2: v[1],
					})
				}
			}
		}
	})

	// Extract familiars
	var familiars []string
	ReaderHTML.Find("#Familiars").Each(func(index int, s *goquery.Selection) {
		s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
			familiarName, exists := s.Attr("title")
			if exists {
				familiars = append(familiars, familiarName)
			}
		})
	})

	// Extract blessings
	var blessings Blessings
	ReaderHTML.Find("#Blessings").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			blessingsAmount := TibiadataStringToIntegerV3(strings.Split(node.FirstChild.Data, " x")[0])
			switch blessingName := node.NextSibling.FirstChild.Data; blessingName {
			case "Adventurer's Blessing":
				blessings.AdventurersBlessing = blessingsAmount
			case "Blood of the Mountain":
				blessings.BloodOfTheMountain = blessingsAmount
			case "Embrace of Tibia":
				blessings.EmbraceOfTibia = blessingsAmount
			case "Fire of the Suns":
				blessings.FireOfTheSuns = blessingsAmount
			case "Heart of the Mountain":
				blessings.HeartOfTheMountain = blessingsAmount
			case "Spark of the Phoenix":
				blessings.SparkOfThePhoenix = blessingsAmount
			case "Spiritual Shielding":
				blessings.SpiritualShielding = blessingsAmount
			case "Twist of Fate":
				blessings.TwistOfFate = blessingsAmount
			case "Wisdom of Solitude":
				blessings.WisdomOfSolitude = blessingsAmount
			}
		})
	})

	// Extract Imbuements
	var imbuements []string
	ReaderHTML.Find("#Imbuements").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No imbuements.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				imbuements = append(imbuements, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Charms
	var charms []Charm
	ReaderHTML.Find("#Charms").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No charms.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				charms = append(charms, Charm{
					Cost: TibiadataStringToIntegerV3(node.FirstChild.Data),
					Name: node.NextSibling.FirstChild.Data,
				})
			}
		})
	})

	// Extract Completed Cyclopedia Map Areas
	var completedCyclopediaMapAreas []string
	ReaderHTML.Find("#CompletedCyclopediaMapAreas").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No areas explored.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				completedCyclopediaMapAreas = append(completedCyclopediaMapAreas, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Completed Quest Lines
	var completedQuestLines []string
	ReaderHTML.Find("#CompletedQuestLines").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				completedQuestLines = append(completedQuestLines, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Titles
	var titles []string
	ReaderHTML.Find("#Titles").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				titles = append(titles, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Achievements
	var achievements []string
	ReaderHTML.Find("#Achievements").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				achievements = append(achievements, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Bestiary Progress
	var bestiaryProgress []BestiaryEntry
	ReaderHTML.Find("#BestiaryProgress").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				bestiaryProgress = append(bestiaryProgress, BestiaryEntry{
					Step:  TibiadataStringToIntegerV3(node.FirstChild.Data),
					Kills: TibiadataStringToIntegerV3(strings.Split(node.NextSibling.FirstChild.Data, " x")[0]),
					Name:  node.NextSibling.NextSibling.FirstChild.Data,
				})
			}
		})
	})

	jsonData := JSONData{
		Auction: Auction{
			Id:                          id,
			Details:                     details,
			General:                     general,
			ItemSummary:                 itemSummary,
			StoreItemSummary:            storeItemSummary,
			Mounts:                      mounts,
			StoreMounts:                 storeMounts,
			Outfits:                     outfits,
			StoreOutfits:                storeOutfits,
			Familiars:                   familiars,
			Blessings:                   blessings,
			Imbuements:                  imbuements,
			Charms:                      charms,
			CompletedCyclopediaMapAreas: completedCyclopediaMapAreas,
			CompletedQuestLines:         completedQuestLines,
			Titles:                      titles,
			Achievements:                achievements,
			BestiaryProgress:            bestiaryProgress,
		},
		Information: Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaCharbazaarAuctionV3", jsonData)
}

func ParseItems(s *goquery.Selection) map[string]int {
	m := make(map[string]int)
	s.Find(".CVIconObject").Each(func(index int, s *goquery.Selection) {

		itemTitle, exists := s.Attr("title")

		if exists {
			var itemAmount int
			var itemName string

			nodes := s.Find(".ObjectAmount").First().Nodes
			if nodes == nil {
				itemName = strings.Split(itemTitle, "\n")[0]
				itemAmount = 1
			} else {
				temp := strings.SplitN(itemTitle, "x ", 2)
				itemName = strings.Split(temp[1], "\n")[0]
				itemAmount = TibiadataStringToIntegerV3(temp[0])
			}
			m[itemName] = itemAmount
		}
	})
	return m
}

func ParseOutfits(s *goquery.Selection) map[string][]bool {
	m := make(map[string][]bool)
	s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
		outfitTitle, exists := s.Attr("title")
		if exists {
			outfitName := strings.Split(outfitTitle, " (")[0]
			hasAddon1 := strings.Contains(outfitTitle, "addon 1")
			hasAddon2 := strings.Contains(outfitTitle, "addon 2")
			m[outfitName] = []bool{hasAddon1, hasAddon2}
		}
	})
	return m
}

func ParseMounts(s *goquery.Selection) []string {
	var mountsList []string
	s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
		mountTitle, exists := s.Attr("title")
		if exists {
			mountsList = append(mountsList, mountTitle)
		}
	})
	return mountsList
}

func AjaxJSONDataCollectorV3(AuctionId int, SectionType int, PageIndex int) string {
	TibiaURL := "https://www.tibia.com/websiteservices/handle_charactertrades.php?" +
		"auctionid=" + strconv.Itoa(AuctionId) +
		"&type=" + strconv.Itoa(SectionType) +
		"&currentpage=" + strconv.Itoa(PageIndex)

	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiadataDebug var
	if TibiadataDebug {
		client.SetDebug(true)
	}

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Content-Type":     "application/json",
		"User-Agent":       TibiadataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	res, err := client.R().Get(TibiaURL)
	if err != nil {
		log.Printf("[error] AjaxJSONDataCollectorV3 (URL: %s) in resp1: %s", TibiaURL, err)
	}

	// Checking if status is something else than 200
	if res.StatusCode() != 200 {
		log.Printf("[warni] AjaxJSONDataCollectorV3 (URL: %s) status code: %s", TibiaURL, res.Status())

		// Check if page is in maintenance mode
		if res.StatusCode() == 302 {
			log.Printf("[info] AjaxJSONDataCollectorV3 (URL: %s): Page tibia.com returns 302, probably maintenance mode enabled?", TibiaURL)
		}
	}

	var result AjaxJSONData
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		log.Printf("[error] AjaxJSONDataCollectorV3 (URL: %s) in deserialization process: %s", TibiaURL, err)
	}

	// Return of extracted html to functions
	return result.AjaxObjects[0].Data
}

// AjaxResponseObject - child of AjaxJSONData
type AjaxResponseObject struct {
	Data     string `json:"Data"`
	DataType string `json:"DataType"`
	Target   string `json:"Target"`
}

// AjaxJSONData - base response for auction items page links
type AjaxJSONData struct {
	AjaxObjects []AjaxResponseObject `json:"AjaxObjects"`
}

const OddEvenSelector = ".Odd,.Even"
const PageLinkSelector = ".PageLink"
const CVIconSelector = ".CVIcon"
const ItemSummarySection = 0
const StoreItemSummarySection = 1
const MountsSection = 2
const StoreMountsSection = 3
const OutfitsSection = 4
const StoreOutfitsSection = 5
