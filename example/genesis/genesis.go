package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/as/slack"
	 "github.com/as/slack/user"
)

var (
	space = flag.String("s", "", `slack workspace name i.e., the "foo" in foo.slack.com`)
	email = flag.String("e", "", "email")
	pass  = flag.String("p", "", "password")
	chid  = flag.String("c", "C029RQSEE", "channel id (must manually resolve for now)")

	token = flag.String("x", "api xox. token (email and password not required if set)", "")
)

func main() {
	flag.Parse()
	if *space == "" {
		log.Fatal("missing flag: -s workspace")
	}
	if *token == "" && (*email == "" || *pass == "") {
		log.Fatal("email and password (or token) required")
	}

	c := slack.NewClient(*space, nil)

	err := slack.Login(c, *email, *pass)
	ck("login", err)

	ts := slack.Ts(time.Now().Add(-time.Hour * 24 * 365 * 25))

	log.Printf("from %s\n", ts)
	m, err := slack.History(c, *chid,
		&slack.HistoryOpt{
			Count: 100,
			Start: ts,
		})
	ck("history", err)

	var umap = make(map[string]string)
	for _, m := range m {
		_, ok := umap[m.User]
		if ok {
			continue
		}
		u, err := user.Info(c, m.User, false)
		if err != nil {
			log.Printf("resolver: %q: %s", m.User, err)
			continue
		}
		umap[m.User] = u.Name
		log.Printf("resolver: %q -> %q", m.User, u.Name)
	}

	for i := range m {
		m[i].Username = umap[m[i].User]
		fmt.Println(m[i])
	}

	// err = SendMsg(c, "#foo", "message")
	// ck("send", err)

}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
