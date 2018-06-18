package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/as/slack"
	user "github.com/as/slack/users"
)

var (
	email = flag.String("e", "", "email")
	pass  = flag.String("p", "", "password")
	token = flag.String("x", "api xox. token (email and password not required if set)", "")
	space = flag.String("s", "", "")
)

func main() {
	flag.Parse()
	if *token == "" && (*email == "" || *pass == "") {
		log.Fatal("email and password (or token) required")
	}

	c := slack.NewClient(*space, nil)

	err := slack.Login(c, *email, *pass)
	ck("login", err)

	m, err := slack.History(c, "C029RQSEE", nil)
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
