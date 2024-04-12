package websocket

import (
	"bytes"
	"io"
)

type Template[PlayerT any] interface {
	// GetName() string
	// Accept(player PlayerT) (bool, any)
	Render(player PlayerT) ([]byte, bool)
}

func NewTemplate[PlayerT any](acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) Template[PlayerT] {
	return &template[PlayerT]{
		acceptFn: acceptFn,
		renderFn: renderFn,
	}
}

type template[PlayerT any] struct {
	name     string
	acceptFn func(player PlayerT) (bool, any)
	renderFn func(w io.Writer, data any)
}

// func (t *template[PlayerT]) GetName() string {
// 	return t.name
// }

// func (t *template[PlayerT]) Accept(player PlayerT) (bool, any) {
// 	return t.acceptFn(player)
// }

func (t *template[PlayerT]) Render(player PlayerT) ([]byte, bool) {
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
