package util

import (
	"fmt"
	"html/template"
	"io"

	"go.uber.org/zap"
)

type TplRenderer interface {
	Render(w io.Writer, name string, data any)
}

func NewTplRenderer(logger *zap.Logger, tpl *template.Template) TplRenderer {
	return &tplRenderer{
		logger: logger,
		tpl:    tpl,
	}
}

type tplRenderer struct {
	logger *zap.Logger
	tpl    *template.Template
}

func (r *tplRenderer) Render(w io.Writer, name string, data any) {
	err := r.tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		r.logger.Warn(fmt.Sprintf("[render] %s: FAILED ( %s )", name, err), zap.Any("data", data), zap.Error(err))
	} else {
		r.logger.Info(fmt.Sprintf("[render] %s: ok", name), zap.Any("data", data))
	}
}
