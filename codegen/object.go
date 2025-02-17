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
// m 需要特殊定义的类型，可以为空；
// c 为对象的字段生成注释，可以为空；
// unexported 是否导出小写字段；
func GoDefine(t reflect.Type, m map[reflect.Type]string, c func(*reflect.StructField) string, unexported bool) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	buf := &errwrap.Buffer{}
	goDefine(buf, 0, t, m, c, unexported, false)
	s := buf.String()

	if strings.HasPrefix(s, "struct {") { // 结构可能由于 m 的关系返回一个非结构体的类型定义，所以只能由开头是否为 struct { 判断是否为结构体。
		return "type " + t.Name() + " " + s
	}
	return buf.String()
}

func goDefine(buf *errwrap.Buffer, indent int, t reflect.Type, m map[reflect.Type]string, c func(*reflect.StructField) string, unexported, anonymous bool) {
	if len(m) > 0 {
		if s, found := m[t]; found {
			buf.WString(s)
			return
		}
	}

	switch t.Kind() {
	case reflect.Func, reflect.Chan: // 忽略
	case reflect.Pointer:
		buf.WByte('*')
		goDefine(buf, indent, t.Elem(), m, c, unexported, anonymous)
	case reflect.Slice:
		buf.WString("[]")
		goDefine(buf, indent, t.Elem(), m, c, unexported, anonymous)
	case reflect.Array:
		buf.WByte('[').WString(strconv.Itoa(t.Len())).WByte(']')
		goDefine(buf, indent, t.Elem(), m, c, unexported, anonymous)
	case reflect.Struct:
		if !anonymous {
			if t.NumField() == 0 {
				buf.WString("struct {}")
				return
			}

			buf.WString("struct {\n")
			indent++
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if f.Anonymous {
				tt := f.Type
				for tt.Kind() == reflect.Pointer { // 匿名字段需要去掉指针类型
					tt = tt.Elem()
				}
				goDefine(buf, indent, tt, m, c, unexported, true)
				continue
			}

			if !unexported && !f.IsExported() {
				continue
			}

			if f.Type.Kind() == reflect.Func || f.Type.Kind() == reflect.Chan {
				continue
			}

			buf.WString(strings.Repeat("\t", indent)).WString(f.Name).WByte('\t')
			goDefine(buf, indent, f.Type, m, c, unexported, false)

			if f.Tag != "" {
				buf.WByte('\t').WByte('`').WString(string(f.Tag)).WByte('`')
			}

			if c != nil {
				if comment := c(&f); comment != "" {
					buf.WString("\t// ").WString(comment)
				}
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
