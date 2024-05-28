package api

import (
	"bytes"
	"html/template"
)

type Renderer interface {
	LoadingDot() template.HTML
	UserBadge() template.HTML
	WsStatusBadge() template.HTML
}

func NewRenderer() Renderer {
	return &renderer{}
}

type renderer struct{}

func (r *renderer) LoadingDot() template.HTML {
	return r.render("loading-dot", nil)
}

func (r *renderer) UserBadge() template.HTML {
	return r.render("user-badge", nil)
}

func (r *renderer) WsStatusBadge() template.HTML {
	return r.render("ws-status-badge", nil)
}

func (r *renderer) render(name string, data any) template.HTML {
	w := &bytes.Buffer{}
	ShareTpl.ExecuteTemplate(w, name, data)
	return template.HTML(w.String())
}
