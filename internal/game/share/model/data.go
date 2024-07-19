package model

type Data map[string]any

func (d Data) With(key string, value any) Data {
	if d == nil {
		return make(Data)
	}
	d[key] = value
	return d
}

func (d Data) Get(key string) any {
	return d[key]
}

func (d Data) Set(key string, value any) {
	d[key] = value
}

func (d Data) Delete(key string) {
	delete(d, key)
}
