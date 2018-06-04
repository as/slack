package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrBadKey = errors.New("bad format for xoxs key")
)

type ErrMissing struct {
	Key string
}

func (e ErrMissing) Error() string {
	return fmt.Sprintf("missing value for key %q", e.Key)
}

func Login(c *Client, email, pass string) (err error) {
	m := url.Values{
		"crumb":    {},
		"email":    {email},
		"password": {pass},
		"redir":    {},
		"signin":   {"1"},
	}
	link := c.BaseURL()
	log.Println("baseurl is", link)

	r, err := c.Get(link)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := parseForm(r.Body, m, "crumb"); err != nil {
		return err
	}

	r, err = c.Do("POST", link, m)
	if err != nil {
		return err
	}
	r.Body.Close()

	r, err = c.Get(link + "/messages")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	sc := bufio.NewScanner(r.Body)
	for sc.Scan() {
		if strings.Contains(sc.Text(), "api_token") {
			break
		}
	}
	if sc.Err() != nil {
		return sc.Err()
	}

	sp := strings.Index(sc.Text(), `api_token: "`)
	ep := strings.LastIndex(sc.Text(), `"`)
	if sp == -1 || ep == -1 {
		return ErrBadKey
	}
	sp += len(`api_token: "`)
	c.token = sc.Text()[sp:ep]
	return nil
}

func parseForm(r io.Reader, m url.Values, require ...string) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		k := ""
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, x := range n.Attr {
				if x.Key == "name" {
					k = x.Val
				} else if x.Key == "value" {
					if _, ok := m[k]; ok {
						if len(m[k]) == 0 {
							m[k] = []string{x.Val}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, x := range require {
		if len(m[x]) == 0 {
			return ErrMissing{x}
		}
	}
	return nil
}
