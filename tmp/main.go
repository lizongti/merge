package main

import (
	"fmt"
	"reflect"

	"github.com/cloudlibraries/merge"
)

func main() {
	var a int = 10
	var b int = 1

	ret := merge.MustMerge(&a, &b,
		merge.WithResolver(merge.ResolverBoth),
		merge.WithCondition(merge.ConditionCoverAll))
	v := reflect.ValueOf(ret)
	fmt.Println(v.Kind(), v.Elem().Kind())
	
}
