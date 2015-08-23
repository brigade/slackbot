package plugins

import (
	"fmt"
)

type PluginsContext struct {
	Plugins       []Plugin
	IsReply       bool
	MessageSender MessageSender
}

type MessageSender interface {
	SendMessage(message string, channel string)
}

func NewPluginsContext(sender MessageSender) *PluginsContext {
	return &PluginsContext{
		Plugins:       []Plugin{},
		IsReply:       true,
		MessageSender: sender,
	}
}

func (ctx *PluginsContext) AddPlugin(key interface{}, val BotMessagePlugin) {
	ctx.Plugins = append(ctx.Plugins, Plugin{key, val})
}

func (ctx *PluginsContext) StopReply() {
	ctx.IsReply = false
}

func (ctx *PluginsContext) StartReply() {
	ctx.IsReply = true
}

func (ctx *PluginsContext) ExecPlugins(message string, channel string) {
	e := NewBotEvent(ctx.MessageSender, message, channel)

	for _, p := range ctx.Plugins {
		ok, m := p.CheckMessage(*e, message)
		if !ok {
			continue
		}

		next := p.DoAction(*e, m)
		if !next {
			break
		}
	}
}

func (ctx *PluginsContext) SendMessage(message string, channel string) {
	if !ctx.IsReply {
		return
	}
	ctx.MessageSender.SendMessage(message, channel)
}

type BotMessagePlugin interface {
	CheckMessage(event BotEvent, message string) (bool, string)
	DoAction(event BotEvent, message string) bool
}

type Plugin struct {
	Key interface{}
	BotMessagePlugin
}

func (p Plugin) Name() string {
	return fmt.Sprintf("%s", p.Key)
}

type BotEvent struct {
	messageSender MessageSender
	text          string
	channel       string
}

func NewBotEvent(sender MessageSender, text, channel string) *BotEvent {
	return &BotEvent{
		messageSender: sender,
		text:          text,
		channel:       channel,
	}
}

func (b *BotEvent) Reply(message string) {
	b.SendMessage(message, b.Channel())
}

func (b *BotEvent) SendMessage(message string, channel string) {
	b.messageSender.SendMessage(message, channel)
}

func (b *BotEvent) BaseText() string {
	return b.text
}

func (b *BotEvent) Channel() string {
	return b.channel
}

var _ MessageSender = (*BotEvent)(nil)
