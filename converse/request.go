package converse

import (
	"fmt"

	"github.com/as/slack"
)

type Kind int

const (
	KindPublic = 1 << iota
	KindPrivate
	KindMPIM
	KindIM
)

var kindTab = map[Kind]string{
	KindPublic:  "public_channel",
	KindPrivate: "private_channel",
	KindMPIM:    "mpim",
	KindIM:      "im",
}

func List(c *slack.Client, include Kind, list *[]Channel) (next string, err error) {
	var R struct {
		Ok      bool      `json:"ok"`
		Channel []Channel `json:"channels"`
		Meta    struct {
			Next string `json:"next_cursor"`
		} `json:"response_metadata"`
	}

	lim := 100
	if cap(*list) != 0 && cap(*list) < 1000 {
		lim = cap(*list)
	}

	r, err := c.Post(c.BaseURL()+"/api/conversations.list", slack.Encode(struct {
		Token string `url:"token"`
		Limit string `url:"limit"`
		Types string `url:"types"`
	}{
		Token: c.Token(),
		Limit: fmt.Sprint(lim),
		Types: "public_channel,private_channel",
	}))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if err = slack.Decode(r, &R); err != nil {
		return "", err
	}
	*list = append(*list, R.Channel...)
	return R.Meta.Next, err
}
