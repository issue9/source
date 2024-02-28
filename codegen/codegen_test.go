// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package codegen

import (
	"io/fs"
	"os"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestDump(t *testing.T) {
	a := assert.New(t, false)
	const file = "./go.g"

	a.NotError(Dump(file, []byte("var x=1"), fs.ModePerm))
	content, err := os.ReadFile(file)
	a.NotError(err)
	a.Equal(string(content), "var x = 1")

	os.Remove(file)
}

func TestDumpFromTemplate(t *testing.T) {
	a := assert.New(t, false)
	const file = "./tpl1.g"
	const tpl = "var x={{.X}}"

	a.NotError(DumpFromTemplate(file, tpl, struct{ X int }{X: 1}, fs.ModePerm))
	content, err := os.ReadFile(file)
	a.NotError(err)
	a.Equal(string(content), "var x = 1")

	os.Remove(file)
}

func TestDumpFromTemplateFile(t *testing.T) {
	a := assert.New(t, false)
	const file = "./tpl2.g"

	a.NotError(DumpFromTemplateFile(file, "./tpl.tpl", struct {
		X int
		H string
	}{X: 1, H: "H"}, fs.ModePerm))
	content, err := os.ReadFile(file)
	a.NotError(err)
	a.Equal(string(content), "// H\n\npackage x\n\nvar x = 1\n")

	os.Remove(file)
}
