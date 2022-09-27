package merge

import "reflect"

type Condition func(dst reflect.Value, src reflect.Value) bool

type Conditions []Condition

func (c Conditions) canCover(dst reflect.Value, src reflect.Value) bool {
	for _, condition := range c {
		if condition(dst, src) {
			return true
		}
	}
	return false
}

func ConditionCoverAll(dst reflect.Value, src reflect.Value) bool {
	return true
}

func ConditionCoverZero(dst reflect.Value, src reflect.Value) bool {
	return dst.IsZero()
}
