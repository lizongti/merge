package merge

import "reflect"

func (m *merger) cover(dst reflect.Value, src reflect.Value) error {
	switch dst.Kind() {
	case reflect.Struct:
		return m.coverStruct(dst, src)
	case reflect.Map:
		return m.coverMap(dst, src)
	case reflect.Slice:
		return m.coverSlice(dst, src)
	}
	return nil
}

func (m *merger) coverStruct(dst reflect.Value, src reflect.Value) error {
	return nil
}


func (m *merger) coverMap(dst reflect.Value, src reflect.Value) error {
	return nil
}

func (m *merger) coverSlice(dst reflect.Value, src reflect.Value) error {
	return nil
}

