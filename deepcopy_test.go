package deepcopy_test

import (
	"bytes"
	"log"
	"reflect"
	"testing"

	"github.com/cloverstd/deepcopy"
)

func TestCCC(t *testing.T) {
	var i = 10
	var a = struct {
		A int
		B *int
		a int
	}{
		A: 10,
		B: &i,
	}
	b, err := deepcopy.Copy(a)
	log.Printf("%#v, %#v, %#v", b, err, a)
}

func TestBasicType(t *testing.T) {

	for _, val := range []struct {
		a   interface{}
		msg string
	}{
		{
			a: 1,
		},
		{
			a: 1.1,
		},
		{
			a: true,
		},
		{
			a: "test",
		},
		{
			a: uint(10),
		},
	} {
		typ := reflect.TypeOf(val.a)
		b, err := deepcopy.Copy(val.a)
		if err != nil {
			t.Errorf("test type %s failed", typ.Kind())
		}
		if reflect.TypeOf(b) != typ {
			t.Errorf("test type %s type failed", typ.Kind())
		}
		if b != val.a {
			t.Errorf("test type %s value failed, %v, %v", typ.Kind(), val.a, b)
		}
	}
}

func TestSlice(t *testing.T) {
	a1, a2, a3, a4, a5 := 1, 2, 3, 4, 5
	for _, val := range []struct {
		a   interface{}
		msg string
	}{
		{
			a: []int{1, 2, 3, 4, 5},
		},
		{
			a: []int{1, 2, 3},
		},
		{
			a: []*int{&a1, &a2, &a3, &a4, &a5},
		},
		{
			a: []interface{}{
				&a1,
				10,
			},
		},
	} {
		typ := reflect.TypeOf(val.a)
		b, err := deepcopy.Copy(val.a)
		if err != nil {
			t.Errorf("test type %s failed, %s", typ.Kind(), err)
		}
		if reflect.TypeOf(b) != typ {
			t.Errorf("test type %s type failed, %s", typ.Kind(), reflect.TypeOf(b))
		}
		for i := 0; i < reflect.ValueOf(val.a).Len(); i++ {
			if reflect.ValueOf(val.a).Index(i).Kind() != reflect.ValueOf(b).Index(i).Kind() {
				t.Errorf("test type %s type failed", typ.Kind())
			}
			if !reflect.DeepEqual(reflect.ValueOf(val.a).Index(i).Interface(), reflect.ValueOf(b).Index(i).Interface()) {
				t.Errorf("test type %s type failed", typ.Kind())
			}
		}
	}
}

func TestArray(t *testing.T) {
	a1, a2, a3, a4, a5 := 1, 2, 3, 4, 5
	for _, val := range []struct {
		a   interface{}
		msg string
	}{
		{
			a: [5]int{1, 2, 3, 4, 5},
		},
		{
			a: [5]int{1, 2, 3},
		},
		{
			a: [5]*int{&a1, &a2, &a3, &a4, &a5},
		},
		{
			a: [5]interface{}{
				&a1,
				10,
			},
		},
	} {
		typ := reflect.TypeOf(val.a)
		b, err := deepcopy.Copy(val.a)
		if err != nil {
			t.Errorf("test type %s failed", typ.Kind())
		}
		if reflect.TypeOf(b) != typ {
			t.Errorf("test type %s type failed", typ.Kind())
		}
		for i := 0; i < reflect.ValueOf(val.a).Len(); i++ {
			if reflect.ValueOf(val.a).Index(i).Kind() != reflect.ValueOf(b).Index(i).Kind() {
				t.Errorf("test type %s type failed", typ.Kind())
			}
			if !reflect.DeepEqual(reflect.ValueOf(val.a).Index(i).Interface(), reflect.ValueOf(b).Index(i).Interface()) {
				t.Errorf("test type %s type failed", typ.Kind())
			}
		}
	}
}

func TestTest(t *testing.T) {
	var i = 10
	b := &i
	typ := reflect.TypeOf(b)
	val := reflect.ValueOf(b)

	log.Println(typ.Kind(), val.Type().Kind(), val.Elem().Kind())
	newVal := reflect.New(typ.Elem())
	log.Println(newVal.Elem().Kind(), newVal.Elem().CanSet())
	newVal.Elem().Set(val.Elem())
	log.Println(*newVal.Interface().(*int))

}

