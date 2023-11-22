package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBoardsBoard(t *testing.T) {
	data, err := os.ReadFile("../testdata/forum/thread.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	threadJson := TibiaForumThread3Impl("4729442", string(data), 1)
	assert := assert.New(t)

	assert.Equal(4729442, threadJson.Thread.ID)
	assert.Equal("ON VACATION DONT SEND 9-13Th -- Farideo's Imbuement Service", threadJson.Thread.Title)

	assert.Equal(1, threadJson.Thread.Pagination.CurrentPage)
	assert.Equal(656, threadJson.Thread.Pagination.TotalResults)
	assert.Equal(33, threadJson.Thread.Pagination.TotalPages)

	assert.Equal(20, len(threadJson.Thread.Posts))
	assert.Equal(38715156, threadJson.Thread.Posts[0].ID)
	assert.Equal("Farideo", threadJson.Thread.Posts[0].Author.CharacterName)
	assert.Equal(false, threadJson.Thread.Posts[0].Author.IsTraded)
	assert.Equal("Antica", threadJson.Thread.Posts[0].Author.Server)
	assert.Equal(524, threadJson.Thread.Posts[0].Author.Level)
	assert.Equal(193, threadJson.Thread.Posts[0].Author.Posts)
	assert.Equal("Its bad", threadJson.Thread.Posts[0].Author.Guild.Name)
	assert.Equal("Dreadful", threadJson.Thread.Posts[0].Author.Guild.Rank)
	assert.Equal("Antica Imbue Guy", threadJson.Thread.Posts[0].Author.Guild.Title)
	assert.Equal("Royal Paladin", threadJson.Thread.Posts[0].Author.Vocation)
	assert.Equal("Edit: July 2022<br/> Over 3000+ Imbuements complete =)<br/> <br/> <b>As this was a service to help out fellow Antican&#39;s for getting their imbuements and not get scammed when no services were available, I believe the server has grown and gained several trusted services at the moment and whilst this is a free service I&#39;ll be keeping an eye on the forums to see if my services are needed or not but I will not be checking my inbox as regularly but if you do need my services you can BUMP this thread to grab my attention and I&#39;ll be happy to serve</b><br/> <br/> Delivery is from 5 minutes to 24 hours depending on the time you sent your package and the time I login. I usually login between 10pm to 2am. (48 hours MAX if i forget to login 1 day)<br/> <br/> For those of you who are new to imbues here&#39;s a link that can give you some knowledge.<br/> <a href=\"?action=externallinkwarning&amp;target=https%3A%2F%2Ftibia.fandom.com%2Fwiki%2FImbuing\" target=\"_blank\" rel=\"noopener noreferrer\">https://tibia.fandom.com/wiki/Imbuing</a><br/> <br/> <b> I Only accept PARCELS if you catch me online I am most probably working on other imbues and will likely logout instantly the moment I am done so I don&#39;t think I&#39;ll be able to help you out if you haven&#39;t already sent before I logged in.</b><i> Id like to ask you to <b>Parcel</b> me the items you need imbued along with the materials needed. All tier 3 imbuements cost 150k per imbue (This is money I don&#39;t take its used in the process of imbuement). <b>Do not forget to include a FRESH parcel and a label with your name on it inside that parcel</b> to make it easier for me to send you back your items. Finally please include a letter with the description of which items need which imbuements.<br/> For example : Crit weapon / Mana Helmet / Life Armor. <br/> </i><br/> <br/> PS:<br/> <br/> Most of you know that I am <b>retired</b> but I keep my account premium renewed only to help you out. I literally log in every night just to check my inbox for imbuement parcels. <br/> <br/> For those of you who treat me like a slave, I am not. <br/> <br/> Please remember this is a free service and any appreciation whatsoever helps me continue this service. Remember I also take time to help you out so the least you can do is throw me a thank you on a letter and post on forums. I sometimes receive items and imbuement materials that are worth 30kk+ and I do not get any appreciation at all. Some of you tip me and that is also truly appreciated but is not a must although of&#39;course it helps me out with paying rent etc.<br/> <br/> Regards,<br/> <br/> Farideo<br/> <br/> <br/> <b>PS: I would like to thank all of you who have been supporting me while I attend to this service.</b><br/>", threadJson.Thread.Posts[0].Message)

	postWithSubject := threadJson.Thread.Posts[15]
	assert.Equal("Trustworthy guy !!!", postWithSubject.Subject)
	assert.Equal("Red face", postWithSubject.Icon)

	postWithIcon := threadJson.Thread.Posts[18]
	assert.Equal("Cool", postWithIcon.Icon)

	tradedCharacter := threadJson.Thread.Posts[2]
	assert.Equal("Dzik Schematix", tradedCharacter.Author.CharacterName)
	assert.Equal(true, tradedCharacter.Author.IsTraded)
}
