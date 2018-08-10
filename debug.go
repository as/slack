package slack

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var debug = os.Getenv("slackdebug") == "on"

func Dump(req interface{}) {
	switch req := req.(type) {
	case *http.Request:
		x, _ := httputil.DumpRequestOut(req, true)
		log.Printf("Request \n%s\n", x)
	case *http.Response:
		x, _ := httputil.DumpResponse(req, true)
		log.Printf("Response \n%s\n", x)
	}
}
