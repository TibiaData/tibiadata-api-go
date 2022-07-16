package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TibiaData/tibiadata-api-go/src/validation"
	//"time"
)

// Child of Character
type Houses struct {
	Name    string `json:"name"`
	Town    string `json:"town"`
	Paid    string `json:"paid"`
	HouseID int    `json:"houseid"`
}

// Child of Character
type CharacterGuild struct {
	GuildName string `json:"name,omitempty"`
	Rank      string `json:"rank,omitempty"`
}

// Child of Characters
type Character struct {
	Name              string         `json:"name"`
	FormerNames       []string       `json:"former_names,omitempty"`
	Traded            bool           `json:"traded,omitempty"`
	DeletionDate      string         `json:"deletion_date,omitempty"`
	Sex               string         `json:"sex"`
	Title             string         `json:"title"`
	UnlockedTitles    int            `json:"unlocked_titles"`
	Vocation          string         `json:"vocation"`
	Level             int            `json:"level"`
	AchievementPoints int            `json:"achievement_points"`
	World             string         `json:"world"`
	FormerWorlds      []string       `json:"former_worlds,omitempty"`
	Residence         string         `json:"residence"`
	MarriedTo         string         `json:"married_to,omitempty"`
	Houses            []Houses       `json:"houses,omitempty"`
	Guild             CharacterGuild `json:"guild"`
	LastLogin         string         `json:"last_login,omitempty"`
	AccountStatus     string         `json:"account_status"`
	Comment           string         `json:"comment,omitempty"`
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

// Child of Deaths
type Killers struct {
	Name   string `json:"name"`
	Player bool   `json:"player"`
	Traded bool   `json:"traded"`
	Summon string `json:"summon"`
}

// Child of Characters
type Deaths struct {
	Time    string    `json:"time"`
	Level   int       `json:"level"`
	Killers []Killers `json:"killers"`
	Assists []Killers `json:"assists"`
	Reason  string    `json:"reason"`
}

// Child of Characters
type AccountInformation struct {
	Position     string `json:"position,omitempty"`
	Created      string `json:"created,omitempty"`
	LoyaltyTitle string `json:"loyalty_title,omitempty"`
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
	AccountBadges      []AccountBadges    `json:"account_badges,omitempty"`
	Achievements       []Achievements     `json:"achievements,omitempty"`
	Deaths             []Deaths           `json:"deaths,omitempty"`
	AccountInformation AccountInformation `json:"account_information,omitempty"`
	OtherCharacters    []OtherCharacters  `json:"other_characters,omitempty"`
}

//
// The base includes two levels, Characters and Information
type CharacterResponse struct {
	Characters  Characters  `json:"characters"`
	Information Information `json:"information"`
}

// From https://pkg.go.dev/golang.org/x/net/html/atom
// This is an Atom. An Atom is an integer code for a string.
// Instead of importing the whole lib, we thought it would be
// best to just simply use the Br constant value.
const Br = 0x202

var (
	deathRegex               = regexp.MustCompile(`<td.*>(.*)<\/td><td>(.*) at Level ([0-9]+) by (.*).<\/td>`)
	summonRegex              = regexp.MustCompile(`(an? .+) of ([^<]+)`)
	accountBadgesRegex       = regexp.MustCompile(`\(this\), &#39;(.*)&#39;, &#39;(.*)&#39;,.*\).*src="(.*)" alt=.*`)
	accountAchievementsRegex = regexp.MustCompile(`<td class="[a-zA-Z0-9_.-]+">(.*)<\/td><td>(.*)?<?.*<\/td>`)
	titleRegex               = regexp.MustCompile(`(.*) \(([0-9]+).*`)
	characterInfoRegex       = regexp.MustCompile(`<td.*<nobr>[0-9]+\..(.*)<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td><td style="width: 70%">(.*)<\/td><td.*`)
)

// TibiaCharactersCharacterV3 func
func TibiaCharactersCharacterV3Impl(BoxContentHTML string) (*CharacterResponse, error) {
	var (
		// local strings used in this function
		localDivQueryString = ".TableContentContainer tr"
		localTradedString   = " (traded)"

		// Declaring vars for later use..
		CharacterInformationData Character
		AccountBadgesData        []AccountBadges
		AchievementsData         []Achievements
		DeathsData               []Deaths
		AccountInformationData   AccountInformation
		OtherCharactersData      []OtherCharacters

		// Errors
		characterNotFound bool
		insideError       error
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("TibiaCharactersCharacterV3Impl failed at  goquery.NewDocumentFromReader, err: %s", err)
	}

	// Running query on every .TableContainer
	ReaderHTML.Find(".TableContainer").EachWithBreak(func(index int, s *goquery.Selection) bool {
		if insideError != nil {
			return false
		}

		SectionTextQuery := s.Find("div.Text")

		SectionName := SectionTextQuery.Nodes[0].FirstChild.Data

		// query current node with goquery
		CharacterDivQuery := goquery.NewDocumentFromNode(s.Nodes[0])

		switch SectionName {
		case "Could not find character":
			characterNotFound = true
			return false
		case "Character Information", "Account Information":
			// Running query over each tr in character content container
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				RowNameQuery := s.Find("td[class^='Label']")

				RowName := RowNameQuery.Nodes[0].FirstChild.Data
				RowData := RowNameQuery.Nodes[0].NextSibling.FirstChild.Data

				switch TibiaDataSanitizeStrings(RowName) {
				case "Name:":
					Tmp := strings.Split(RowData, "<")
					CharacterInformationData.Name = strings.TrimSpace(Tmp[0])
					if strings.Contains(Tmp[0], ", will be deleted at") {
						Tmp2 := strings.Split(Tmp[0], ", will be deleted at ")
						CharacterInformationData.Name = Tmp2[0]
						CharacterInformationData.DeletionDate = TibiaDataDatetimeV3(strings.TrimSpace(Tmp2[1]))
					}
					if strings.Contains(RowData, localTradedString) {
						CharacterInformationData.Traded = true
						CharacterInformationData.Name = strings.Replace(CharacterInformationData.Name, localTradedString, "", -1)
					}
				case "Former Names:":
					CharacterInformationData.FormerNames = strings.Split(RowData, ", ")
				case "Sex:":
					CharacterInformationData.Sex = RowData
				case "Title:":
					subma1t := titleRegex.FindAllStringSubmatch(RowData, -1)
					CharacterInformationData.Title = subma1t[0][1]
					CharacterInformationData.UnlockedTitles = TibiaDataStringToIntegerV3(subma1t[0][2])
				case "Vocation:":
					CharacterInformationData.Vocation = RowData
				case "Level:":
					CharacterInformationData.Level = TibiaDataStringToIntegerV3(RowData)
				case "nobr", "Achievement Points:":
					CharacterInformationData.AchievementPoints = TibiaDataStringToIntegerV3(RowData)
				case "World:":
					CharacterInformationData.World = RowData
				case "Former World:":
					CharacterInformationData.FormerWorlds = strings.Split(RowData, ", ")
				case "Residence:":
					CharacterInformationData.Residence = RowData
				case "Account Status:":
					CharacterInformationData.AccountStatus = RowData
				case "Married To:":
					AnchorQuery := s.Find("a")
					CharacterInformationData.MarriedTo = AnchorQuery.Nodes[0].FirstChild.Data
				case "House:":
					AnchorQuery := s.Find("a")
					HouseName := AnchorQuery.Nodes[0].FirstChild.Data
					HouseHref := AnchorQuery.Nodes[0].Attr[0].Val
					//substring from houseid= to &character in the href for the house
					HouseId := HouseHref[strings.Index(HouseHref, "houseid")+8 : strings.Index(HouseHref, "&character")]
					HouseRawData := RowNameQuery.Nodes[0].NextSibling.LastChild.Data
					HouseTown := HouseRawData[strings.Index(HouseRawData, "(")+1 : strings.Index(HouseRawData, ")")]
					HousePaidUntil := HouseRawData[strings.Index(HouseRawData, "is paid until ")+14:]

					CharacterInformationData.Houses = append(CharacterInformationData.Houses, Houses{
						Name:    HouseName,
						Town:    HouseTown,
						Paid:    TibiaDataDateV3(HousePaidUntil),
						HouseID: TibiaDataStringToIntegerV3(HouseId),
					})
				case "Guild Membership:":
					CharacterInformationData.Guild.Rank = strings.TrimSuffix(RowData, " of the ")

					//TODO: I don't understand why the unicode nbsp is there...
					CharacterInformationData.Guild.GuildName = TibiaDataSanitizeStrings(RowNameQuery.Nodes[0].NextSibling.LastChild.LastChild.Data)
				case "Last Login:":
					if RowData != "never logged in" {
						CharacterInformationData.LastLogin = TibiaDataDatetimeV3(RowData)
					}
				case "Comment:":
					node := RowNameQuery.Nodes[0].NextSibling.FirstChild

					stringBuilder := strings.Builder{}
					for node != nil {
						if node.DataAtom == Br {
							//It appears we can ignore br because either the encoding or goquery adds an \n for us
							//stringBuilder.WriteString("\n")
						} else {
							stringBuilder.WriteString(node.Data)
						}

						node = node.NextSibling
					}

					CharacterInformationData.Comment = stringBuilder.String()
				case "Loyalty Title:":
					if RowData != "(no title)" {
						AccountInformationData.LoyaltyTitle = RowData
					}
				case "Created:":
					AccountInformationData.Created = TibiaDataDatetimeV3(RowData)
				case "Position:":
					TmpPosition := strings.Split(RowData, "<")
					AccountInformationData.Position = strings.TrimSpace(TmpPosition[0])
				default:
					log.Println("LEFT OVER: `" + RowName + "` = `" + RowData + "`")
				}
			})
		case "Account Badges":
			// Running query over each tr in list
			CharacterDivQuery.Find(".TableContentContainer tr td").EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterV3Impl failed at s.Html() inside Account Badges, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaksV3(CharacterListHTML)

				// prevent failure of regex that parses account badges
				if CharacterListHTML != "There are no account badges set to be displayed for this character." {
					subma1 := accountBadgesRegex.FindAllStringSubmatch(CharacterListHTML, -1)

					AccountBadgesData = append(AccountBadgesData, AccountBadges{
						Name:        subma1[0][1],
						IconURL:     subma1[0][3],
						Description: subma1[0][2],
					})
				}

				return true
			})
		case "Account Achievements":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterV3Impl failed at s.Html() inside Account Achievements, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaksV3(CharacterListHTML)

				subma1a := accountAchievementsRegex.FindAllStringSubmatch(CharacterListHTML, -1)
				if len(subma1a) > 0 {
					// fixing encoding for achievement name
					subma1a[0][2] = TibiaDataSanitizeEscapedString(subma1a[0][2])

					// get the name of the achievement (and ignore the secret image on the right)
					Name := strings.Split(subma1a[0][2], "<img")

					AchievementsData = append(AchievementsData, Achievements{
						Name:   Name[0],
						Grade:  strings.Count(subma1a[0][1], "achievement-grade-symbol"),
						Secret: strings.Contains(subma1a[0][2], "achievement-secret-symbol"),
					})
				}

				return true
			})
		case "Character Deaths":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterV3Impl failed at s.Html() inside Character Deaths, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaksV3(CharacterListHTML)
				CharacterListHTML = strings.ReplaceAll(CharacterListHTML, ".<br/>Assisted by", ". Assisted by")

				// Regex to get data for deaths
				subma1 := deathRegex.FindAllStringSubmatch(CharacterListHTML, -1)

				if len(subma1) > 0 {
					// defining responses
					DeathKillers := []Killers{}
					DeathAssists := []Killers{}

					// store for reply later on.. and sanitizing string
					ReasonString := TibiaDataSanitizeStrings(RemoveHtmlTag(subma1[0][2] + " at Level " + subma1[0][3] + " by " + subma1[0][4] + "."))

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
					DeathsData = append(DeathsData, Deaths{
						Time:    TibiaDataDatetimeV3(subma1[0][1]),
						Level:   TibiaDataStringToIntegerV3(subma1[0][3]),
						Killers: DeathKillers,
						Assists: DeathAssists,
						Reason:  ReasonString,
					})
				}

				return true
			})
		case "Characters":
			// Running query over each tr in character list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterV3Impl failed at s.Html() inside Characters, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaksV3(CharacterListHTML)

				subma1 := characterInfoRegex.FindAllStringSubmatch(CharacterListHTML, -1)

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
						Tmp := strings.Split(TmpCharacterName, "<")
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
						Name:    TibiaDataSanitizeStrings(TmpCharacterName),
						World:   subma1[0][2],
						Status:  TmpStatus,
						Deleted: TmpDeleted,
						Main:    TmpMain,
						Traded:  TmpTraded,
					})
				}

				return true
			})
		}

		return true
	})

	// Build the character data
	charData := Characters{
		CharacterInformationData,
		AccountBadgesData,
		AchievementsData,
		DeathsData,
		AccountInformationData,
		OtherCharactersData,
	}

	// Search for errors
	switch {
	case characterNotFound:
		return nil, validation.ErrorCharacterNotFound
	case insideError != nil:
		return nil, insideError
	case reflect.DeepEqual(charData, Characters{}):
		// There are some rare cases where a character name would
		// bug out tibia.com (tíbia, for example) and then we would't
		// receive the character not found error, for these edge cases
		// we check if the char structure is empty, if it is, it means
		// the character has not been found
		//
		// Validating those names would also be a pain because of old
		// tibian names such as Kolskägg, which for whatever reason is valid
		return nil, validation.ErrorCharacterNotFound
	}

	//
	// Build the data-blob
	return &CharacterResponse{
		charData,
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}

