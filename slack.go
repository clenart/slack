package slack

import (
	"errors"
	"log"
	"net/url"
	"os"
	"fmt"
)

var logger *log.Logger // A logger that can be set by consumers
/*
  Added as a var so that we can change this for testing purposes
*/
var SLACK_API string
var SLACK_WEBSOCKET_REFERER string
var SLACK_WEB_API_FORMAT string = "https://%s.slack.com/api/users.admin.%s?t=%s"

func init() {
	OVERRIDE_SLACK_HOST := os.Getenv("SLACK_HOST")
	var SLACK_HOST string
	if OVERRIDE_SLACK_HOST != nil {
		SLACK_HOST = OVERRIDE_SLACK_HOST
	} else {
		SLACK_HOST = "slack.com"
	}

	OVERRIDE_SLACK_SCHEME := os.Getenv("SLACK_SCHEME")
	var SLACK_SCHEME string
	if OVERRIDE_SLACK_SCHEME != nil {
		SLACK_SCHEME = OVERRIDE_SLACK_SCHEME
	} else {
		SLACK_SCHEME = "https"
	}

	SLACK_API = fmt.Sprintf("%s://%s/api/", SLACK_SCHEME, SLACK_HOST)
	SLACK_WEBSOCKET_REFERER = fmt.Sprintf("%s://%s", SLACK_SCHEME, SLACK_HOST)
}

type SlackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type AuthTestResponse struct {
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type authTestResponseFull struct {
	SlackResponse
	AuthTestResponse
}

type Client struct {
	config struct {
		token string
	}
	info  Info
	debug bool
}

// SetLogger let's library users supply a logger, so that api debugging
// can be logged along with the application's debugging info.
func SetLogger(l *log.Logger) {
	logger = l
}

func New(token string) *Client {
	s := &Client{}
	s.config.token = token
	return s
}

// AuthTest tests if the user is able to do authenticated requests or not
func (api *Client) AuthTest() (response *AuthTestResponse, error error) {
	responseFull := &authTestResponseFull{}
	err := post("auth.test", url.Values{"token": {api.config.token}}, responseFull, api.debug)
	if err != nil {
		return nil, err
	}
	if !responseFull.Ok {
		return nil, errors.New(responseFull.Error)
	}
	return &responseFull.AuthTestResponse, nil
}

// SetDebug switches the api into debug mode
// When in debug mode, it logs various info about what its doing
// If you ever use this in production, don't call SetDebug(true)
func (api *Client) SetDebug(debug bool) {
	api.debug = debug
	if debug && logger == nil {
		logger = log.New(os.Stdout, "nlopes/slack", log.LstdFlags | log.Lshortfile)
	}
}

func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		logger.Printf(format, v...)
	}
}

func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		logger.Println(v...)
	}
}
