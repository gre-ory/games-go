package websocket

import (
	"bytes"
	"io"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

type TplRenderer[PlayerT any] interface {
	Render(player PlayerT) ([]byte, bool)
}

func NewTplRenderer[PlayerT any](acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) TplRenderer[PlayerT] {
	return &tplRenderer[PlayerT]{
		acceptFn: acceptFn,
		renderFn: renderFn,
	}
}

type tplRenderer[PlayerT any] struct {
	acceptFn func(player PlayerT) (bool, model.Data)
	renderFn func(w io.Writer, data model.Data)
}

func (t *tplRenderer[PlayerT]) Render(player PlayerT) ([]byte, bool) {
	if t.renderFn == nil {
		return nil, false
	}
	ok := true
	var data model.Data
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
