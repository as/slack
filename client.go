package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

var DefaultConfig = Config{
	Token: "",
	Jar:   jar,
}

type Config struct {
	Token string
	Jar   *cookiejar.Jar
}

type Client struct {
	*http.Client
	token     string
	workspace string
}

// NewClient returns a slack Client for the given workspace. If config
// is nil, the DefaultConfig is used, which includes a cookiejar with
// common prefixes.
//
// The token and cookie jar can be set in the config to avoid calling
// the Login package function.
func NewClient(workspace string, conf *Config) *Client {
	if conf == nil {
		conf = &DefaultConfig
	}
	return &Client{
		workspace: workspace,
		token:     conf.Token,
		Client: &http.Client{
			Jar: conf.Jar,
		},
	}
}

func (c *Client) BaseURL() string {
	return fmt.Sprintf("https://%s.slack.com", c.workspace)
}

func (c *Client) Token() string {
	return c.token
}

// Do is like http.Client.Do, except it sets the bearer token (if present) and assumes the
// content type is json (marshalling the body as-needed if the body is non-nil). If the
// body is a url.Values, Do calls the underlying clients' PostForm method (currently this
// does not add any authorization headers)
func (c *Client) Do(method string, path string, body interface{}) (r *http.Response, err error) {
	if m, ok := body.(url.Values); ok {
		return c.Client.PostForm(path, m)
	}

	data, err := c.marshal(body)
	if err != nil {
		return c.fail(err)
	}

	tx, err := http.NewRequest(method, path, data)
	if err != nil {
		return c.fail(err)
	}

	if c.token != "" {
		tx.Header.Set("Authorization", "Bearer "+c.token)
	}
	if data != nil {
		tx.Header.Set("Content-Type", "application/json")
	}

	if r, err = c.Client.Do(tx); err != nil {
		return c.fail(err)
	}

	return r, nil
}

// Get issues a post request, see Do for details
func (c *Client) Get(path string) (r *http.Response, err error) {
	return c.Do("GET", path, nil)
}

// Post issues a post request, see Do for details
func (c *Client) Post(path string, body interface{}) (r *http.Response, err error) {
	return c.Do("POST", path, body)
}

func (c *Client) fail(err error) (*http.Response, error) {
	return nil, err
}

func (c *Client) marshal(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

var jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
