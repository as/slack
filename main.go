package main

import (
	"flag"
	"fmt"
	"log"
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

	c := NewClient(*space, nil)

	err := Login(c, *email, *pass)
	ck("login", err)
	fmt.Println("api_token", c.token)

	// err = SendMsg(c, "#foo", "message")
	// ck("send", err)

}

func ck(where string, err error) {
	if err != nil {
		log.Fatalf("%s: ", err)
	}
}
