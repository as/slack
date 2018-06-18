package slack

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendMsg(c *Client, to string, text string) (err error) {
	r, err := c.Post(c.BaseURL()+"/api/chat.postMessage", &Message{
		Channel: to,
		Text:    text,
	})
	if err != nil {
		return err
	}
	defer r.Body.Close()
	io.Copy(os.Stdout, r.Body)
	return err
}

type Attachment struct {
	Text       string `json:"text"`
	Type       string `json:"attachment_type"`
	Fallback   string `json:"fallback"`
	Color      string `json:"color"`
	CallbackID string `json:"callback_id"`
	Actions    []struct {
		Name       string `json:"name"`
		Text       string `json:"text"`
		Type       string `json:"type"`
		DataSource string `json:"data_source"`
	} `json:"actions"`
}

type Message struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`

	// Below only used for incoming messages
	Username  string `json:"username,omitempty"`
	User      string `json:"user,omitempty"`
	Subtype   string `json:"subtype,omitempty"`
	Ts        Ts     `json:"ts,omitempty"`
	Type      string `json:"type,omitempty"`
	BotID     string `json:"bot_id,omitempty"`
	IsStarred bool   `json:"is_starred,omitempty"`
	Reactions []struct {
		Users []string `json:"users"`
		Name  string   `json:"name"`
		Count int      `json:"count"`
	} `json:"reactions,omitempty"`

	Attachments []struct {
		Text     string `json:"text"`
		Id       int    `json:"id"`
		Fallback string `json:"fallback"`
	} `json:"attachments,omitempty"`
}

type Ts time.Time

func (t Ts) String() string {
	return time.Time(t).Format("20060102.150405")
}

func (t Ts) MarshalText() (data []byte, err error) {
	tm := time.Time(t)
	return []byte(fmt.Sprintf("%d.%06d", tm.Unix(), tm.UnixNano()%int64(time.Millisecond))), nil
}
func (t *Ts) UnmarshalText(data []byte) (err error) {
	x := strings.Split(string(data), ".")
	if len(x) < 2 {
		return fmt.Errorf("bad timestamp: %s", data)
	}
	s, err := strconv.Atoi(x[0])
	if err != nil {
		return err
	}
	ms, err := strconv.Atoi(x[1])
	if err != nil {
		return err
	}

	*t = Ts(time.Unix(int64(s), int64(ms)))
	return nil
}

func (m Message) String() string {
	if m.Username != "" {
		return fmt.Sprintf("%s\t%s: %q", m.Ts, m.Username, m.Text)
	}
	return fmt.Sprintf("%s\t%s: %q", m.Ts, m.User, m.Text)
}
