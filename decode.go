package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

var ErrNotOk = errors.New("not ok")

func Decode(r *http.Response, v interface{}) (err error) {
	var b bytes.Buffer
	b.ReadFrom(r.Body)
	rd := bytes.NewReader(b.Bytes())

	var ok struct {
		Ok bool `json:"ok"`
	}
	if err = json.NewDecoder(rd).Decode(&ok); err != nil || !ok.Ok {
		rd.Seek(0, io.SeekStart)
		io.Copy(os.Stderr, rd)
		if err != nil {
			return err
		}
		return ErrNotOk
	}
	rd.Seek(0, io.SeekStart)
	return json.NewDecoder(rd).Decode(v)
}
