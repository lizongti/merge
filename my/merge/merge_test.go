package merge_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/cloudlibraries/merge"
	"gopkg.in/yaml.v3"
)

type simpleTest struct {
	Value int
}

type complexTest struct {
	ID string
	St simpleTest
	sz int
}

type mapTest struct {
	M map[int]int
}

type ifcTest struct {
	I interface{}
}

type moreComplextText struct {
	Ct complexTest
	St simpleTest
	Nt simpleTest
}

type pointerTest struct {
	C *simpleTest
}

type sliceTest struct {
	S []int
}

func TestKb(t *testing.T) {
	type testStruct struct {
		KeyValue map[string]interface{}
		Name     string
	}

	akv := make(map[string]interface{})
	akv["Key1"] = "not value 1"
	akv["Key2"] = "value2"
	a := testStruct{}
	a.Name = "A"
	a.KeyValue = akv

	bkv := make(map[string]interface{})
	bkv["Key1"] = "value1"
	bkv["Key3"] = "value3"
	b := testStruct{}
	b.Name = "B"
	b.KeyValue = bkv

	ekv := make(map[string]interface{})
	ekv["Key1"] = "value1"
	ekv["Key2"] = "value2"
	ekv["Key3"] = "value3"
	expected := testStruct{}
	expected.Name = "B"
	expected.KeyValue = ekv

	if err := merge.Merge(&b, a); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(b, expected) {
		t.Errorf("Actual: %#v did not match \nExpected: %#v", b, expected)
	}
}

func TestNil(t *testing.T) {
	if err := merge.Merge(nil, nil); err != merge.ErrNilArguments {
		t.Fail()
	}
}

func TestDifferentTypes(t *testing.T) {
	a := simpleTest{42}
	b := 42
	if err := merge.Merge(&a, b); err != merge.ErrDifferentArgumentsTypes {
		t.Fail()
	}
}

