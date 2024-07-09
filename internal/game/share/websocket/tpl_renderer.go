package websocket

import (
	"bytes"
	"io"
)

type TplRenderer[PlayerT any] interface {
	Render(player PlayerT) ([]byte, bool)
}

func NewTplRenderer[PlayerT any](acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) TplRenderer[PlayerT] {
	return &tplRenderer[PlayerT]{
		acceptFn: acceptFn,
		renderFn: renderFn,
	}
}

type tplRenderer[PlayerT any] struct {
	name     string
	acceptFn func(player PlayerT) (bool, any)
	renderFn func(w io.Writer, data any)
}

func (t *tplRenderer[PlayerT]) Render(player PlayerT) ([]byte, bool) {
	if t.renderFn == nil {
		return nil, false
	}
	ok := true
	var data any
	if t.acceptFn != nil {
		ok, data = t.acceptFn(player)
	}
	if !ok {
		return nil, false
	}
	buf := &bytes.Buffer{}
	t.renderFn(buf, data)
	return buf.Bytes(), true
}
