package merge

import "reflect"

var conditionsDefault = Conditions{
	ConditionSrcIsValid,
}

type Condition func(dst reflect.Value, src reflect.Value) bool

type Conditions []Condition

func newConditions() Conditions {
	return append([]Condition{}, conditionsDefault...)
}

func (c Conditions) Check(dst reflect.Value, src reflect.Value) bool {
	for _, condition := range c {
		if !condition(dst, src) {
			return false
		}
	}
	return true
}

func ConditionDstIsZero(dst reflect.Value, src reflect.Value) bool {
	return dst.IsZero()
}

func ConditionTypeCheck(dst reflect.Value, src reflect.Value) bool {
	return dst.Type() == src.Type()
}

func ConditionSrcIsValid(dst reflect.Value, src reflect.Value) bool {
	return src.IsValid()
}
