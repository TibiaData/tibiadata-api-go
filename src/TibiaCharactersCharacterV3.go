package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	//"time"
)

// TibiaCharactersCharacterV3 func
func TibiaCharactersCharacterV3(character string) string {

	// Child of Character
	type Houses struct {
		Name    string `json:"name"`
		Town    string `json:"town"`
		Paid    string `json:"paid"`
		HouseID int    `json:"houseid"`
	}

	// Child of Character
	type Guild struct {
		GuildName string `json:"name"`
		Rank      string `json:"rank"`
	}

	// Child of Characters
	type Character struct {
		Name              string   `json:"name"`
		FormerNames       []string `json:"former_names"`
		Traded            bool     `json:"traded"`
		DeletionDate      string   `json:"deletion_date"`
		Sex               string   `json:"sex"`
		Title             string   `json:"title"`
		UnlockedTitles    int      `json:"unlocked_titles"`
		Vocation          string   `json:"vocation"`
		Level             int      `json:"level"`
		AchievementPoints int      `json:"achievement_points"`
		World             string   `json:"world"`
		FormerWorlds      []string `json:"former_worlds"`
		Residence         string   `json:"residence"`
		MarriedTo         string   `json:"married_to"`
		Houses            []Houses `json:"houses"`
		Guild             Guild    `json:"guild"`
		LastLogin         string   `json:"last_login"`
		AccountStatus     string   `json:"account_status"`
		Comment           string   `json:"comment"`
	}

	// Child of Characters
	type AccountBadges struct {
		Name        string `json:"name"`
		IconURL     string `json:"icon_url"`
		Description string `json:"description"`
	}

	// Child of Characters
	type Achievements struct {
		Name   string `json:"name"`
		Grade  int    `json:"grade"`
		Secret bool   `json:"secret"`
	}

	// Child of DeathEntries
	type Killers struct {
		Name   string `json:"name"`
		Player bool   `json:"player"`
		Traded bool   `json:"traded"`
		Summon string `json:"summon"`
	}

	// Child of Deaths
	type DeathEntries struct {
		Time    string    `json:"time"`
		Level   int       `json:"level"`
		Killers []Killers `json:"killers"`
		Assists []Killers `json:"assists"`
		Reason  string    `json:"reason"`
	}

	// Child of Characters
	type Deaths struct {
		DeathEntries    []DeathEntries `json:"death_entries"`
		TruncatedDeaths bool           `json:"truncated"` // TODO: when are those relevant..
	}

	// Child of Characters
	type AccountInformation struct {
		Position     string `json:"position"`
		Created      string `json:"created"`
		LoyaltyTitle string `json:"loyalty_title"`
	}

	// Child of Characters
	type OtherCharacters struct {
		Name    string `json:"name"`
		World   string `json:"world"`
		Status  string `json:"status"`  // online/offline
		Deleted bool   `json:"deleted"` // don't know how to do that yet..
		Main    bool   `json:"main"`
		Traded  bool   `json:"traded"`
	}

	// Child of JSONData
	type Characters struct {
		Character          Character          `json:"character"`
		AccountBadges      []AccountBadges    `json:"account_badges"`
		Achievements       []Achievements     `json:"achievements"`
		Deaths             Deaths             `json:"deaths"`
		AccountInformation AccountInformation `json:"account_information"`
		OtherCharacters    []OtherCharacters  `json:"other_characters"`
	}

	//
	// The base includes two levels, Characters and Information
	type JSONData struct {
		Characters  Characters  `json:"characters"`
		Information Information `json:"information"`
	}

	// Declaring vars for later use..
	var CharacterInformationData Character
	CharacterInformationData.Houses = []Houses{}
	CharacterInformationData.FormerNames = []string{}
	CharacterInformationData.FormerWorlds = []string{}

	var AccountBadgesData []AccountBadges
	AccountBadgesData = []AccountBadges{}

	var AchievementsData []Achievements
	AchievementsData = []Achievements{}

	var DeathsData Deaths
	DeathsData.DeathEntries = []DeathEntries{}

	var AccountInformationData AccountInformation
	var OtherCharactersData []OtherCharacters
	OtherCharactersData = []OtherCharacters{}

	var CharacterSection string

	// Sanatizing some on the character..
	character = strings.ReplaceAll(character, "+", " ")

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/community/?subtopic=characters&name=" + TibiadataQueryEscapeStringV3(character))

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

		if strings.Contains(CharacterDivHTML, "Text\">Character Information") {
			// Character Information
			CharacterSection = "characterinformation"
		} else if strings.Contains(CharacterDivHTML, "Text\">Account Badges") {
			// Account Badges
			CharacterSection = "accountbadges"
		} else if strings.Contains(CharacterDivHTML, "Text\">Account Achievements") {
			// Account Achievements
			CharacterSection = "accountachievements"
		} else if strings.Contains(CharacterDivHTML, "Text\">Character Deaths") {
			// Character Deaths
			CharacterSection = "characterdeaths"
		} else if strings.Contains(CharacterDivHTML, "Text\">Account Information") {
			// Account Information
			CharacterSection = "accountinformation"
		} else if strings.Contains(CharacterDivHTML, "Text\">Search Character") {
			// Search Character
			CharacterSection = "searchcharacter"
		} else if strings.Contains(CharacterDivHTML, "Text\">Characters") {
			// Characters
			CharacterSection = "characters"
		}

		// parsing CharacterDivHTML to goquery format
		CharacterDivQuery, _ := goquery.NewDocumentFromReader(strings.NewReader(CharacterDivHTML))

		if CharacterSection == "characterinformation" || CharacterSection == "accountinformation" {

			// Running query over each tr in character content container
			CharacterDivQuery.Find(".TableContentContainer tr").Each(func(index int, s *goquery.Selection) {

				// Storing HTML into CharacterTrHTML
				CharacterTrHTML, err := s.Html()
				if err != nil {
					log.Fatal(err)
				}

				// Removing line breaks
				CharacterTrHTML = TibiadataHTMLRemoveLinebreaksV3(CharacterTrHTML)

				// Regex to get data for fansites
				regex1 := regexp.MustCompile(`<td.*class="[a-zA-Z0-9_.-]+".*>(.*):<\/.*td><td>(.*)<\/td>`)
				subma1 := regex1.FindAllStringSubmatch(CharacterTrHTML, -1)

				if len(subma1) > 0 {
					subma1[0][2] = html.UnescapeString(subma1[0][2])

					if subma1[0][1] == "Name" {
						Tmp := strings.Split(subma1[0][2], "<")
						CharacterInformationData.Name = strings.TrimSpace(Tmp[0])
						if strings.Contains(Tmp[0], ", will be deleted at") {
							Tmp2 := strings.Split(Tmp[0], ", will be deleted at ")
							CharacterInformationData.Name = Tmp2[0]
							CharacterInformationData.DeletionDate = TibiadataDatetimeV3(strings.TrimSpace(Tmp2[1]))
						}
						if strings.Contains(subma1[0][2], "(traded)") {
							CharacterInformationData.Traded = true
							CharacterInformationData.Name = strings.Replace(CharacterInformationData.Name, " (traded)", "", -1)
						}
					} else if subma1[0][1] == "Former Names" {
						CharacterInformationData.FormerNames = strings.Split(subma1[0][2], ", ")
					} else if subma1[0][1] == "Sex" {
						CharacterInformationData.Sex = subma1[0][2]
					} else if subma1[0][1] == "Title" {
						regex1t := regexp.MustCompile(`(.*) \(([0-9]+).*`)
						subma1t := regex1t.FindAllStringSubmatch(subma1[0][2], -1)
						CharacterInformationData.Title = subma1t[0][1]
						CharacterInformationData.UnlockedTitles = TibiadataStringToIntegerV3(subma1t[0][2])
					} else if subma1[0][1] == "Vocation" {
						CharacterInformationData.Vocation = subma1[0][2]
					} else if subma1[0][1] == "Level" {
						CharacterInformationData.Level = TibiadataStringToIntegerV3(subma1[0][2])
					} else if subma1[0][1] == "Achievement Points" {
						CharacterInformationData.AchievementPoints = TibiadataStringToIntegerV3(subma1[0][2])
					} else if subma1[0][1] == "World" {
						CharacterInformationData.World = subma1[0][2]
					} else if subma1[0][1] == "Former World" {
						CharacterInformationData.FormerWorlds = strings.Split(subma1[0][2], ", ")
					} else if subma1[0][1] == "Residence" {
						CharacterInformationData.Residence = subma1[0][2]
					} else if strings.Contains(subma1[0][1], "Account") && strings.Contains(subma1[0][1], "Status") {
						// } else if subma1[0][1] == "Account Status" {
						// TODO this does not work.. somehow.. -.-
						CharacterInformationData.AccountStatus = subma1[0][2]
					} else if subma1[0][1] == "Married To" {
						CharacterInformationData.MarriedTo = TibiadataRemoveURLsV3(subma1[0][2])
					} else if subma1[0][1] == "House" {
						regex1h := regexp.MustCompile(`.*houseid=([0-9]+).*character=.*>(.*)</a> \((.*)\) is paid until (.*)`)
						subma1h := regex1h.FindAllStringSubmatch(subma1[0][2], -1)
						CharacterInformationData.Houses = append(CharacterInformationData.Houses, Houses{
							Name:    subma1h[0][2],
							Town:    subma1h[0][3],
							Paid:    TibiadataDateV3(subma1h[0][4]),
							HouseID: TibiadataStringToIntegerV3(subma1h[0][1]),
						})
					} else if strings.Contains(subma1[0][1], "Guild") && strings.Contains(subma1[0][1], "Membership") {
						// } else if subma1[0][1] == "Guild Membership" {
						// TODO this does not work.. somehow.. -.-
						Tmp := strings.Split(subma1[0][2], " of the <a href=")
						CharacterInformationData.Guild.Rank = Tmp[0]
						CharacterInformationData.Guild.GuildName = TibiadataRemoveURLsV3("<a href=" + Tmp[1])
					} else if subma1[0][1] == "Last Login" {
						if subma1[0][2] != "never logged in" {
							CharacterInformationData.LastLogin = TibiadataDatetimeV3(subma1[0][2])
						}
					} else if subma1[0][1] == "Comment" {
						CharacterInformationData.Comment = strings.ReplaceAll(subma1[0][2], "<br/>", "\n")
					} else if subma1[0][1] == "Loyalty Title" {
						AccountInformationData.LoyaltyTitle = subma1[0][2]
					} else if subma1[0][1] == "Created" {
						AccountInformationData.Created = TibiadataDatetimeV3(subma1[0][2])
					} else if subma1[0][1] == "Position" {
						TmpPosition := strings.Split(subma1[0][2], "<")
						AccountInformationData.Position = strings.TrimSpace(TmpPosition[0])
					} else {
						log.Println("LEFT OVER: `" + subma1[0][1] + "` = `" + subma1[0][2] + "`")
					}

				}

			})

		} else if CharacterSection == "accountbadges" {
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

				AccountBadgesData = append(AccountBadgesData, AccountBadges{
					Name:        subma1[0][1],
					IconURL:     subma1[0][3],
					Description: subma1[0][2],
				})
			})

		} else if CharacterSection == "accountachievements" {
			// Running query over each tr in list
			CharacterDivQuery.Find(".TableContentContainer tr").Each(func(index int, s *goquery.Selection) {

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
					subma1a[0][2] = html.UnescapeString(subma1a[0][2])

					// get the name of the achievement (and ignore the secret image on the right)
					Name := strings.Split(subma1a[0][2], "<img")

					AchievementsData = append(AchievementsData, Achievements{
						Name:   Name[0],
						Grade:  strings.Count(subma1a[0][1], "achievement-grade-symbol"),
						Secret: strings.Contains(subma1a[0][2], "achievement-secret-symbol"),
					})

				}

			})

		} else if CharacterSection == "characterdeaths" {
			// Running query over each tr in list
			CharacterDivQuery.Find(".TableContentContainer tr").Each(func(index int, s *goquery.Selection) {

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
					DeathKillers := []Killers{}
					DeathAssists := []Killers{}

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
							DeathAssists = append(DeathAssists, Killers{
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
						DeathKillers = append(DeathKillers, Killers{
							Name:   name,
							Player: isPlayer,
							Traded: isTraded,
							Summon: theSummon,
						})
					}

					// append deadentry to death list
					DeathsData.DeathEntries = append(DeathsData.DeathEntries, DeathEntries{
						Time:    TibiadataDatetimeV3(subma1[0][1]),
						Level:   TibiadataStringToIntegerV3(subma1[0][3]),
						Killers: DeathKillers,
						Assists: DeathAssists,
						Reason:  ReasonString,
					})

				}

			})

		} else if CharacterSection == "characters" {
			// Running query over each tr in character list
			CharacterDivQuery.Find(".TableContentContainer tr").Each(func(index int, s *goquery.Selection) {

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
					if strings.Contains(TmpCharacterName, " (traded)") {
						TmpTraded = true
						TmpCharacterName = strings.ReplaceAll(TmpCharacterName, " (traded)", "")
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
					OtherCharactersData = append(OtherCharactersData, OtherCharacters{
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
		Characters{
			CharacterInformationData,
			AccountBadgesData,
			AchievementsData,
			DeathsData,
			AccountInformationData,
			OtherCharactersData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	js, _ := json.Marshal(jsonData)
	if TibiadataDebug {
		fmt.Printf("%s\n", js)
	}
	return string(js)
}

// TibiaDataParseKiller func - insert a html string and get the killers back
func TibiaDataParseKiller(data string) (string, bool, bool, string) {
	var isPlayer, isTraded bool
	var theSummon string

	// check if killer is a traded player
	if strings.Contains(data, " (traded)") {
		isPlayer = true
		isTraded = true
		data = strings.ReplaceAll(data, " (traded)", "")
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
