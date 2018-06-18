package user

import (
	"encoding/json"
	"log"

	"github.com/as/slack"
)

// Info returns details for the user identified by the given user id. If locale
// is true, additional user local information is returned in the user response.
func Info(c *slack.Client, uid string, locale bool) (*User, error) {
	var R struct {
		Ok   bool  `json:"ok"`
		User *User `json:"user"`
	}

	r, err := c.Post(c.BaseURL()+"/api/users.info", slack.Encode(struct {
		Token  string `url:"token"`
		User   string `url:"user"`
		Locale bool   `url:"include_locale"`
	}{
		Token:  c.Token(),
		User:   uid,
		Locale: locale,
	}))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&R); err != nil {
		return nil, err
	}
	if R.Ok {
		log.Println("no ok!")
	}
	log.Printf("history: %#v", R)
	return R.User, nil
}
