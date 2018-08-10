package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/as/slack"
	"github.com/as/slack/converse"
)

var (
	space  = flag.String("s", os.Getenv("slackworkspace"), `slack workspace name i.e., the "foo" in foo.slack.com`)
	email  = flag.String("e", "", "email")
	pass   = flag.String("p", "", "password")
	raw    = flag.String("f", "", "raw file to save json results (contains more details)")
	replay = flag.Bool("r", false, "replay the old contents in the raw file as if the messages were obtained online")
	sleep  = flag.Duration("sleep", time.Second, "duration to wait before downloading each next page")
	token  = flag.String("x", os.Getenv("slacktoken"), "api xox. token (email and password not required if set)")
)

func main() {
	flag.Parse()
	if *space == "" {
		log.Fatal("missing flag: -s workspace")
	}
	if *token == "" && (*email == "" || *pass == "") {
		log.Fatal("email and password (or token) required")
	}

	c := slack.NewClient(*space, &slack.Config{
		Token: *token,
	})
	if *token == "" {
		err := slack.Login(c, *email, *pass)
		ck("login", err)
	}
	var (
		err  error
		list = []converse.Channel{}
		next string
	)
	for err == nil {
		next, err = converse.List(c, 0, &list, next)
	}
	for _, v := range list {
		fmt.Println(v.Name, v.Id, v.LastRead, v.NumMembers, v.Topic, v.Purpose, v.Created, v.Creator)
	}
	if *raw != "" {
		fd, err := os.Create(*raw)
		ck("raw", err)
		defer fd.Close()
		enc := json.NewEncoder(fd)
		for _, v := range list {
			enc.Encode(v)
		}
	}
}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
