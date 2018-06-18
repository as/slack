package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/as/slack"
)

var (
	email = flag.String("e", "", "email")
	pass  = flag.String("p", "", "password")
	space = flag.String("s", "", "slack workspace")
	ch    = flag.String("c", "C029RQSEE", "slack channel")

	token = flag.String("x", "api xox. token (email and password not required if set)", "")
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
	err := slack.Login(c, *email, *pass)
	ck("login", err)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		err = slack.SendMsg(c, *ch, sc.Text())
		ck("send", err)
	}
	ck("done", sc.Err())
}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