func TestSimpleStruct(t *testing.T) {
	a := simpleTest{}
	b := simpleTest{42}
	if err := merge.Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.Value != 42 {
		t.Errorf("b not merged in properly: a.Value(%d) != b.Value(%d)", a.Value, b.Value)
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

func TestComplexStruct(t *testing.T) {
	a := complexTest{}
	a.ID = "athing"
	b := complexTest{"bthing", simpleTest{42}, 1}
	if err := merge.Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.St.Value != 42 {
		t.Errorf("b not merged in properly: a.St.Value(%d) != b.St.Value(%d)", a.St.Value, b.St.Value)
	}
	if a.sz == 1 {
		t.Errorf("a's private field sz not preserved from merge: a.sz(%d) == b.sz(%d)", a.sz, b.sz)
	}
	if a.ID == b.ID {
		t.Errorf("a's field ID merged unexpectedly: a.ID(%s) == b.ID(%s)", a.ID, b.ID)
	}
}

func TestComplexStructWithOverwrite(t *testing.T) {
	a := complexTest{"do-not-overwrite-with-empty-value", simpleTest{1}, 1}
	b := complexTest{"", simpleTest{42}, 2}

	expect := complexTest{"do-not-overwrite-with-empty-value", simpleTest{42}, 1}
	if err := merge.Merge(&a, b, merge.WithOverwrite()); err != nil {
		t.FailNow()
	}

	if !reflect.DeepEqual(a, expect) {
		t.Errorf("Test failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", a, expect)
	}
}

func TestPointerStruct(t *testing.T) {
	s1 := simpleTest{}
	s2 := simpleTest{19}
	a := pointerTest{&s1}
	b := pointerTest{&s2}
	if err := merge.Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.C.Value != b.C.Value {
		t.Errorf("b not merged in properly: a.C.Value(%d) != b.C.Value(%d)", a.C.Value, b.C.Value)
	}
}

type embeddingStruct struct {
	embeddedStruct
}

type embeddedStruct struct {
	A string
}

func TestEmbeddedStruct(t *testing.T) {
	tests := []struct {
		src      embeddingStruct
		dst      embeddingStruct
		expected embeddingStruct
	}{
		{
			src: embeddingStruct{
				embeddedStruct{"foo"},
			},
			dst: embeddingStruct{
				embeddedStruct{""},
			},
			expected: embeddingStruct{
				embeddedStruct{"foo"},
			},
		},
		{
			src: embeddingStruct{
				embeddedStruct{""},
			},
			dst: embeddingStruct{
				embeddedStruct{"bar"},
			},
			expected: embeddingStruct{
				embeddedStruct{"bar"},
			},
		},
		{
			src: embeddingStruct{
				embeddedStruct{"foo"},
			},
			dst: embeddingStruct{
				embeddedStruct{"bar"},
			},
			expected: embeddingStruct{
				embeddedStruct{"bar"},
			},
		},
	}

	for _, test := range tests {
		err := merge.Merge(&test.dst, test.src)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if !reflect.DeepEqual(test.dst, test.expected) {
			t.Errorf("unexpected output\nexpected:\n%+v\nsaw:\n%+v\n", test.expected, test.dst)
		}
	}
}

func TestPointerStructNil(t *testing.T) {
	a := pointerTest{nil}
	b := pointerTest{&simpleTest{19}}
	if err := merge.Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.C.Value != b.C.Value {
		t.Errorf("b not merged in a properly: a.C.Value(%d) != b.C.Value(%d)", a.C.Value, b.C.Value)
	}
}

func testSlice(t *testing.T, a []int, b []int, e []int, opts ...merge.Option) {
	t.Helper()
	bc := b

	sa := sliceTest{a}
	sb := sliceTest{b}
	if err := merge.Merge(&sa, sb, opts...); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(sb.S, bc) {
		t.Errorf("Source slice was modified %d != %d", sb.S, bc)
	}
	if !reflect.DeepEqual(sa.S, e) {
		t.Errorf("b not merged in a proper way %d != %d", sa.S, e)
	}

	ma := map[string][]int{"S": a}
	mb := map[string][]int{"S": b}
	if err := merge.Merge(&ma, mb, opts...); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(mb["S"], bc) {
		t.Errorf("map value: Source slice was modified %d != %d", mb["S"], bc)
	}
	if !reflect.DeepEqual(ma["S"], e) {
		t.Errorf("map value: b not merged in a proper way %d != %d", ma["S"], e)
	}

	if a == nil {
		// test case with missing dst key
		ma := map[string][]int{}
		mb := map[string][]int{"S": b}
		if err := merge.Merge(&ma, mb); err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(mb["S"], bc) {
			t.Errorf("missing dst key: Source slice was modified %d != %d", mb["S"], bc)
		}
		if !reflect.DeepEqual(ma["S"], e) {
			t.Errorf("missing dst key: b not merged in a proper way %d != %d", ma["S"], e)
		}
	}

	if b == nil {
		// test case with missing src key
		ma := map[string][]int{"S": a}
		mb := map[string][]int{}
		if err := merge.Merge(&ma, mb); err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(mb["S"], bc) {
			t.Errorf("missing src key: Source slice was modified %d != %d", mb["S"], bc)
		}
		if !reflect.DeepEqual(ma["S"], e) {
			t.Errorf("missing src key: b not merged in a proper way %d != %d", ma["S"], e)
		}
	}
}

func TestSlice(t *testing.T) {
	testSlice(t, nil, []int{1, 2, 3}, []int{1, 2, 3})
	testSlice(t, []int{}, []int{1, 2, 3}, []int{1, 2, 3})
	testSlice(t, []int{1}, []int{2, 3}, []int{1})
	testSlice(t, []int{1}, []int{}, []int{1})
	testSlice(t, []int{1}, nil, []int{1})
	testSlice(t, nil, []int{1, 2, 3}, []int{1, 2, 3}, merge.WithAppendSlice())
	testSlice(t, []int{}, []int{1, 2, 3}, []int{1, 2, 3}, merge.WithAppendSlice())
	testSlice(t, []int{1}, []int{2, 3}, []int{1, 2, 3}, merge.WithAppendSlice())
	testSlice(t, []int{1}, []int{2, 3}, []int{1, 2, 3}, merge.WithAppendSlice(), merge.WithOverwrite())
	testSlice(t, []int{1}, []int{}, []int{1}, merge.WithAppendSlice())
	testSlice(t, []int{1}, nil, []int{1}, merge.WithAppendSlice())
}

func TestEmptyMaps(t *testing.T) {
	a := mapTest{}
	b := mapTest{
		map[int]int{},
	}
	if err := merge.Merge(&a, b); err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

func TestEmptyToEmptyMaps(t *testing.T) {
	a := mapTest{}
	b := mapTest{}
	if err := merge.Merge(&a, b); err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

func TestEmptyToNotEmptyMaps(t *testing.T) {
	a := mapTest{map[int]int{
		1: 2,
		3: 4,
	}}
	aa := mapTest{map[int]int{
		1: 2,
		3: 4,
	}}
	b := mapTest{
		map[int]int{},
	}
	if err := merge.Merge(&a, b); err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(a, aa) {
		t.FailNow()
	}
}

func TestMapsWithOverwrite(t *testing.T) {
	m := map[string]simpleTest{
		"a": {},   // overwritten by 16
		"b": {42}, // overwritten by 0, as map Value is not addressable and it doesn't check for b is set or not set in `n`
		"c": {13}, // overwritten by 12
		"d": {61},
	}
	n := map[string]simpleTest{
		"a": {16},
		"b": {},
		"c": {12},
		"e": {14},
	}
	expect := map[string]simpleTest{
		"a": {16},
		"b": {},
		"c": {12},
		"d": {61},
		"e": {14},
	}

	if err := merge.Merge(&m, n, merge.WithOverwrite()); err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(m, expect) {
		t.Errorf("Test failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", m, expect)
	}
}

func TestMapWithEmbeddedStructPointer(t *testing.T) {
	m := map[string]*simpleTest{
		"a": {},   // overwritten by 16
		"b": {42}, // not overwritten by empty value
		"c": {13}, // overwritten by 12
		"d": {61},
	}
	n := map[string]*simpleTest{
		"a": {16},
		"b": {},
		"c": {12},
		"e": {14},
	}
	expect := map[string]*simpleTest{
		"a": {16},
		"b": {42},
		"c": {12},
		"d": {61},
		"e": {14},
	}

	if err := merge.Merge(&m, n, merge.WithOverwrite()); err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(m, expect) {
		t.Errorf("Test failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", m, expect)
	}
}

func TestMergeUsingStructAndMap(t *testing.T) {
	type multiPtr struct {
		Text   string
		Number int
	}
	type final struct {
		Msg1 string
		Msg2 string
	}
	type params struct {
		Multi *multiPtr
		Final *final
		Name  string
	}
	type config struct {
		Params *params
		Foo    string
		Bar    string
	}

	cases := []struct {
		changes   *config
		target    *config
		output    *config
		name      string
		overwrite bool
	}{
		{
			name:      "Should overwrite values in target for non-nil values in source",
			overwrite: true,
			changes: &config{
				Bar: "from changes",
				Params: &params{
					Final: &final{
						Msg1: "from changes",
						Msg2: "from changes",
					},
				},
			},
			target: &config{
				Foo: "from target",
				Params: &params{
					Name: "from target",
					Multi: &multiPtr{
						Text:   "from target",
						Number: 5,
					},
					Final: &final{
						Msg1: "from target",
						Msg2: "",
					},
				},
			},
			output: &config{
				Foo: "from target",
				Bar: "from changes",
				Params: &params{
					Name: "from target",
					Multi: &multiPtr{
						Text:   "from target",
						Number: 5,
					},
					Final: &final{
						Msg1: "from changes",
						Msg2: "from changes",
					},
				},
			},
		},
		{
			name:      "Should not overwrite values in target for non-nil values in source",
			overwrite: false,
			changes: &config{
				Bar: "from changes",
				Params: &params{
					Final: &final{
						Msg1: "from changes",
						Msg2: "from changes",
					},
				},
			},
			target: &config{
				Foo: "from target",
				Params: &params{
					Name: "from target",
					Multi: &multiPtr{
						Text:   "from target",
						Number: 5,
					},
					Final: &final{
						Msg1: "from target",
						Msg2: "",
					},
				},
			},
			output: &config{
				Foo: "from target",
				Bar: "from changes",
				Params: &params{
					Name: "from target",
					Multi: &multiPtr{
						Text:   "from target",
						Number: 5,
					},
					Final: &final{
						Msg1: "from target",
						Msg2: "from changes",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			if tc.overwrite {
				err = merge.Merge(tc.target, *tc.changes, merge.WithOverwrite())
			} else {
				err = merge.Merge(tc.target, *tc.changes)
			}
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tc.target, tc.output) {
				t.Errorf("Test failed:\ngot  :\n%+v\n\nwant :\n%+v\n\n", tc.target.Params, tc.output.Params)
			}
		})
	}
}
func TestMaps(t *testing.T) {
	m := map[string]simpleTest{
		"a": {},
		"b": {42},
		"c": {13},
		"d": {61},
	}
	n := map[string]simpleTest{
		"a": {16},
		"b": {},
		"c": {12},
		"e": {14},
	}
	expect := map[string]simpleTest{
		"a": {0},
		"b": {42},
		"c": {13},
		"d": {61},
		"e": {14},
	}

	if err := merge.Merge(&m, n); err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(m, expect) {
		t.Errorf("Test failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", m, expect)
	}
	if m["a"].Value != 0 {
		t.Errorf(`n merged in m because I solved non-addressable map values TODO: m["a"].Value(%d) != n["a"].Value(%d)`, m["a"].Value, n["a"].Value)
	}
	if m["b"].Value != 42 {
		t.Errorf(`n wrongly merged in m: m["b"].Value(%d) != n["b"].Value(%d)`, m["b"].Value, n["b"].Value)
	}
	if m["c"].Value != 13 {
		t.Errorf(`n overwritten in m: m["c"].Value(%d) != n["c"].Value(%d)`, m["c"].Value, n["c"].Value)
	}
}

func TestMapsWithNilPointer(t *testing.T) {
	m := map[string]*simpleTest{
		"a": nil,
		"b": nil,
	}
	n := map[string]*simpleTest{
		"b": nil,
		"c": nil,
	}
	expect := map[string]*simpleTest{
		"a": nil,
		"b": nil,
		"c": nil,
	}

	if err := merge.Merge(&m, n, merge.WithOverwrite()); err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(m, expect) {
		t.Errorf("Test failed:\ngot   :\n%#v\n\nwant :\n%#v\n\n", m, expect)
	}
}

func TestTwoPointerValues(t *testing.T) {
	a := &simpleTest{}
	b := &simpleTest{42}
	if err := merge.Merge(a, b); err != nil {
		t.Errorf(`Boom. You crossed the streams: %s`, err)
	}
}

func TestMap(t *testing.T) {
	a := complexTest{}
	a.ID = "athing"
	c := moreComplextText{a, simpleTest{}, simpleTest{}}
	b := map[string]interface{}{
		"ct": map[string]interface{}{
			"st": map[string]interface{}{
				"value": 42,
			},
			"sz": 1,
			"id": "bthing",
		},
		"st": &simpleTest{144}, // Mapping a reference
		"zt": simpleTest{299},  // Mapping a missing field (zt doesn't exist)
		"nt": simpleTest{3},
	}
	if err := merge.Map(&c, b); err != nil {
		t.FailNow()
	}
	m := b["ct"].(map[string]interface{})
	n := m["st"].(map[string]interface{})
	o := b["st"].(*simpleTest)
	p := b["nt"].(simpleTest)
	if c.Ct.St.Value != 42 {
		t.Errorf("b not merged in properly: c.Ct.St.Value(%d) != b.Ct.St.Value(%d)", c.Ct.St.Value, n["value"])
	}
	if c.St.Value != 144 {
		t.Errorf("b not merged in properly: c.St.Value(%d) != b.St.Value(%d)", c.St.Value, o.Value)
	}
	if c.Nt.Value != 3 {
		t.Errorf("b not merged in properly: c.Nt.Value(%d) != b.Nt.Value(%d)", c.St.Value, p.Value)
	}
	if c.Ct.sz == 1 {
		t.Errorf("a's private field sz not preserved from merge: c.Ct.sz(%d) == b.Ct.sz(%d)", c.Ct.sz, m["sz"])
	}
	if c.Ct.ID == m["id"] {
		t.Errorf("a's field ID merged unexpectedly: c.Ct.ID(%s) == b.Ct.ID(%s)", c.Ct.ID, m["id"])
	}
}

func TestSimpleMap(t *testing.T) {
	a := simpleTest{}
	b := map[string]interface{}{
		"value": 42,
	}
	if err := merge.Map(&a, b); err != nil {
		t.FailNow()
	}
	if a.Value != 42 {
		t.Errorf("b not merged in properly: a.Value(%d) != b.Value(%v)", a.Value, b["value"])
	}
}

func TestIfcMap(t *testing.T) {
	a := ifcTest{}
	b := ifcTest{42}
	if err := merge.Map(&a, b); err != nil {
		t.FailNow()
	}
	if a.I != 42 {
		t.Errorf("b not merged in properly: a.I(%d) != b.I(%d)", a.I, b.I)
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

func TestIfcMapNoOverwrite(t *testing.T) {
	a := ifcTest{13}
	b := ifcTest{42}
	if err := merge.Map(&a, b); err != nil {
		t.FailNow()
	}
	if a.I != 13 {
		t.Errorf("a not left alone: a.I(%d) == b.I(%d)", a.I, b.I)
	}
}

func TestIfcMapWithOverwrite(t *testing.T) {
	a := ifcTest{13}
	b := ifcTest{42}
	if err := merge.Map(&a, b, merge.WithOverwrite()); err != nil {
		t.FailNow()
	}
	if a.I != 42 {
		t.Errorf("b not merged in properly: a.I(%d) != b.I(%d)", a.I, b.I)
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

type pointerMapTest struct {
	B      *simpleTest
	A      int
	hidden int
}

func TestBackAndForth(t *testing.T) {
	pt := pointerMapTest{&simpleTest{66}, 42, 1}
	m := make(map[string]interface{})
	if err := merge.Map(&m, pt); err != nil {
		t.FailNow()
	}
	var (
		v  interface{}
		ok bool
	)
	if v, ok = m["a"]; v.(int) != pt.A || !ok {
		t.Errorf("pt not merged in properly: m[`a`](%d) != pt.A(%d)", v, pt.A)
	}
	if v, ok = m["b"]; !ok {
		t.Errorf("pt not merged in properly: B is missing in m")
	}
	var st *simpleTest
	if st = v.(*simpleTest); st.Value != 66 {
		t.Errorf("something went wrong while mapping pt on m, B wasn't copied")
	}
	bpt := pointerMapTest{}
	if err := merge.Map(&bpt, m); err != nil {
		t.Error(err)
	}
	if bpt.A != pt.A {
		t.Errorf("pt not merged in properly: bpt.A(%d) != pt.A(%d)", bpt.A, pt.A)
	}
	if bpt.hidden == pt.hidden {
		t.Errorf("pt unexpectedly merged: bpt.hidden(%d) == pt.hidden(%d)", bpt.hidden, pt.hidden)
	}
	if bpt.B.Value != pt.B.Value {
		t.Errorf("pt not merged in properly: bpt.B.Value(%d) != pt.B.Value(%d)", bpt.B.Value, pt.B.Value)
	}
}

func TestEmbeddedPointerUnpacking(t *testing.T) {
	tests := []struct{ input pointerMapTest }{
		{pointerMapTest{nil, 42, 1}},
		{pointerMapTest{&simpleTest{66}, 42, 1}},
	}
	newValue := 77
	m := map[string]interface{}{
		"b": map[string]interface{}{
			"value": newValue,
		},
	}
	for _, test := range tests {
		pt := test.input
		if err := merge.Map(&pt, m, merge.WithOverwrite()); err != nil {
			t.FailNow()
		}
		if pt.B.Value != newValue {
			t.Errorf("pt not mapped properly: pt.A.Value(%d) != m[`b`][`value`](%d)", pt.B.Value, newValue)
		}

	}
}

type structWithTimePointer struct {
	Birth *time.Time
}

func TestTime(t *testing.T) {
	now := time.Now()
	dataStruct := structWithTimePointer{
		Birth: &now,
	}
	dataMap := map[string]interface{}{
		"Birth": &now,
	}
	b := structWithTimePointer{}
	if err := merge.Merge(&b, dataStruct); err != nil {
		t.FailNow()
	}
	if b.Birth.IsZero() {
		t.Errorf("time.Time not merged in properly: b.Birth(%v) != dataStruct['Birth'](%v)", b.Birth, dataStruct.Birth)
	}
	if b.Birth != dataStruct.Birth {
		t.Errorf("time.Time not merged in properly: b.Birth(%v) != dataStruct['Birth'](%v)", b.Birth, dataStruct.Birth)
	}
	b = structWithTimePointer{}
	if err := merge.Map(&b, dataMap); err != nil {
		t.FailNow()
	}
	if b.Birth.IsZero() {
		t.Errorf("time.Time not merged in properly: b.Birth(%v) != dataMap['Birth'](%v)", b.Birth, dataMap["Birth"])
	}
}

type simpleNested struct {
	A int
}

type structWithNestedPtrValueMap struct {
	NestedPtrValue map[string]*simpleNested
}

func TestNestedPtrValueInMap(t *testing.T) {
	src := &structWithNestedPtrValueMap{
		NestedPtrValue: map[string]*simpleNested{
			"x": {
				A: 1,
			},
		},
	}
	dst := &structWithNestedPtrValueMap{
		NestedPtrValue: map[string]*simpleNested{
			"x": {},
		},
	}
	if err := merge.Map(dst, src); err != nil {
		t.FailNow()
	}
	if dst.NestedPtrValue["x"].A == 0 {
		t.Errorf("Nested Ptr value not merged in properly: dst.NestedPtrValue[\"x\"].A(%v) != src.NestedPtrValue[\"x\"].A(%v)", dst.NestedPtrValue["x"].A, src.NestedPtrValue["x"].A)
	}
}

func loadYAML(path string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	raw, _ := ioutil.ReadFile(path)
	_ = yaml.Unmarshal(raw, &m)
	return
}

type structWithMap struct {
	m map[string]structWithUnexportedProperty
}

type structWithUnexportedProperty struct {
	s string
}

func TestUnexportedProperty(t *testing.T) {
	a := structWithMap{map[string]structWithUnexportedProperty{
		"key": {"hello"},
	}}
	b := structWithMap{map[string]structWithUnexportedProperty{
		"key": {"hi"},
	}}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not have panicked")
		}
	}()
	merge.Merge(&a, b)
}

type structWithBoolPointer struct {
	C *bool
}

func TestBooleanPointer(t *testing.T) {
	bt, bf := true, false
	src := structWithBoolPointer{
		&bt,
	}
	dst := structWithBoolPointer{
		&bf,
	}
	if err := merge.Merge(&dst, src); err != nil {
		t.FailNow()
	}
	if dst.C == src.C {
		t.Errorf("dst.C should be a different pointer than src.C")
	}
	if *dst.C != *src.C {
		t.Errorf("dst.C should be true")
	}
}

func TestMergeDifferentSlicesIsNotSupported(t *testing.T) {
	src := []string{"a", "b"}
	dst := []int{1, 2}

	if err := merge.Merge(&src, &dst, merge.WithOverwrite(), merge.WithAppendSlice()); err != merge.ErrDifferentArgumentsTypes {
		t.Errorf("expected %q, got %q", merge.ErrNotSupported, err)
	}
}

type transformer struct {
	m map[reflect.Type]func(dst, src reflect.Value) error
}

func (s *transformer) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if fn, ok := s.m[t]; ok {
		return fn
	}
	return nil
}

type foo struct {
	Bar *bar
	s   string
}

type bar struct {
	s map[string]string
	i int
}

func TestMergeWithTransformerNilStruct(t *testing.T) {
	a := foo{s: "foo"}
	b := foo{Bar: &bar{i: 2, s: map[string]string{"foo": "bar"}}}

	if err := merge.Merge(&a, &b, merge.WithOverwrite(), merge.WithTransformers(&transformer{
		m: map[reflect.Type]func(dst, src reflect.Value) error{
			reflect.TypeOf(&bar{}): func(dst, src reflect.Value) error {
				// Do sthg with Elem
				t.Log(dst.Elem().FieldByName("i"))
				t.Log(src.Elem())
				return nil
			},
		},
	})); err != nil {
		t.Error(err)
	}

	if a.s != "foo" {
		t.Errorf("b not merged in properly: a.s.Value(%s) != expected(%s)", a.s, "foo")
	}

	if a.Bar == nil {
		t.Errorf("b not merged in properly: a.Bar shouldn't be nil")
	}
}

func TestMergeNonPointer(t *testing.T) {
	dst := bar{
		i: 1,
	}
	src := bar{
		i: 2,
		s: map[string]string{
			"a": "1",
		},
	}
	want := merge.ErrNonPointerAgument

	if got := merge.Merge(dst, src); got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestMapNonPointer(t *testing.T) {
	dst := make(map[string]bar)
	src := map[string]bar{
		"a": {
			i: 2,
			s: map[string]string{
				"a": "1",
			},
		},
	}
	want := merge.ErrNonPointerAgument
	if got := merge.Merge(dst, src); got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestYAMLMaps(t *testing.T) {
	thing := loadYAML("testdata/thing.yml")
	license := loadYAML("testdata/license.yml")
	ft := thing["fields"].(map[string]interface{})
	fl := license["fields"].(map[string]interface{})
	// license has one extra field (site) and another already existing in thing (author) that merge won't overwrite.
	expectedLength := len(ft) + len(fl) - 1
	if err := merge.Merge(&license, thing); err != nil {
		t.Error(err.Error())
	}
	currentLength := len(license["fields"].(map[string]interface{}))
	if currentLength != expectedLength {
		t.Errorf(`thing not merged in license properly, license must have %d elements instead of %d`, expectedLength, currentLength)
	}
	fields := license["fields"].(map[string]interface{})
	if _, ok := fields["id"]; !ok {
		t.Errorf(`thing not merged in license properly, license must have a new id field from thing`)
	}
}

func TestIssue17MergeWithOverwrite(t *testing.T) {
	var (
		request    = `{"timestamp":null, "name": "foo"}`
		maprequest = map[string]interface{}{
			"timestamp": nil,
			"name":      "foo",
			"newStuff":  "foo",
		}
	)

	var something map[string]interface{}
	if err := json.Unmarshal([]byte(request), &something); err != nil {
		t.Errorf("Error while Unmarshalling maprequest: %s", err)
	}

	if err := merge.Merge(&something, maprequest, merge.WithOverwrite()); err != nil {
		t.Errorf("Error while merging: %s", err)
	}
}

type issue23Document struct {
	Created *time.Time
}

func TestIssue23MergeWithOverwrite(t *testing.T) {
	now := time.Now()
	dst := issue23Document{
		&now,
	}
	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := issue23Document{
		&expected,
	}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !dst.Created.Equal(*src.Created) { //--> https://golang.org/pkg/time/#pkg-overview
		t.Errorf("Created not merged in properly: dst.Created(%v) != src.Created(%v)", dst.Created, src.Created)
	}
}

type issue33Foo struct {
	Str    string
	Bslice []byte
}

func TestIssue33Merge(t *testing.T) {
	dest := issue33Foo{Str: "a"}
	toMerge := issue33Foo{
		Str:    "b",
		Bslice: []byte{1, 2},
	}

	if err := merge.Merge(&dest, toMerge); err != nil {
		t.Errorf("Error while merging: %s", err)
	}
	// Merge doesn't overwrite an attribute if in destination it doesn't have a zero value.
	// In this case, Str isn't a zero value string.
	if dest.Str != "a" {
		t.Errorf("dest.Str should have not been override as it has a non-zero value: dest.Str(%v) != 'a'", dest.Str)
	}
	// If we want to override, we must use MergeWithOverwrite or Merge using WithOverride.
	if err := merge.Merge(&dest, toMerge, merge.WithOverwrite()); err != nil {
		t.Errorf("Error while merging: %s", err)
	}

	if dest.Str != toMerge.Str {
		t.Errorf("dest.Str should have been override: dest.Str(%v) != toMerge.Str(%v)", dest.Str, toMerge.Str)
	}
}

type issue38StructWithoutTimePointer struct {
	Created time.Time
}

func TestIssue38Merge(t *testing.T) {
	dst := issue38StructWithoutTimePointer{
		time.Now(),
	}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := issue38StructWithoutTimePointer{
		expected,
	}

	if err := merge.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Created == src.Created {
		t.Errorf("Created merged unexpectedly: dst.Created(%v) == src.Created(%v)", dst.Created, src.Created)
	}
}

func TestIssue38MergeEmptyStruct(t *testing.T) {
	dst := issue38StructWithoutTimePointer{}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := issue38StructWithoutTimePointer{
		expected,
	}

	if err := merge.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Created == src.Created {
		t.Errorf("Created merged unexpectedly: dst.Created(%v) == src.Created(%v)", dst.Created, src.Created)
	}
}

func TestIssue38MergeWithOverwrite(t *testing.T) {
	dst := issue38StructWithoutTimePointer{
		time.Now(),
	}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := issue38StructWithoutTimePointer{
		expected,
	}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Created != src.Created {
		t.Errorf("Created not merged in properly: dst.Created(%v) != src.Created(%v)", dst.Created, src.Created)
	}
}

type issue50TestStruct struct {
	time.Duration
}

func TestIssue50Merge(t *testing.T) {
	to := issue50TestStruct{}
	from := issue50TestStruct{}

	if err := merge.Merge(&to, from); err != nil {
		t.Fail()
	}
}

type issue52StructWithTime struct {
	Birth time.Time
}

type issue52TimeTransfomer struct {
	overwrite bool
}

func (t issue52TimeTransfomer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if t.overwrite {
					isZero := src.MethodByName("IsZero")

					result := isZero.Call([]reflect.Value{})
					if !result[0].Bool() {
						dst.Set(src)
					}
				} else {
					isZero := dst.MethodByName("IsZero")

					result := isZero.Call([]reflect.Value{})
					if result[0].Bool() {
						dst.Set(src)
					}
				}
			}
			return nil
		}
	}
	return nil
}

func TestIssue52OverwriteZeroSrcTime(t *testing.T) {
	now := time.Now()
	dst := issue52StructWithTime{now}
	src := issue52StructWithTime{}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.FailNow()
	}

	if !dst.Birth.IsZero() {
		t.Errorf("dst should have been overwritten: dst.Birth(%v) != now(%v)", dst.Birth, now)
	}
}

func TestIssue52OverwriteZeroSrcTimeWithTransformer(t *testing.T) {
	now := time.Now()
	dst := issue52StructWithTime{now}
	src := issue52StructWithTime{}

	if err := merge.Merge(&dst, src, merge.WithTransformers(issue52TimeTransfomer{true}), merge.WithOverwrite()); err != nil {
		t.FailNow()
	}

	if dst.Birth.IsZero() {
		t.Errorf("dst should not have been overwritten: dst.Birth(%v) != now(%v)", dst.Birth, now)
	}
}

func TestIssue52OverwriteZeroDstTime(t *testing.T) {
	now := time.Now()
	dst := issue52StructWithTime{}
	src := issue52StructWithTime{now}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.FailNow()
	}

	if dst.Birth.IsZero() {
		t.Errorf("dst should have been overwritten: dst.Birth(%v) != zero(%v)", dst.Birth, time.Time{})
	}
}

func TestIssue52ZeroDstTime(t *testing.T) {
	now := time.Now()
	dst := issue52StructWithTime{}
	src := issue52StructWithTime{now}

	if err := merge.Merge(&dst, src); err != nil {
		t.FailNow()
	}

	if !dst.Birth.IsZero() {
		t.Errorf("dst should not have been overwritten: dst.Birth(%v) != zero(%v)", dst.Birth, time.Time{})
	}
}

func TestIssue52ZeroDstTimeWithTransformer(t *testing.T) {
	now := time.Now()
	dst := issue52StructWithTime{}
	src := issue52StructWithTime{now}

	if err := merge.Merge(&dst, src, merge.WithTransformers(issue52TimeTransfomer{})); err != nil {
		t.FailNow()
	}

	if dst.Birth.IsZero() {
		t.Errorf("dst should have been overwritten: dst.Birth(%v) != now(%v)", dst.Birth, now)
	}
}

func TestIssue61MergeNilMap(t *testing.T) {
	type T struct {
		I map[string][]string
	}
	t1 := T{}
	t2 := T{I: map[string][]string{"hi": {"there"}}}

	if err := merge.Merge(&t1, t2); err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(t2, T{I: map[string][]string{"hi": {"there"}}}) {
		t.FailNow()
	}
}

type issue64Student struct {
	Name  string
	Books []string
}

type issue64TestData struct {
	S1            issue64Student
	S2            issue64Student
	ExpectedSlice []string
}

func issue64Data() []issue64TestData {
	return []issue64TestData{
		{issue64Student{"Jack", []string{"a", "B"}}, issue64Student{"Tom", []string{"1"}}, []string{"a", "B"}},
		{issue64Student{"Jack", []string{"a", "B"}}, issue64Student{"Tom", []string{}}, []string{"a", "B"}},
		{issue64Student{"Jack", []string{}}, issue64Student{"Tom", []string{"1"}}, []string{"1"}},
		{issue64Student{"Jack", []string{}}, issue64Student{"Tom", []string{}}, []string{}},
	}
}

func TestIssue64MergeSliceWithOverride(t *testing.T) {
	for _, data := range issue64Data() {
		err := merge.Merge(&data.S2, data.S1, merge.WithOverwrite())
		if err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if len(data.S2.Books) != len(data.ExpectedSlice) {
			t.Errorf("Got %d elements in slice, but expected %d", len(data.S2.Books), len(data.ExpectedSlice))
		}

		for i, val := range data.S2.Books {
			if val != data.ExpectedSlice[i] {
				t.Errorf("Expected %s, but got %s while merging slice with override", data.ExpectedSlice[i], val)
			}
		}
	}
}

type issue66PrivateSliceTest struct {
	PublicStrings  []string
	privateStrings []string
}

func TestIssue66PrivateSlice(t *testing.T) {
	p1 := issue66PrivateSliceTest{
		PublicStrings:  []string{"one", "two", "three"},
		privateStrings: []string{"four", "five"},
	}
	p2 := issue66PrivateSliceTest{
		PublicStrings: []string{"six", "seven"},
	}

	if err := merge.Merge(&p1, p2); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if len(p1.PublicStrings) != 3 {
		t.Error("3 elements should be in 'PublicStrings' field, when no append")
	}

	if len(p1.privateStrings) != 2 {
		t.Error("2 elements should be in 'privateStrings' field")
	}
}

func TestIssue66PrivateSliceWithAppendSlice(t *testing.T) {
	p1 := issue66PrivateSliceTest{
		PublicStrings:  []string{"one", "two", "three"},
		privateStrings: []string{"four", "five"},
	}
	p2 := issue66PrivateSliceTest{
		PublicStrings: []string{"six", "seven"},
	}

	if err := merge.Merge(&p1, p2, merge.WithAppendSlice()); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if len(p1.PublicStrings) != 5 {
		t.Error("5 elements should be in 'PublicStrings' field")
	}

	if len(p1.privateStrings) != 2 {
		t.Error("2 elements should be in 'privateStrings' field")
	}
}

type issue83My struct {
	Data []int
}

func TestIssue83(t *testing.T) {
	dst := issue83My{Data: []int{1, 2, 3}}
	new := issue83My{}
	if err := merge.Merge(&dst, new, merge.WithOverwriteWithEmptyValue()); err != nil {
		t.Error(err)
	}
	if len(dst.Data) > 0 {
		t.Errorf("expected empty slice, got %v", dst.Data)
	}
}

type issue84DstStruct struct {
	A int
	B int
	C int
}

type issue84DstNestedStruct struct {
	A struct {
		A int
		B int
		C int
	}
	B int
	C int
}

func TestIssue84MergeMapWithNilValueToStructWithOverride(t *testing.T) {
	p1 := issue84DstStruct{
		A: 0, B: 1, C: 2,
	}
	p2 := map[string]interface{}{
		"A": 3, "B": 4, "C": 0,
	}

	if err := merge.Map(&p1, p2, merge.WithOverwrite()); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if p1.C != 0 {
		t.Error("C field should become '0'")
	}
}

func TestIssue84MergeMapWithoutKeyExistsToStructWithOverride(t *testing.T) {
	p1 := issue84DstStruct{
		A: 0, B: 1, C: 2,
	}
	p2 := map[string]interface{}{
		"A": 3, "B": 4,
	}

	if err := merge.Map(&p1, p2, merge.WithOverwrite()); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if p1.C != 2 {
		t.Error("C field should be '2'")
	}
}

func TestIssue84MergeNestedMapWithNilValueToStructWithOverride(t *testing.T) {
	p1 := issue84DstNestedStruct{
		A: struct {
			A int
			B int
			C int
		}{A: 1, B: 2, C: 0},
		B: 0,
		C: 2,
	}
	p2 := map[string]interface{}{
		"A": map[string]interface{}{
			"A": 0, "B": 0, "C": 5,
		}, "B": 4, "C": 0,
	}

	if err := merge.Map(&p1, p2, merge.WithOverwrite()); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if p1.B != 4 {
		t.Error("A.C field should become '4'")
	}

	if p1.A.C != 5 {
		t.Error("A.C field should become '5'")
	}

	if p1.A.B != 0 || p1.A.A != 0 {
		t.Error("A.A and A.B field should become '0'")
	}
}

func TestIssue89Boolean(t *testing.T) {
	type Foo struct {
		Bar bool `json:"bar"`
	}

	src := Foo{Bar: true}
	dst := Foo{Bar: false}

	if err := merge.Merge(&dst, src); err != nil {
		t.Error(err)
	}
	if dst.Bar == false {
		t.Errorf("expected true, got false")
	}
}

func TestIssue89MergeWithEmptyValue(t *testing.T) {
	p1 := map[string]interface{}{
		"A": 3, "B": "note", "C": true,
	}
	p2 := map[string]interface{}{
		"B": "", "C": false,
	}
	if err := merge.Merge(&p1, p2, merge.WithOverwriteWithEmptyValue()); err != nil {
		t.Error(err)
	}
	testCases := []struct {
		expected interface{}
		key      string
	}{
		{
			3,
			"A",
		},
		{
			"",
			"B",
		},
		{
			false,
			"C",
		},
	}
	for _, tC := range testCases {
		if p1[tC.key] != tC.expected {
			t.Errorf("expected %v in p1[%q], got %v", tC.expected, tC.key, p1[tC.key])
		}
	}
}

type issue90structWithStringMap struct {
	Data map[string]string
}

func TestIssue90(t *testing.T) {
	dst := map[string]issue90structWithStringMap{
		"struct": {
			Data: nil,
		},
	}
	src := map[string]issue90structWithStringMap{
		"struct": {
			Data: map[string]string{
				"foo": "bar",
			},
		},
	}
	expected := map[string]issue90structWithStringMap{
		"struct": {
			Data: map[string]string{
				"foo": "bar",
			},
		},
	}

	err := merge.Merge(&dst, src, merge.WithOverwrite())
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("expected: %#v\ngot: %#v", expected, dst)
	}
}

type issue100s struct {
	Member interface{}
}

func TestIssue100(t *testing.T) {
	m := make(map[string]interface{})
	m["Member"] = "anything"

	st := &issue100s{}
	if err := merge.Map(st, m); err != nil {
		t.Error(err)
	}
}

type issue104Record struct {
	Data    map[string]interface{}
	Mapping map[string]string
}

func issue104StructToRecord(in interface{}) *issue104Record {
	rec := issue104Record{}
	rec.Data = make(map[string]interface{})
	rec.Mapping = make(map[string]string)
	typ := reflect.TypeOf(in)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbFieldName := field.Tag.Get("db")
		if dbFieldName != "" {
			rec.Mapping[field.Name] = dbFieldName
		}
	}

	if err := merge.Map(&rec.Data, in); err != nil {
		panic(err)
	}
	return &rec
}

func TestIssue104StructToRecord(t *testing.T) {
	type A struct {
		Name string `json:"name" db:"name"`
		CIDR string `json:"cidr" db:"cidr"`
	}
	type Record struct {
		Data    map[string]interface{}
		Mapping map[string]string
	}
	a := A{Name: "David", CIDR: "10.0.0.0/8"}
	rec := issue104StructToRecord(a)
	if len(rec.Mapping) < 2 {
		t.Fatalf("struct to record failed, no mapping, struct missing tags?, rec: %+v, a: %+v ", rec, a)
	}
}

func TestIssue121(t *testing.T) {
	dst := map[string]interface{}{
		"inter": map[string]interface{}{
			"a": "1",
			"b": "2",
		},
	}

	src := map[string]interface{}{
		"inter": map[string]interface{}{
			"a": "3",
			"c": "4",
		},
	}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if dst["inter"].(map[string]interface{})["a"].(string) != "3" {
		t.Error("inter.a should equal '3'")
	}

	if dst["inter"].(map[string]interface{})["c"].(string) != "4" {
		t.Error("inter.c should equal '4'")
	}
}

func TestIssue123(t *testing.T) {
	src := map[string]interface{}{
		"col1": nil,
		"col2": 4,
		"col3": nil,
	}
	dst := map[string]interface{}{
		"col1": 2,
		"col2": 3,
		"col3": 3,
	}

	// Expected behavior
	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		expected interface{}
		key      string
	}{
		{
			nil,
			"col1",
		},
		{
			4,
			"col2",
		},
		{
			nil,
			"col3",
		},
	}
	for _, tC := range testCases {
		if dst[tC.key] != tC.expected {
			t.Fatalf("expected %v in dst[%q], got %v", tC.expected, tC.key, dst[tC.key])
		}
	}
}

