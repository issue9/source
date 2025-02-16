// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package codegen

import (
	"reflect"
	"testing"

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

func TestGoDefine(t *testing.T) {
	a := assert.New(t, false)

	wont := "{\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tArray\t[5]int\n\tSlice\t[]string\n\tByte\tuint8\n}"
	a.Equal(GoDefine(reflect.TypeFor[object1]()), wont)

	wont = "{\n\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\tObject\t{\n\t\tInt\tint\t`json:\"int\" yaml:\"int\"`\n\t\tArray\t[5]int\n\t\tSlice\t[]string\n\t\tByte\tuint8\n\t}\t`json:\"object\"`\n}"
	a.Equal(GoDefine(reflect.TypeFor[*object2]()), wont)

	a.Equal(GoDefine(reflect.TypeFor[int]()), "int")

	a.Equal(GoDefine(reflect.TypeFor[string]()), "string")

	a.Equal(GoDefine(reflect.TypeFor[func()]()), "")
}
