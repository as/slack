package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/as/slack"
)

var (
	del = flag.Bool("d", false,  "delete")
	email = flag.String("e", "", "email")
	pass  = flag.String("p", "", "password")
	space = flag.String("s", "gophers", "slack workspace")
	ch    = flag.String("c", "GB1KBRGKA", "slack channel")
	token      = flag.String("x", os.Getenv("slacktoken"), "api xox. token (email and password not required if set)")
)

func main() {
	flag.Parse()
	if *space == "" || *ch == "" {
		log.Fatal("workspace and channel required (set -s and -c)")
	}
	if *token == "" && (*email == "" || *pass == "") {
		log.Fatal("email and password (or token) required")
	}

	c := slack.NewClient(*space, nil)
	if *token == "" || *token == "dump"{
		err := slack.Login(c, *email, *pass)
		ck("login", err)
		if *token == "dump"{
			log.Fatalln(c.Token())
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		m, err := slack.SendMsg(c, *ch, sc.Text())
		ck("send", err)
		if *del{
			if err = slack.DelMsg(c, *ch, m.TS); err != nil{
				log.Println("slack: delete:", err)
			}
		}
	}
	ck("done", sc.Err())
}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