type issue125Settings struct {
	FirstSlice  []string `json:"FirstSlice"`
	SecondSlice []string `json:"SecondSlice"`
}

func TestIssue125MergeWithOverwrite(t *testing.T) {
	var (
		defaultSettings = issue125Settings{
			FirstSlice:  []string{},
			SecondSlice: []string{},
		}
		something issue125Settings
		data      = `{"FirstSlice":[], "SecondSlice": null}`
	)

	if err := json.Unmarshal([]byte(data), &something); err != nil {
		t.Errorf("Error while Unmarshalling maprequest: %s", err)
	}

	if err := merge.Merge(&something, defaultSettings, merge.WithOverwriteSliceWithEmptyValue()); err != nil {
		t.Errorf("Error while merging: %s", err)
	}

	if something.FirstSlice == nil {
		t.Error("Invalid merging first slice")
	}

	if something.SecondSlice == nil {
		t.Error("Invalid merging second slice")
	}
}

func TestIssue129Boolean(t *testing.T) {
	type Foo struct {
		A bool
		B bool
	}

	src := Foo{
		A: true,
		B: false,
	}
	dst := Foo{
		A: false,
		B: true,
	}

	// Standard behavior
	if err := merge.Merge(&dst, src); err != nil {
		t.Error(err)
	}
	if dst.A != true {
		t.Errorf("expected true, got false")
	}
	if dst.B != true {
		t.Errorf("expected true, got false")
	}

	// Expected behavior
	dst = Foo{
		A: false,
		B: true,
	}
	if err := merge.Merge(&dst, src, merge.WithOverwriteWithEmptyValue()); err != nil {
		t.Error(err)
	}
	if dst.A != true {
		t.Errorf("expected true, got false")
	}
	if dst.B != false {
		t.Errorf("expected false, got true")
	}
}

