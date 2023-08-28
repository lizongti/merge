package main

import (
	"encoding/json"
	"fmt"

	"github.com/cloudlibraries/merge"
)

var (
	s1 = `{"a":[{"b":1,"c":2},[3,4],null,5,6],"d":{"e":7,"f":"g","h":{"i":8,"j":9}},"k":[10],"l":[11], "m":[12]}`
	s2 = `{"a":[{"b":-1},[-3],-4],"d":{"e":-7,"f":-8},"k":null,"l":[]}`
	// s2 = `{"k":null}`.
)

func init() {
	fmt.Println("[s1]:", s1)
	fmt.Println("[s2]:", s2)
}

func main() {
	// makeExample("None")
	makeExample("Overwrite", merge.WithOverwrite())
	// makeExample("AppendSlice", merge.WithAppendSlice())
	// makeExample("Overwrite,AppendSlice", merge.WithOverwrite(), merge.WithAppendSlice())
	// makeExample("OverwriteWithEmptyValue", merge.WithOverwriteWithEmptyValue())
	makeExample("OverwriteEmptySlice", merge.WithOverwriteSliceWithEmptyValue())
}

func makeExample(name string, options ...merge.Option) {
	var m1, m2 map[string]interface{}
	if err := json.Unmarshal([]byte(s1), &m1); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(s2), &m2); err != nil {
		panic(err)
	}

	if err := merge.Merge(&m1, m2, options...); err != nil {
		panic(err)
	}

	if b1, err := json.Marshal(m1); err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("[%s]:", name), string(b1))
	}
}
