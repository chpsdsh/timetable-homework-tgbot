package telegram

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(ctx context.Context, u tgbotapi.Update)

type Router struct {
	cmd   map[string]Handler
	text  map[string]Handler
	state map[string]Handler
	def   Handler
}

func NewRouter() *Router {
	return &Router{
		cmd:   map[string]Handler{},
		text:  map[string]Handler{},
		state: map[string]Handler{},
		def:   func(context.Context, tgbotapi.Update) {},
	}
}

func (r *Router) OnCommand(cmd string, h Handler) { r.cmd[cmd] = h }
func (r *Router) OnText(txt string, h Handler)    { r.text[txt] = h }
func (r *Router) OnState(st string, h Handler)    { r.state[st] = h }
func (r *Router) Default(h Handler)               { r.def = h }

func matchState(r *Router, st string) (Handler, bool) {
	if h, ok := r.state[st]; ok {
		return h, true
	}
	for key, h := range r.state {
		if strings.HasPrefix(st, key) {
			return h, true
		}
	}
	return nil, false
}