type issue131Foz struct {
	A *bool
	B string
}

func TestIssue131MergeWithOverwriteWithEmptyValue(t *testing.T) {
	src := issue131Foz{
		A: func(v bool) *bool { return &v }(false),
		B: "src",
	}
	dest := issue131Foz{
		A: func(v bool) *bool { return &v }(true),
		B: "dest",
	}
	if err := merge.Merge(&dest, src, merge.WithOverwriteWithEmptyValue()); err != nil {
		t.Error(err)
	}
	if *src.A != *dest.A {
		t.Errorf("dest.A not merged in properly: %v != %v", *src.A, *dest.A)
	}
	if src.B != dest.B {
		t.Errorf("dest.B not merged in properly: %v != %v", src.B, dest.B)
	}
}

type issue136EmbeddedTestA struct {
	Name string
	Age  uint8
}

type issue136EmbeddedTestB struct {
	Address string
	issue136EmbeddedTestA
}

func TestIssue136MergeEmbedded(t *testing.T) {
	var (
		err error
		a   = &issue136EmbeddedTestA{
			"Suwon", 16,
		}
		b = &issue136EmbeddedTestB{}
	)

	if err := merge.Merge(&b.issue136EmbeddedTestA, *a); err != nil {
		t.Error(err)
	}

	if b.Name != "Suwon" {
		t.Errorf("%v %v", b.Name, err)
	}
}

