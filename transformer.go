package merge

import "reflect"

type Transformer interface {
	Transformer(reflect.Type) func(dst, src reflect.Value) error
}
