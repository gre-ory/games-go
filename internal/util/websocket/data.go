package websocket

type Data map[string]any

func (d Data) With(key string, value any) Data {
	if d == nil {
		d = make(Data)
	}
	d[key] = value
	return d
}