const issue138configuration string = `
{
	"Port": 80
}
`

func TestIssue138(t *testing.T) {
	type config struct {
		Port uint16
	}
	type compatibleConfig struct {
		Port float64
	}

	foo := make(map[string]interface{})
	// encoding/json unmarshals numbers as float64
	// https://golang.org/pkg/encoding/json/#Unmarshal
	json.Unmarshal([]byte(issue138configuration), &foo)

	err := merge.Map(&config{}, foo)
	if err == nil {
		t.Error("expected type mismatch error, got nil")
	} else {
		if err.Error() != "type mismatch on Port field: found float64, expected uint16" {
			t.Errorf("expected type mismatch error, got %q", err)
		}
	}

	c := compatibleConfig{}
	if err := merge.Map(&c, foo); err != nil {
		t.Error(err)
	}
}

func TestIssue143(t *testing.T) {
	testCases := []struct {
		expected func(map[string]interface{}) error
		options  []merge.Option
	}{
		{
			options: []merge.Option{merge.WithOverwrite()},
			expected: func(m map[string]interface{}) error {
				properties := m["properties"].(map[string]interface{})
				if properties["field1"] != "wrong" {
					return fmt.Errorf("expected %q, got %v", "wrong", properties["field1"])
				}
				return nil
			},
		},
		{
			options: []merge.Option{},
			expected: func(m map[string]interface{}) error {
				properties := m["properties"].(map[string]interface{})
				if properties["field1"] == "wrong" {
					return fmt.Errorf("expected a map, got %v", "wrong")
				}
				return nil
			},
		},
	}
	for _, tC := range testCases {
		base := map[string]interface{}{
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "text",
				},
			},
		}

		err := merge.Map(
			&base,
			map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": "wrong",
				},
			},
			tC.options...,
		)
		if err != nil {
			t.Error(err)
		}
		if err := tC.expected(base); err != nil {
			t.Error(err)
		}
	}
}

