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

func TestGoDefine(t *testing.T) {
	a := assert.New(t, false)

	wont := "{\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tArray\t[5]int\n\tSlice\t[]string\n\tByte\tuint8\n}"
	a.Equal(GoDefine(reflect.TypeFor[object1](), nil, false), wont)

	wont = "{\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tObject\t{\n\t\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\t\tArray\t[5]int\n\t\tSlice\t[]string\n\t\tByte\tuint8\n\t}\t`json:\"object\"`\n}"
	a.Equal(GoDefine(reflect.TypeFor[*object2](), nil, false), wont)

	a.Equal(GoDefine(reflect.TypeFor[int](), nil, false), "int")

	a.Equal(GoDefine(reflect.TypeFor[string](), nil, false), "string")

	a.Equal(GoDefine(reflect.TypeFor[func()](), nil, false), "")

	m := map[reflect.Type]string{reflect.TypeFor[time.Time](): "string"}

	a.Equal(GoDefine(reflect.TypeFor[time.Time](), m, false), "string")
	a.Equal(GoDefine(reflect.TypeFor[*time.Time](), m, false), "string")

	a.Equal(GoDefine(reflect.TypeFor[*object3](), m, false), "{\n\tT\tstring\t`json:\"t\"`\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object3](), m, true), "{\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n}")

	a.Equal(GoDefine(reflect.TypeFor[*object4](), m, true), "{\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object4](), m, false), "{\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n}")

	a.Equal(GoDefine(reflect.TypeFor[*object5](), m, true), "{\n\tint\tint\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n\tStr\tstring\n}")
	a.Equal(GoDefine(reflect.TypeFor[*object5](), m, false), "{\n\tT\tstring\t`json:\"t\"`\n\tInt\tint\n\tStr\tstring\n}")
}
