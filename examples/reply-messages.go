package main

import (
	"context"

	go_lark_bot "github.com/chyroc/go-lark-bot"
	"github.com/chyroc/lark"
)

func main() {
	larkClient := lark.New(lark.WithAppCredential("", ""))
	example(context.Background(), larkClient)
}

func example(ctx context.Context, larkClient *lark.Lark) {
	r := go_lark_bot.NewClient(larkClient)

	r.Text("/hi", func(ctx context.Context, header *go_lark_bot.EventHeader, message *go_lark_bot.EventMessage) error {
		_, err := message.ReplyText(ctx, "hi")
		return err
	})

	r.TextRegex(`/translate (.*)`, func(ctx context.Context, header *go_lark_bot.EventHeader, message *go_lark_bot.EventMessage) error {
		_, err := message.ReplyText(ctx, "translate for: "+message.Args[0])
		return err
	})

	r.ListenCallback(ctx, nil, nil)
}
