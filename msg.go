package main

import (
	"io"
	"os"
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
	Channel string `json:"channel"`
	Text    string `json:"text,omitempty"`
}
