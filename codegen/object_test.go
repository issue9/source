// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package codegen

import (
	"reflect"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
)

type object1 struct {
	Int   int `json:"int" yaml:"int"`
	Array [5]int
	Slice []string
	Byte  byte
	Chan  chan int
}

type object2 struct {
	Int    int      `json:"int" yaml:"int"`
	Object *object1 `json:"object"`
}

type object3 struct {
	int int
	T   time.Time `json:"t"`
}

type object4 struct {
	*object3
	Int int
}

type object5 struct {
	*object4
	Str string
}

type object6 struct {
	XMLName struct{} `json:"root"`
	Str     []*object3
}

func TestGoDefine(t *testing.T) {
	a := assert.New(t, false)

	c := func(f *reflect.StructField) string {
		return f.Name
	}

	wont := "type object1 struct {\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tArray\t[5]int\n\tSlice\t[]string\n\tByte\tuint8\n}"
	a.Equal(GoDefine(reflect.TypeFor[object1](), nil, nil, false), wont)

	wont = "type object2 struct {\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tObject\t*struct {\n\t\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\t\tArray\t[5]int\n\t\tSlice\t[]string\n\t\tByte\tuint8\n\t}\t`json:\"object\"`\n}"
	a.Equal(GoDefine(reflect.TypeFor[*object2](), nil, nil, false), wont)

	a.Equal(GoDefine(reflect.TypeFor[int](), nil, nil, false), "int")

	a.Equal(GoDefine(reflect.TypeFor[string](), nil, nil, false), "string")

	a.Equal(GoDefine(reflect.TypeFor[func()](), nil, nil, false), "")

	m := map[reflect.Type]string{reflect.TypeFor[time.Time](): "string"}

	a.Equal(GoDefine(reflect.TypeFor[time.Time](), m, nil, false), "string")
	a.Equal(GoDefine(reflect.TypeFor[*time.Time](), m, nil, false), "string")

	a.Equal(GoDefine(reflect.TypeFor[*object3](), m, nil, false), "type object3 struct {\n\tT\tstring\t`json:\"t\"`\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object3](), m, nil, true), "type object3 struct {\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n}")

	a.Equal(GoDefine(reflect.TypeFor[*object4](), m, nil, true), "type object4 struct {\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object4](), m, nil, false), "type object4 struct {\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n}")

	a.Equal(GoDefine(reflect.TypeFor[*object5](), m, nil, true), "type object5 struct {\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n\tStr\tstring\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object5](), m, c, false), "type object5 struct {\n\tT\tstring\t`json:\"t\"`\t// T\n\tInt\tint\t// Int\n\tStr\tstring\t// Str\n}")

	a.Equal(GoDefine(reflect.TypeFor[time.Time](), m, nil, false), "string")
	a.Equal(GoDefine(reflect.TypeFor[time.Time](), nil, nil, false), "type Time struct {\n}")

	a.Equal(GoDefine(reflect.TypeFor[object6](), m, nil, false), "type object6 struct {\n\tXMLName\tstruct {}\t`json:\"root\"`\n\tStr\t[]*struct {\n\t\tT\tstring\t`json:\"t\"`\n\t}\n}")
}