func TestPtr(t *testing.T) {
	var i = 10
	var ii *int = nil
	a := &ii

	b, err := deepcopy.Copy(a)
	if err != nil {
		t.Error("test nil ptr failed, ", err)
	}
	if ba, ok := b.(**int); !ok {
		t.Errorf("test nil ptr type failed, %s", reflect.TypeOf(b))
	} else if *ba != nil {
		t.Errorf("test nil ptr value should be nil, %v", *ba)
	}
	var c = 10
	ii = &c

	b, err = deepcopy.Copy(a)
	if err != nil {
		t.Error("test ptr failed, ", err)
	}
	if ba, ok := b.(**int); !ok {
		t.Errorf("test ptr type failed, %s", reflect.TypeOf(b))
	} else if *ba == nil {
		t.Errorf("test ptr value should not be nil")
	} else if **ba != i {
		t.Errorf("test ptr value failed, %d", **ba)
	} else {
		if *a != &c {
			t.Error("change value failed")
		}
		if &c == *ba {
			t.Error("change value failed")
		}
	}

}

func TestMap(t *testing.T) {
	a := map[string]int{
		"1": 1,
	}
	b, err := deepcopy.Copy(a)
	if err != nil {
		t.Error("test map copy failed, ", err)
	}
	if !reflect.DeepEqual(b, a) {
		t.Error("test map copy failed")
	}
	a["1"] = 2
	a["2"] = 3
	if bb, ok := b.(map[string]int); !ok || bb["1"] == 2 {
		t.Error("test map copy value failed")
	} else if _, exist := bb["2"]; exist {
		t.Error("test map copy value failed")
	}

	aa := map[string]map[string]map[string]int{
		"1": map[string]map[string]int{
			"1": map[string]int{
				"1": 1,
			},
		},
	}
	b, err = deepcopy.Copy(aa)
	if err != nil {
		t.Error("test map deep copy failed, ", err)
	}
	if !reflect.DeepEqual(b, aa) {
		t.Error("test map deep copy failed")
	}
	aa["2"] = map[string]map[string]int{
		"1": map[string]int{
			"1": 1,
		},
	}
	aa["1"]["2"] = map[string]int{
		"1": 1,
	}
	aa["1"]["1"]["1"] = 10
	if bb, ok := b.(map[string]map[string]map[string]int); !ok {
		t.Error("test map deep copy value failed")
	} else if _, exist := bb["2"]; exist {
		t.Error("test map deep copy change failed")
	} else if bb["1"]["1"]["1"] == 10 {
		t.Error("test map deep copy change failed")
	}

	var aaa map[string]int
	b, err = deepcopy.Copy(aaa)
	if err != nil {
		t.Error("test map nil copy failed,", err)
	}
	if bb, ok := b.(map[string]int); !ok {
		t.Error("test map nil copy failed", reflect.TypeOf(b))
	} else if len(bb) != 0 {
		t.Error("test map nil copy failed", bb)
	}

}

func TestHybird(t *testing.T) {
	var i = 10
	var f = 10.0
	var s = "1"
	var b = []byte("10")
	a := map[string]interface{}{
		"slice": []int{1, 2, 3},
		"array": [3]int{1, 2, 3},
		"map": map[string]int{
			"1": 10,
		},
		"struct": struct {
			A   int
			Map map[string]int
		}{
			A: i,
			Map: map[string]int{
				"1": 10,
			},
		},
		"ptr": map[string]interface{}{
			"1": &i,
			"slice": []interface{}{
				&f, &s,
			},
			"byte": b,
		},
	}
	to, err := deepcopy.Copy(a)
	if err != nil {
		t.Error("test map deep copy failed, ", err)
	}

	if !reflect.DeepEqual(to, a) {
		t.Error("test map deep copy failed")
	}

	if !bytes.Equal(to.(map[string]interface{})["ptr"].(map[string]interface{})["byte"].([]byte), b) {
		t.Error("test byte copy failed")
	}

	if *to.(map[string]interface{})["ptr"].(map[string]interface{})["1"].(*int) != i {
		t.Error("test byte copy failed")
	}

}

func TestCycle(t *testing.T) {
	a := map[string]interface{}{
		"1": 10,
	}
	a["a"] = a
	_, err := deepcopy.Copy(a)
	if err != deepcopy.ErrMapDepth {
		t.Error("it should be ErrMapDepth")
	}

	type Cycle struct {
		Cycle *Cycle
		Int   int
	}
	aa := Cycle{
		Int: 10,
	}
	aa.Cycle = &aa
	b, err := deepcopy.Copy(aa)
	if err != deepcopy.ErrMapDepth {
		t.Error("it should be ErrMapDepth", err)
	}
	log.Println("aaa", aa, b)
}

func TestError(t *testing.T) {
	a := map[string]interface{}{
		"chan": make(chan struct{}),
	}
	_, err := deepcopy.Copy(a)
	if err == nil {
		t.Error("copy chan should not support")
	}

}
