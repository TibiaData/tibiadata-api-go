package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strings"
)

type ForumThreadPostAuthorGuild struct {
	Name  string `json:"name,omitempty"`
	Rank  string `json:"rank,omitempty"`
	Title string `json:"title,omitempty"`
}

type ForumThreadPostAuthor struct {
	CharacterName string                      `json:"character_name"`
	IsTraded      bool                        `json:"is_traded"`
	Server        string                      `json:"server,omitempty"`
	Vocation      string                      `json:"vocation,omitempty"`
	Level         int                         `json:"level,omitempty"`
	Posts         int                         `json:"posts,omitempty"`
	Guild         *ForumThreadPostAuthorGuild `json:"guild,omitempty"`
}

type ForumThreadPost struct {
	ID      int                   `json:"id"`
	Subject string                `json:"subject,omitempty"`
	Icon    string                `json:"icon,omitempty"`
	Message string                `json:"message"`
	Author  ForumThreadPostAuthor `json:"author"`
}

type ForumThreadPagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_records"`
}

type ForumThread struct {
	ID         int                   `json:"id"`
	Title      string                `json:"title"`
	Posts      []ForumThreadPost     `json:"posts"`
	Pagination ForumThreadPagination `json:"pagination"`
}

type ForumThreadResponse struct {
	Thread      ForumThread `json:"thread"`
	Information Information `json:"information"`
}

var (
	threadTitleRegex             = regexp.MustCompile(`<div class="ForumTitleText"><b>(.*)<\/b><\/div><\/div>`)
	threadResultsRegex           = regexp.MustCompile(`<td class="PageNavigation">.*Results: ([0-9,]+)<\/b><\/div><\/small><\/td>`)
	totalPagesRegex              = regexp.MustCompile(`<span class="PageLink.*&pagenumber=([0-9]+).*<\/span>`)
	postIdRegex                  = regexp.MustCompile(`<a.*postid=([0-9]+)">Quote<\/a>`)
	postSubjectRegex             = regexp.MustCompile(`<div class="PostText">.*<b>(.*)<\/b>`)
	postIconRegex                = regexp.MustCompile(`<div class="PostText"><img src=".*" border="0" width="15" height="15" alt="([a-zA-Z ]+)"\/>`)
	postTextRegex                = regexp.MustCompile(`<div class="PostText">.*<div class="ForumDesktopElements"><br/></div><br/>(?s)(.*)</div></div><div class="PostLower">`)
	postCharacterNameRegex       = regexp.MustCompile(`<div><b><a href="https://www.tibia.com/community/\?subtopic=characters&amp;name=[a-zA-Z+]+">([a-zA-Z ]+)</a>`)
	postCharacterNameTradedRegex = regexp.MustCompile(`<div><b>([a-zA-Z ]+) \(traded\)</b>`)
	postCharacterServerRegex     = regexp.MustCompile(`Inhabitant of ([a-zA-Z]+)<br/>`)
	postCharacterLevelRegex      = regexp.MustCompile(`Level: ([0-9,]+)<br/>`)
	postCharacterPostsRegex      = regexp.MustCompile(`Posts: ([0-9,]+)<br/>`)
	postCharacterVocationRegex   = regexp.MustCompile(`Vocation: ([a-zA-Z ]+)<br/>`)
	postCharacterGuildRegex      = regexp.MustCompile(`<br/><br/><font class="ff_smallinfo">(.*) of the <a href="https://www.tibia.com/community/\?subtopic=guilds&amp;page=view&amp;GuildName=.*">(.*)</a>(?:( \(.*\)))</font>`)
	spaceRegex                   = regexp.MustCompile(`\s+`)
)

func TibiaForumThread3Impl(threadId string, BoxContentHTML string, currentPage int) ForumThreadResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	var (
		PostsData                 []ForumThreadPost
		threadTitle               string
		totalPages, threadResults int
	)

	subma1 := threadTitleRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		threadTitle = subma1[0][1]
	}

	subma1 = threadResultsRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		threadResults = TibiaDataStringToIntegerV3(subma1[0][1])
	}

	subma1 = totalPagesRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		totalPages = TibiaDataStringToIntegerV3(subma1[0][1])
		if currentPage > totalPages {
			totalPages = currentPage
		}
	} else {
		totalPages = 1
	}

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent td.CipPost").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CreatureDivHTML
		PostDivHtml, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		PostAuthorHtml, err := s.Find(".PostCharacterText").Html()
		if err != nil {
			log.Fatal(err)
		}

		subma2 := postIdRegex.FindAllStringSubmatch(PostDivHtml, -1)
		if len(subma2) == 0 {
			return
		}

		var (
			postSubject, postIcon, postText, authorName, authorServer, authorVocation, authorGuildTitle string
			postIsTraded                                                                                bool
			authorLevel, authorPosts                                                                    int
			authorGuild                                                                                 *ForumThreadPostAuthorGuild
		)

		subma3 := postSubjectRegex.FindAllStringSubmatch(PostDivHtml, -1)
		if len(subma3) > 0 {
			postSubject = subma3[0][1]
		}

		subma3 = postIconRegex.FindAllStringSubmatch(PostDivHtml, -1)
		if len(subma3) > 0 {
			postIcon = subma3[0][1]
		}

		subma3 = postTextRegex.FindAllStringSubmatch(PostDivHtml, -1)
		if len(subma3) > 0 {
			postText = strings.Split(subma3[0][1], "________________")[0]
			postText = spaceRegex.ReplaceAllString(postText, " ")
		}

		subma3 = postCharacterNameRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			authorName = subma3[0][1]
		} else {
			subma3 = postCharacterNameTradedRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
			authorName = subma3[0][1]
			postIsTraded = true
		}

		subma3 = postCharacterServerRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			authorServer = subma3[0][1]
		}

		subma3 = postCharacterLevelRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			authorLevel = TibiaDataStringToIntegerV3(subma3[0][1])
		}

		subma3 = postCharacterVocationRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			authorVocation = subma3[0][1]
		}

		subma3 = postCharacterPostsRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			authorPosts = TibiaDataStringToIntegerV3(subma3[0][1])
		}

		subma3 = postCharacterGuildRegex.FindAllStringSubmatch(PostAuthorHtml, -1)
		if len(subma3) > 0 {
			if len(subma3[0]) > 2 {
				authorGuildTitle = strings.Trim(subma3[0][3], "() ")
			}

			authorGuild = &ForumThreadPostAuthorGuild{
				Name:  subma3[0][2],
				Rank:  subma3[0][1],
				Title: authorGuildTitle,
			}
		}

		PostsData = append(PostsData, ForumThreadPost{
			ID:      TibiaDataStringToIntegerV3(subma2[0][1]),
			Subject: postSubject,
			Icon:    postIcon,
			Message: postText,
			Author: ForumThreadPostAuthor{
				CharacterName: authorName,
				IsTraded:      postIsTraded,
				Server:        authorServer,
				Level:         authorLevel,
				Vocation:      authorVocation,
				Posts:         authorPosts,
				Guild:         authorGuild,
			},
		})
	})

	return ForumThreadResponse{
		Thread: ForumThread{
			ID:    TibiaDataStringToIntegerV3(threadId),
			Title: threadTitle,
			Posts: PostsData,
			Pagination: ForumThreadPagination{
				CurrentPage:  currentPage,
				TotalPages:   totalPages,
				TotalResults: threadResults,
			},
		},
		Information: Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
