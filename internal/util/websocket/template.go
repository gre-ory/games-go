package websocket

type Template[PlayerT any] interface {
	GetName() string
	Accept(player PlayerT) (bool, any)
}

func NewTemplate[PlayerT any](name string, acceptFn func(player PlayerT) (bool, any)) Template[PlayerT] {
	return &template[PlayerT]{
		name:     name,
		acceptFn: acceptFn,
	}
}

type template[PlayerT any] struct {
	name     string
	acceptFn func(player PlayerT) (bool, any)
}

func (t *template[PlayerT]) GetName() string {
	return t.name
}

func (t *template[PlayerT]) Accept(player PlayerT) (bool, any) {
	return t.acceptFn(player)
}