type issue149User struct {
	Name string
}

type issue149Token struct {
	User  *issue149User
	Token *string
}

func TestIssue149(t *testing.T) {
	dest := &issue149Token{
		User: &issue149User{
			Name: "destination",
		},
		Token: nil,
	}
	tokenValue := "Issue149"
	src := &issue149Token{
		User:  nil,
		Token: &tokenValue,
	}
	if err := merge.Merge(dest, src, merge.WithOverwriteWithEmptyValue()); err != nil {
		t.Error(err)
	}
	if dest.User != nil {
		t.Errorf("expected nil User, got %q", dest.User)
	}
	if dest.Token == nil {
		t.Errorf("expected not nil Token, got %q", *dest.Token)
	}
}

type issue174StructWithBlankField struct {
	_ struct{}
	A struct{}
}

func TestIssue174(t *testing.T) {
	dst := issue174StructWithBlankField{}
	src := issue174StructWithBlankField{}

	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Error(err)
	}
}

func TestIssue209(t *testing.T) {
	dst := []string{"a", "b"}
	src := []string{"c", "d"}

	if err := merge.Merge(&dst, src, merge.WithAppendSlice()); err != nil {
		t.Error(err)
	}

	expected := []string{"a", "b", "c", "d"}
	if len(dst) != len(expected) {
		t.Errorf("arrays not equal length")
	}
	for i := range expected {
		if dst[i] != expected[i] {
			t.Errorf("array elements at %d are not equal", i)
		}
	}
}

