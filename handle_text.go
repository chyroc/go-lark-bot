package go_lark_bot

import (
	"context"
	"regexp"
)

type HandleMessage func(ctx context.Context, header *EventHeader, message *EventMessage) error

func (r *Client) Text(textText string, f HandleMessage) {
	r.textMatch[textText] = f
}

func (r *Client) TextStart(text string, f HandleMessage) {
	r.textStart[text] = f
}

func (r *Client) TextRegex(text string, f HandleMessage) {
	r.textRegex[text] = f
	r.textRegexList[text] = regexp.MustCompile(text)
}
