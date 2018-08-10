package converse_test

import (
	"io"

	"github.com/as/slack"
	"github.com/as/slack/converse"

	"flag"
	"log"
	"os"
	"testing"
)

var (
	integ = flag.Bool("integration", false, "run integration test (requires -x and -w) or ($slacktestxox and $slacktestworkspace) set")
	xox   = flag.String("x", os.Getenv("slacktestxox"), "slack test token")
	ws    = flag.String("w", os.Getenv("slacktestworkspace"), "slack test workspace")

	client *slack.Client
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *integ && *xox == "" {
		*integ = false
		log.Println("skipping integration tests: no token set")
	} else {
		client = slack.NewClient(*ws, &slack.Config{Token: *xox})
	}
	os.Exit(m.Run())
}

func TestIntList(t *testing.T) {
	if !*integ {
		t.Skip("integration test")
	}
	var (
		err  error
		next string
		list = []converse.Channel{}
	)
	for err == nil {
		next, err = converse.List(client, 0, &list, next)
	}
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	for i, v := range list {
		t.Logf("%d: %#v\n", i, v)
	}

}