var issueTemplateTestDataS = []struct {
	S1            issue64Student
	S2            issue64Student
	ExpectedSlice []string
}{
	{issue64Student{"Jack", []string{"a", "B"}}, issue64Student{"Tom", []string{"1"}}, []string{"1", "a", "B"}},
	{issue64Student{"Jack", []string{"a", "B"}}, issue64Student{"Tom", []string{}}, []string{"a", "B"}},
	{issue64Student{"Jack", []string{}}, issue64Student{"Tom", []string{"1"}}, []string{"1"}},
	{issue64Student{"Jack", []string{}}, issue64Student{"Tom", []string{}}, []string{}},
}

func TestIssueTemplateMergeSliceWithOverrideWithAppendSlice(t *testing.T) {
	for _, data := range issueTemplateTestDataS {
		err := merge.Merge(&data.S2, data.S1, merge.WithOverwrite(), merge.WithAppendSlice())
		if err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if len(data.S2.Books) != len(data.ExpectedSlice) {
			t.Errorf("Got %d elements in slice, but expected %d", len(data.S2.Books), len(data.ExpectedSlice))
		}

		for i, val := range data.S2.Books {
			if val != data.ExpectedSlice[i] {
				t.Errorf("Expected %s, but got %s while merging slice with override", data.ExpectedSlice[i], val)
			}
		}
	}
}

