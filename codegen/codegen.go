// SPDX-License-Identifier: MIT

// Package codegen 简单的代码生成工具
package codegen

import (
	"bytes"
	"go/format"
	"io/fs"
	"os"
	"text/template"
)

// Dump 格式化并输出 Go 源码到 path
func Dump(path string, content []byte, mode fs.FileMode) error {
	src, err := format.Source(content)
	if err != nil {
		return err
	}
	return os.WriteFile(path, src, mode)
}

// DumpFromTemplate 根据模板内容生成代码
func DumpFromTemplate(path, tpl string, data any, mode fs.FileMode) error {
	t, err := template.New("codegen").Parse(tpl)
	if err != nil {
		return nil
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return err
	}
	return Dump(path, buf.Bytes(), mode)
}

// DumpFromTemplateFile 根据 tplFile 指向的模板内容生成代码
func DumpFromTemplateFile(path string, tplFile string, data any, mode fs.FileMode) error {
	b, err := os.ReadFile(tplFile)
	if err != nil {
		return err
	}

	return DumpFromTemplate(path, string(b), data, mode)
}
