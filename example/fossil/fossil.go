package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/as/slack"
	"github.com/as/slack/user"
)

var (
	space = flag.String("s", "", `slack workspace name i.e., the "foo" in foo.slack.com`)
	email = flag.String("e", "", "email")
	pass  = flag.String("p", "", "password")
	chid  = flag.String("c", "C029RQSEE", "channel id (must manually resolve for now)")
	raw   = flag.String("f", "", "raw file to save json results (contains more details)")
	sleep = flag.Duration("sleep", time.Second, "duration to wait before downloading each next page")
	token = flag.String("x", "", "api xox. token (email and password not required if set)")
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
		umap = make(map[string]string)
		// ts = slack.Ts(time.Date(2016, time.August, 25, 0, 0, 0, 0, time.UTC))
		ts  = slack.Ts(time.Now().Add(-time.Hour * 24 * 365 * 25)) // 25 years ago should do it
		fd  *os.File
		enc *json.Encoder
		dec *json.Decoder
	)

	if *raw != "" {
		fd, err = os.OpenFile(*raw, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		ck("file", err)
		defer fd.Close()
		enc = json.NewEncoder(fd)
		dec = json.NewDecoder(fd)

		var m slack.Message
		i := 0
		for ; ; i++ {
			if err = dec.Decode(&m); err != nil {
				break
			}
			if name := umap[m.User]; name == "" && m.Username != "" {
				umap[m.User] = m.Username
				log.Printf("file: found mapping %q->%q", m.User, m.Username)
			}
			ts = m.Ts
		}
		log.Printf("read through %d messages (last ts=%s)", i, ts)
	}

	for {
		log.Printf("from %s\n", ts)
		m, err := slack.History(c, *chid, &slack.HistoryOpt{
			Count: 1000,
			Start: ts,
		})
		ck("history", err)

		if len(m) == 0 {
			log.Println("no more messages")
			break
		}

		ts = m[0].Ts
		for _, m := range m {
			name, ok := umap[m.User]
			if ok {
				log.Printf("resolver: already know %q->%q", m.User, name)
				continue
			}
			log.Printf("resolver: dont know %q->%q", m.User, name)
			time.Sleep(*sleep / 3)
			u, err := user.Info(c, m.User, false)
			if err != nil {
				log.Printf("resolver: %q: %s", m.User, err)
				continue
			}
			umap[m.User] = u.Name
			log.Printf("resolver: %q -> %q", m.User, u.Name)

		}

		for i := len(m); i != 0; i-- {
			u := &m[i-1]
			u.Username = umap[u.User]
			fmt.Println(u)
			if enc != nil {
				enc.Encode(u)
			}
		}

		time.Sleep(*sleep + time.Duration(rand.Intn(int(time.Second))))
	}
}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
