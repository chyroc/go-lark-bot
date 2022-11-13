package go_lark_bot

import (
	"context"
	"strings"

	"github.com/chyroc/lark"
)

type EventHeader struct {
	EventID    string         `json:"event_id,omitempty"`    // 事件 ID
	EventType  lark.EventType `json:"event_type,omitempty"`  // 事件类型
	CreateTime string         `json:"create_time,omitempty"` // 事件创建时间戳（单位：毫秒）
	Token      string         `json:"token,omitempty"`       // 事件 Token
	AppID      string         `json:"app_id,omitempty"`      // 应用 ID
	TenantKey  string         `json:"tenant_key,omitempty"`  // 租户 Key
}

type EventMessage struct {
	lark      *lark.Lark
	MsgType   lark.MsgType  `json:"msg_type,omitempty"` // 消息类型. 如: text
	RootID    string        `json:"root_id,omitempty"`
	ParentID  string        `json:"parent_id,omitempty"`
	ChatID    string        `json:"chat_id"`
	ChatType  lark.ChatType `json:"chat_type,omitempty"` // 私聊private, 群聊group. 如: private
	OpenID    string        `json:"open_id,omitempty"`   // 如: ou_18eac85d35a26f989317ad4f02e8bbbb
	UserID    string        `json:"user_id"`
	UnionID   string        `json:"union_id,omitempty"` // 如: xxx
	MessageID string        `json:"message_id"`
	Text      string        `json:"text,omitempty"`      // 消息文本, 可能包含被@的人/机器人。. 如: <at open_id="xxx">@小助手</at> 消息内容 <at open_id="yyy">@张三</at>
	Title     string        `json:"title,omitempty"`     // 富文本标题
	ImageKey  string        `json:"image_key,omitempty"` // 图片内容
	FileKey   string        `json:"file_key,omitempty"`  // 文件内容
	Args      []string      `json:"args"`
}

func (r *EventMessage) ReplyText(ctx context.Context, text string) (string, error) {
	res, _, err := r.lark.Message.Reply(r.MessageID).SendText(ctx, text)
	if err != nil {
		return "", err
	}
	return res.MessageID, nil
}

func (r *Client) makeMessageHandle() HandleMessage {
	return func(ctx context.Context, header *EventHeader, message *EventMessage) error {
		switch message.MsgType {
		case lark.MsgTypeText:
			for k, v := range r.textMatch {
				if k == message.Text {
					return v(ctx, header, message)
				}
			}
			for k, v := range r.textStart {
				if strings.HasPrefix(message.Text, k) {
					return v(ctx, header, message)
				}
			}
			for k, v := range r.textRegex {
				reg := r.textRegexList[k]
				match := reg.FindStringSubmatch(message.Text)
				if len(match) >= 1 {
					if len(match) >= 2 {
						message.Args = match[1:]
					}
					return v(ctx, header, message)
				}
			}
		}
		return nil
	}
}

func (r *Client) handleMessage() {
	f := r.makeMessageHandle()
	if r.receiveMessageV1 {
		r.larkClient.EventCallback.HandlerEventV1ReceiveMessage(func(ctx context.Context, cli *lark.Lark, schema string, header *lark.EventHeaderV1, event *lark.EventV1ReceiveMessage) (string, error) {
			return "", f(ctx, wrapHeader(header, nil), wrapV1Message(cli, event))
		})
	} else {
		r.larkClient.EventCallback.HandlerEventV2IMMessageReceiveV1(func(ctx context.Context, cli *lark.Lark, schema string, header *lark.EventHeaderV2, event *lark.EventV2IMMessageReceiveV1) (string, error) {
			content, err := lark.UnwrapMessageContent(event.Message.MessageType, event.Message.Content)
			if err != nil {
				return "", err
			}
			return "", f(ctx, wrapHeader(nil, header), wrapV2Message(cli, event, content))
		})
	}
}

func wrapV1Message(cli *lark.Lark, event *lark.EventV1ReceiveMessage) *EventMessage {
	return &EventMessage{
		lark:      cli,
		MsgType:   event.MsgType,
		RootID:    event.RootID,
		ParentID:  event.ParentID,
		ChatID:    event.OpenChatID,
		ChatType:  event.ChatType,
		OpenID:    event.OpenID,
		UserID:    event.EmployeeID,
		UnionID:   event.UnionID,
		MessageID: event.OpenMessageID,
		Text:      event.TextWithoutAtBot,
		Title:     event.Title,
		ImageKey:  event.ImageKey,
		FileKey:   event.FileKey,
	}
}

func wrapV2Message(cli *lark.Lark, event *lark.EventV2IMMessageReceiveV1, content *lark.MessageContent) *EventMessage {
	msg := event.Message
	res := &EventMessage{
		lark:      cli,
		MsgType:   msg.MessageType,
		RootID:    msg.RootID,
		ParentID:  msg.ParentID,
		ChatID:    msg.ChatID,
		ChatType:  msg.ChatType,
		MessageID: msg.MessageID,
	}
	if event.Sender != nil && event.Sender.SenderID != nil {
		res.OpenID = event.Sender.SenderID.OpenID
		res.UserID = event.Sender.SenderID.UserID
		res.UnionID = event.Sender.SenderID.UnionID
	}
	switch msg.MessageType {
	case lark.MsgTypeText:
		res.Text = content.Text.Text
	case lark.MsgTypePost:
		res.Title = content.Post.Title
	case lark.MsgTypeImage:
		res.ImageKey = content.Image.ImageKey
	case lark.MsgTypeFile:
		res.FileKey = content.File.FileKey
	}

	return res
}

func wrapHeader(v1 *lark.EventHeaderV1, v2 *lark.EventHeaderV2) *EventHeader {
	if v1 != nil {
		return &EventHeader{
			EventID:    v1.UUID,
			EventType:  v1.EventType,
			CreateTime: v1.TS,
			Token:      v1.Token,
			AppID:      "",
			TenantKey:  "",
		}
	} else {
		return &EventHeader{
			EventID:    v2.EventID,
			EventType:  v2.EventType,
			CreateTime: v2.CreateTime,
			Token:      v2.Token,
			AppID:      v2.AppID,
			TenantKey:  v2.TenantKey,
		}
	}
}
