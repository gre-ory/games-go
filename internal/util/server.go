package util

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Server interface {
	RegisterRoutes(router *httprouter.Router)
}

type HxServer interface {
	TplRenderer
	Redirect(w http.ResponseWriter, url string)
	RefreshTarget(w http.ResponseWriter, target string, url string)
	RenderInfo(w io.Writer, info string)
	RenderError(w io.Writer, err error)
}

func NewHxServer(logger *zap.Logger, tpl *template.Template) HxServer {
	return &hxServer{
		TplRenderer: NewTplRenderer(logger, tpl),
		logger:      logger,
	}
}

type hxServer struct {
	TplRenderer
	logger *zap.Logger
}

// //////////////////////////////////////////////////
// redirect

func (s *hxServer) Redirect(w http.ResponseWriter, url string) {
	// w.Header().Set("HX-Location", url)
	w.Header().Set("HX-Redirect", url)
	w.Header().Set("HX-Replace-Url", "true")
}

// //////////////////////////////////////////////////
// refresh target

func (s *hxServer) RefreshTarget(w http.ResponseWriter, target string, url string) {
	w.Header().Set("HX-Location", fmt.Sprintf("{\"path\":\"%s\", \"target\":\"%s\"}", url, target))
}

// //////////////////////////////////////////////////
// render

func (s *hxServer) RenderInfo(w io.Writer, info string) {
	s.Render(w, "info", map[string]any{
		"info": info,
	})
}

func (s *hxServer) RenderError(w io.Writer, err error) {
	s.Render(w, "error", map[string]any{
		"error": err.Error(),
	})
}