// TibiaDataParseKiller func - insert a html string and get the killers back
func TibiaDataParseKiller(data string) (string, bool, bool, string) {
	var (
		// local strings used in this function
		localTradedString = " (traded)"

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
	if strings.HasPrefix(data, "a ") || strings.HasPrefix(data, "an ") {
		if containsCreaturesWithOf(data) {
			// this is not a summon, since it is a creature with a of in the middle
		} else {
			rs := summonRegex.FindAllStringSubmatch(data, -1)
			if len(rs) >= 1 {
				theSummon = rs[0][1]
				data = rs[0][2]
			}
		}
	}

	// sanitizing string
	data = TibiaDataSanitizeStrings(data)

	return data, isPlayer, isTraded, theSummon
}

// containsCreaturesWithOf checks if creature is present in special creatures list
func containsCreaturesWithOf(str string) bool {
	// this list should be based on the https://assets.tibiadata.com/data.json creatures name and plural_name field (currently only singular version)
	creaturesWithOf := []string{
		"acolyte of the cult",
		"adept of the cult",
		"cloak of terror",
		"energuardian of tales",
		"enlightened of the cult",
		"guardian of tales",
		"hand of cursed fate",
		"monk of the order",
		"novice of the cult",
		"priestess of the wild sun",
		"sight of surrender",
		"son of verminor",
	}

	// trim away "an " and "a "
	str = strings.TrimPrefix(strings.TrimPrefix(str, "an "), "a ")

	for _, v := range creaturesWithOf {
		if v == str {
			return true
		}
	}
	return false
}
