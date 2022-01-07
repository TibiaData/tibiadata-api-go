package main

import (
	"log"
	"regexp"
	"strings"
	"tibiadata-api-go/src/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	//"time"
)

// TibiaCharactersCharacterV3 func
func TibiaCharactersCharacterV3(c *gin.Context) {
	var (
		// local strings used in this function
		localDivQueryString = ".TableContentContainer tr"
		localTradedString   = " (traded)"

		// Declaring vars for later use..
		CharacterInformationData structs.Character
		AccountBadgesData        []structs.AccountBadges
		AchievementsData         []structs.Achievements
		DeathsData               structs.Deaths
		AccountInformationData   structs.AccountInformation
		OtherCharactersData      []structs.OtherCharacters

		CharacterSection string
	)

	// getting params from URL
	character := c.Param("character")

	// The base includes two levels, Characters and Information
	type JSONData struct {
		Characters  structs.Characters  `json:"characters"`
		Information structs.Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=characters&name=" + TibiadataQueryEscapeStringV3(character)
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Running query on every .TableContainer
	ReaderHTML.Find(".TableContainer").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CharacterDivHTML
		CharacterDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case strings.Contains(CharacterDivHTML, "Text\">Character Information"):
			// Character Information
			CharacterSection = "characterinformation"
		case strings.Contains(CharacterDivHTML, "Text\">Account Badges"):
			// Account Badges
			CharacterSection = "accountbadges"
		case strings.Contains(CharacterDivHTML, "Text\">Account Achievements"):
			// Account Achievements
			CharacterSection = "accountachievements"
		case strings.Contains(CharacterDivHTML, "Text\">Character Deaths"):
			// Character Deaths
			CharacterSection = "characterdeaths"
		case strings.Contains(CharacterDivHTML, "Text\">Account Information"):
			// Account Information
			CharacterSection = "accountinformation"
		case strings.Contains(CharacterDivHTML, "Text\">Search Character"):
			// Search Character
			CharacterSection = "searchcharacter"
		case strings.Contains(CharacterDivHTML, "Text\">Characters"):
			// Characters
			CharacterSection = "characters"
		}

		// parsing CharacterDivHTML to goquery format
		CharacterDivQuery, _ := goquery.NewDocumentFromReader(strings.NewReader(CharacterDivHTML))

		switch CharacterSection {
		case "characterinformation", "accountinformation":
			// Running query over each tr in character content container
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				// Storing HTML into CharacterTrHTML
				CharacterTrHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterTrHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterTrHTML)
				// Unescape hmtl string
				CharacterTrHTML = TibiaDataSanitizeEscapedString(CharacterTrHTML)

				// Regex to get data for fansites
				regex1 := regexp.MustCompile(`<td.*class=.[a-zA-Z0-9_.-]+..*>(.*):<\/.*td><td>(.*)<\/td>`)
				subma1 := regex1.FindAllStringSubmatch(CharacterTrHTML, -1)

				if len(subma1) > 0 {
					switch TibiaDataSanitizeNbspSpaceString(subma1[0][1]) {
					case "Name":
						Tmp := strings.Split(subma1[0][2], "<")
						CharacterInformationData.Name = strings.TrimSpace(Tmp[0])
						if strings.Contains(Tmp[0], ", will be deleted at") {
							Tmp2 := strings.Split(Tmp[0], ", will be deleted at ")
							CharacterInformationData.Name = Tmp2[0]
							CharacterInformationData.DeletionDate = TibiadataDatetimeV3(strings.TrimSpace(Tmp2[1]))
						}
						if strings.Contains(subma1[0][2], localTradedString) {
							CharacterInformationData.Traded = true
							CharacterInformationData.Name = strings.Replace(CharacterInformationData.Name, localTradedString, "", -1)
						}
					case "Former Names":
						CharacterInformationData.FormerNames = strings.Split(subma1[0][2], ", ")
					case "Sex":
						CharacterInformationData.Sex = subma1[0][2]
					case "Title":
						regex1t := regexp.MustCompile(`(.*) \(([0-9]+).*`)
						subma1t := regex1t.FindAllStringSubmatch(subma1[0][2], -1)
						CharacterInformationData.Title = subma1t[0][1]
						CharacterInformationData.UnlockedTitles = TibiadataStringToIntegerV3(subma1t[0][2])
					case "Vocation":
						CharacterInformationData.Vocation = subma1[0][2]
					case "Level":
						CharacterInformationData.Level = TibiadataStringToIntegerV3(subma1[0][2])
					case "Achievement Points":
						CharacterInformationData.AchievementPoints = TibiadataStringToIntegerV3(subma1[0][2])
					case "World":
						CharacterInformationData.World = subma1[0][2]
					case "Former World":
						CharacterInformationData.FormerWorlds = strings.Split(subma1[0][2], ", ")
					case "Residence":
						CharacterInformationData.Residence = subma1[0][2]
					case "Account Status":
						CharacterInformationData.AccountStatus = subma1[0][2]
					case "Married To":
						CharacterInformationData.MarriedTo = TibiadataRemoveURLsV3(subma1[0][2])
					case "House":
						regex1h := regexp.MustCompile(`.*houseid=([0-9]+).*character=.*>(.*)</a> \((.*)\) is paid until (.*)`)
						subma1h := regex1h.FindAllStringSubmatch(subma1[0][2], -1)
						CharacterInformationData.Houses = append(CharacterInformationData.Houses, structs.CharacterHouses{
							Name:    subma1h[0][2],
							Town:    subma1h[0][3],
							Paid:    TibiadataDateV3(subma1h[0][4]),
							HouseID: TibiadataStringToIntegerV3(subma1h[0][1]),
						})
					case "Guild Membership":
						Tmp := strings.Split(subma1[0][2], " of the <a href=")
						CharacterInformationData.Guild.Rank = Tmp[0]
						CharacterInformationData.Guild.GuildName = TibiadataRemoveURLsV3("<a href=" + Tmp[1])
					case "Last Login":
						if subma1[0][2] != "never logged in" {
							CharacterInformationData.LastLogin = TibiadataDatetimeV3(subma1[0][2])
						}
					case "Comment":
						CharacterInformationData.Comment = strings.ReplaceAll(subma1[0][2], "<br/>", "\n")
					case "Loyalty Title":
						AccountInformationData.LoyaltyTitle = subma1[0][2]
					case "Created":
						AccountInformationData.Created = TibiadataDatetimeV3(subma1[0][2])
					case "Position":
						TmpPosition := strings.Split(subma1[0][2], "<")
						AccountInformationData.Position = strings.TrimSpace(TmpPosition[0])
					default:
						log.Println("LEFT OVER: `" + subma1[0][1] + "` = `" + subma1[0][2] + "`")
					}
				}
			})
		case "accountbadges":
			// Running query over each tr in list
			CharacterDivQuery.Find(".TableContentContainer tr td").Each(func(index int, s *goquery.Selection) {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterListHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterListHTML)

				// Regex to get account badges of account
				regex1 := regexp.MustCompile(`\(this\), &#39;(.*)&#39;, &#39;(.*)&#39;,.*\).*src="(.*)" alt=.*`)
				subma1 := regex1.FindAllStringSubmatch(CharacterListHTML, -1)

				AccountBadgesData = append(AccountBadgesData, structs.AccountBadges{
					Name:        subma1[0][1],
					IconURL:     subma1[0][3],
					Description: subma1[0][2],
				})
			})
		case "accountachievements":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterListHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterListHTML)

				// Regex to get other characters of same account
				regex1a := regexp.MustCompile(`<td class="[a-zA-Z0-9_.-]+">(.*)<\/td><td>(.*)?<?.*<\/td>`)
				subma1a := regex1a.FindAllStringSubmatch(CharacterListHTML, -1)
				if len(subma1a) > 0 {

					// fixing encoding for achievement name
					subma1a[0][2] = TibiaDataSanitizeEscapedString(subma1a[0][2])

					// get the name of the achievement (and ignore the secret image on the right)
					Name := strings.Split(subma1a[0][2], "<img")

					AchievementsData = append(AchievementsData, structs.Achievements{
						Name:   Name[0],
						Grade:  strings.Count(subma1a[0][1], "achievement-grade-symbol"),
						Secret: strings.Contains(subma1a[0][2], "achievement-secret-symbol"),
					})
				}
			})
		case "characterdeaths":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterListHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterListHTML)
				//log.Println(CharacterListHTML)
				CharacterListHTML = strings.ReplaceAll(CharacterListHTML, ".<br/>Assisted by", ". Assisted by")

				// Regex to get data for deaths
				regex1 := regexp.MustCompile(`<td.*>(.*)<\/td><td>(.*) at Level ([0-9]+) by (.*).<\/td>`)
				subma1 := regex1.FindAllStringSubmatch(CharacterListHTML, -1)

				if len(subma1) > 0 {
					// defining responses
					DeathKillers := []structs.Killers{}
					DeathAssists := []structs.Killers{}

					// store for reply later on..
					ReasonString := RemoveHtmlTag(subma1[0][2] + " at Level " + subma1[0][3] + " by " + subma1[0][4] + ".")

					// if kill is with assist..
					if strings.Contains(subma1[0][4], ". Assisted by ") {
						TmpListOfDeath := strings.Split(subma1[0][4], ". Assisted by ")
						subma1[0][4] = TmpListOfDeath[0]
						TmpAssist := TmpListOfDeath[1]

						// get a list of killers
						ListOfAssists := strings.Split(TmpAssist, ", ")

						// extract if "and" is in last ss1
						ListOfAssistsTmp := strings.Split(ListOfAssists[len(ListOfAssists)-1], " and ")

						// if there is an "and", then we split it..
						if len(ListOfAssistsTmp) > 1 {
							ListOfAssists[len(ListOfAssists)-1] = ListOfAssistsTmp[0]
							ListOfAssists = append(ListOfAssists, ListOfAssistsTmp[1])
						}

						// loop through all killers and append to result
						for i := range ListOfAssists {
							name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfAssists[i])
							DeathAssists = append(DeathAssists, structs.Killers{
								Name:   name,
								Player: isPlayer,
								Traded: isTraded,
								Summon: theSummon,
							})
						}
					}

					// get a list of killers
					ListOfKillers := strings.Split(subma1[0][4], ", ")

					// extract if "and" is in last ss1
					ListOfKillersTmp := strings.Split(ListOfKillers[len(ListOfKillers)-1], " and ")

					// if there is an "and", then we split it..
					if len(ListOfKillersTmp) > 1 {
						ListOfKillers[len(ListOfKillers)-1] = ListOfKillersTmp[0]
						ListOfKillers = append(ListOfKillers, ListOfKillersTmp[1])
					}

					// loop through all killers and append to result
					for i := range ListOfKillers {
						name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfKillers[i])
						DeathKillers = append(DeathKillers, structs.Killers{
							Name:   name,
							Player: isPlayer,
							Traded: isTraded,
							Summon: theSummon,
						})
					}

					// append deadentry to death list
					DeathsData.DeathEntries = append(DeathsData.DeathEntries, structs.DeathEntries{
						Time:    TibiadataDatetimeV3(subma1[0][1]),
						Level:   TibiadataStringToIntegerV3(subma1[0][3]),
						Killers: DeathKillers,
						Assists: DeathAssists,
						Reason:  ReasonString,
					})
				}
			})
		case "characters":
			// Running query over each tr in character list
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterListHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterListHTML)

				// Regex to get data for fansites
				regex1 := regexp.MustCompile(`<td.*<nobr>[0-9]+\..(.*)<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td><td style="width: 70%">(.*)<\/td><td.*`)
				subma1 := regex1.FindAllStringSubmatch(CharacterListHTML, -1)

				if len(subma1) > 0 {
					TmpCharacterName := subma1[0][1]

					var TmpTraded bool
					if strings.Contains(TmpCharacterName, localTradedString) {
						TmpTraded = true
						TmpCharacterName = strings.ReplaceAll(TmpCharacterName, localTradedString, "")
					}

					// If this character is the main character of the account
					TmpMain := false
					if strings.Contains(TmpCharacterName, "Main Character") {
						TmpMain = true
						Tmp := strings.Split(subma1[0][1], "<")
						TmpCharacterName = strings.TrimSpace(Tmp[0])
					}

					// If this character is online or offline
					TmpStatus := "offline"
					if strings.Contains(subma1[0][3], "<b class=\"green\">online</b>") {
						TmpStatus = "online"
					}

					// Is this character is deleted
					TmpDeleted := false
					if strings.Contains(subma1[0][3], "deleted") {
						TmpDeleted = true
					}

					// Create the character and append it to the other characters list
					OtherCharactersData = append(OtherCharactersData, structs.OtherCharacters{
						Name:    TmpCharacterName,
						World:   subma1[0][2],
						Status:  TmpStatus,
						Deleted: TmpDeleted,
						Main:    TmpMain,
						Traded:  TmpTraded,
					})
				}
			})
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		structs.Characters{
			Character:          CharacterInformationData,
			AccountBadges:      AccountBadgesData,
			Achievements:       AchievementsData,
			Deaths:             DeathsData,
			AccountInformation: AccountInformationData,
			OtherCharacters:    OtherCharactersData,
		},
		structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaCharactersCharacterV3", jsonData)
}

// TibiaDataParseKiller func - insert a html string and get the killers back
func TibiaDataParseKiller(data string) (string, bool, bool, string) {
	// local strings used in this function
	var (
		localTradedString  = " (traded)"
		isPlayer, isTraded bool
		theSummon          string
	)

	// check if killer is a traded player
	if strings.Contains(data, localTradedString) {
		isPlayer = true
		isTraded = true
		data = strings.ReplaceAll(data, localTradedString, "")
	}

	// check if killer is a player
	if strings.Contains(data, "https://www.tibia.com") {
		isPlayer = true
		data = RemoveHtmlTag(data)
	}

	// get summon information
	re := regexp.MustCompile(`(an? .+) of ([^<]+)`)
	rs := re.FindAllStringSubmatch(data, -1)
	if len(rs) >= 1 {
		theSummon = rs[0][1]
		data = rs[0][2]
	}

	return data, isPlayer, isTraded, theSummon
}
