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
//
// 对于结构类型会自动展开。
//
// t 需要转换的类型；
// m 需要特殊定义的类型；
// unexported 是否导出小写字段；
func GoDefine(t reflect.Type, m map[reflect.Type]string, unexported bool) string {
	buf := &errwrap.Buffer{}
	goDefine(buf, 0, t, m, unexported, false)
	return buf.String()
}

func goDefine(buf *errwrap.Buffer, indent int, t reflect.Type, m map[reflect.Type]string, unexported, anonymous bool) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if len(m) > 0 {
		if s, found := m[t]; found {
			buf.WString(s)
			return
		}
	}

	switch t.Kind() {
	case reflect.Func, reflect.Chan: // 忽略
	case reflect.Slice:
		buf.WString("[]").WString(t.Elem().Name())
	case reflect.Array:
		buf.WByte('[').WString(strconv.Itoa(t.Len())).WByte(']').WString(t.Elem().Name())
	case reflect.Struct:
		if !anonymous {
			buf.WString("{\n")
			indent++
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if f.Anonymous {
				goDefine(buf, indent, f.Type, m, unexported, true)
				continue
			}

			if !unexported && !f.IsExported() {
				continue
			}

			if f.Type.Kind() == reflect.Func || f.Type.Kind() == reflect.Chan {
				continue
			}

			buf.WString(strings.Repeat("\t", indent)).WString(f.Name).WByte('\t')
			goDefine(buf, indent, f.Type, m, unexported, false)

			if f.Tag != "" {
				buf.WByte('\t').WByte('`').WString(string(f.Tag)).WByte('`')
			}

			buf.WByte('\n')
		}

		if !anonymous {
			indent--
			buf.WString(strings.Repeat("\t", indent)).WByte('}')
		}
	default:
		buf.WString(t.Name())
	}
}
