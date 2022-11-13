package go_lark_bot

import (
	"context"
	"io"
	"net/http"
	"regexp"

	"github.com/chyroc/lark"
)

type Client struct {
	larkClient       *lark.Lark
	receiveMessageV1 bool

	textMatch     map[string]HandleMessage
	textStart     map[string]HandleMessage
	textRegex     map[string]HandleMessage
	textRegexList map[string]*regexp.Regexp
}

func NewClient(larkClient *lark.Lark) *Client {
	res := &Client{
		larkClient:       larkClient,
		receiveMessageV1: false,
		textMatch:        map[string]HandleMessage{},
		textStart:        map[string]HandleMessage{},
		textRegex:        map[string]HandleMessage{},
		textRegexList:    map[string]*regexp.Regexp{},
	}
	return res
}

func (r *Client) ListenCallback(ctx context.Context, reader io.Reader, writer http.ResponseWriter) {
	r.handleMessage()
	r.larkClient.EventCallback.ListenCallback(ctx, reader, writer)
}
