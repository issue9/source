// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package codegen

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/issue9/errwrap"
)

// GoDefine 返回对象 v 在 Go 源码中的定义方式
func GoDefine(t reflect.Type) string {
	buf := &errwrap.Buffer{}
	goDefine(buf, 0, t)
	return buf.String()
}

func goDefine(buf *errwrap.Buffer, indent int, t reflect.Type) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Func, reflect.Chan: // 忽略
	case reflect.Slice:
		buf.WString("[]").WString(t.Elem().Name())
	case reflect.Array:
		buf.WByte('[').WString(strconv.Itoa(t.Len())).WByte(']').WString(t.Elem().Name())
	case reflect.Struct:
		buf.WString("{\n")
		indent++
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Type.Kind() == reflect.Func || f.Type.Kind() == reflect.Chan {
				continue
			}

			buf.WString(strings.Repeat("\t", indent)).WString(f.Name).WByte('\t')
			goDefine(buf, indent, f.Type)

			if f.Tag != "" {
				buf.WByte('\t').WByte('`').WString(string(f.Tag)).WByte('`')
			}

			buf.WByte('\n')
		}

		indent--
		buf.WString(strings.Repeat("\t", indent)).WByte('}')
	default:
		buf.WString(t.Name())
	}
}