type pr80MapInterface map[string]interface{}

func TestPr80MergeMapsEmptyString(t *testing.T) {
	a := pr80MapInterface{"s": ""}
	b := pr80MapInterface{"s": "foo"}
	if err := merge.Merge(&a, b); err != nil {
		t.Error(err)
	}
	if a["s"] != "foo" {
		t.Errorf("b not merged in properly: a.s.Value(%s) != expected(%s)", a["s"], "foo")
	}
}

func TestPr81MapInterfaceWithMultipleLayer(t *testing.T) {
	m1 := map[string]interface{}{
		"k1": map[string]interface{}{
			"k1.1": "v1",
		},
	}

	m2 := map[string]interface{}{
		"k1": map[string]interface{}{
			"k1.1": "v2",
			"k1.2": "v3",
		},
	}

	if err := merge.Map(&m1, m2, merge.WithOverwrite()); err != nil {
		t.Errorf("Error merging: %v", err)
	}

	// Check overwrite of sub map works
	expected := "v2"
	actual := m1["k1"].(map[string]interface{})["k1.1"].(string)
	if actual != expected {
		t.Errorf("Expected %v but got %v",
			expected,
			actual)
	}

	// Check new key is merged
	expected = "v3"
	actual = m1["k1"].(map[string]interface{})["k1.2"].(string)
	if actual != expected {
		t.Errorf("Expected %v but got %v",
			expected,
			actual)
	}
}

func TestPr211MergeWithTransformerZeroValue(t *testing.T) {
	// This test specifically tests that a transformer can be used to
	// prevent overwriting a zero value (in this case a bool). This would fail prior to #211
	type fooWithBoolPtr struct {
		b *bool
	}
	var Bool = func(b bool) *bool { return &b }
	a := fooWithBoolPtr{b: Bool(false)}
	b := fooWithBoolPtr{b: Bool(true)}

	if err := merge.Merge(&a, &b, merge.WithTransformers(&transformer{
		m: map[reflect.Type]func(dst, src reflect.Value) error{
			reflect.TypeOf(Bool(false)): func(dst, src reflect.Value) error {
				if dst.CanSet() && dst.IsNil() {
					dst.Set(src)
				}
				return nil
			},
		},
	})); err != nil {
		t.Error(err)
	}

	if *a.b != false {
		t.Errorf("b not merged in properly: a.b(%v) != expected(%v)", a.b, false)
	}
}

type v039Inner struct {
	A int
}

type v039Outer struct {
	v039Inner
	B int
}

func TestV039Issue139(t *testing.T) {
	dst := v039Outer{
		v039Inner: v039Inner{A: 1},
		B:         2,
	}
	src := v039Outer{
		v039Inner: v039Inner{A: 10},
		B:         20,
	}
	err := merge.Merge(&dst, src, merge.WithOverwrite())
	if err != nil {
		panic(err.Error())
	}
	if dst.v039Inner.A == 1 {
		t.Errorf("expected %d, got %d", src.v039Inner.A, dst.v039Inner.A)
	}
}

func TestV039Issue152(t *testing.T) {
	dst := map[string]interface{}{
		"properties": map[string]interface{}{
			"field1": map[string]interface{}{
				"type": "text",
			},
			"field2": "ohai",
		},
	}
	src := map[string]interface{}{
		"properties": map[string]interface{}{
			"field1": "wrong",
		},
	}
	if err := merge.Map(&dst, src, merge.WithOverwrite()); err != nil {
		t.Error(err)
	}
}

type v039Foo struct {
	B map[string]v039Bar
	A string
}

type v039Bar struct {
	C *string
	D *string
}

func TestV039Issue146(t *testing.T) {
	var (
		s1 = "asd"
		s2 = "sdf"
	)
	dst := v039Foo{
		A: "two",
		B: map[string]v039Bar{
			"foo": {
				C: &s1,
			},
		},
	}
	src := v039Foo{
		A: "one",
		B: map[string]v039Bar{
			"foo": {
				D: &s2,
			},
		},
	}
	if err := merge.Merge(&dst, src, merge.WithOverwrite()); err != nil {
		t.Error(err)
	}
	if dst.B["foo"].D == nil {
		t.Errorf("expected %v, got nil", &s2)
	}
}
