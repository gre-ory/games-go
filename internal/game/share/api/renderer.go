package api

import (
	"bytes"
	"html/template"
)

type Renderer interface {
	LoadingDot() template.HTML
	UserBadge() template.HTML
	WsStatusBadge() template.HTML
	Render(name string, data any) template.HTML
}

func NewRenderer() Renderer {
	return &renderer{}
}

type renderer struct{}

func (r *renderer) LoadingDot() template.HTML {
	return r.Render("loading-dot", nil)
}

func (r *renderer) UserBadge() template.HTML {
	return r.Render("user-badge", nil)
}

func (r *renderer) WsStatusBadge() template.HTML {
	return r.Render("ws-status-badge", nil)
}

func (r *renderer) Render(name string, data any) template.HTML {
	w := &bytes.Buffer{}
	ShareTpl.ExecuteTemplate(w, name, data)
	return template.HTML(w.String())
}
