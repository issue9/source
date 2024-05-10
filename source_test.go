// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package source

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestCurrentPath(t *testing.T) {
	a := assert.New(t, false)

	dir, err := filepath.Abs("./file.go")
	a.NotError(err).NotEmpty(dir)

	d, err := filepath.Abs(CurrentPath("./file.go"))
	a.NotError(err).NotEmpty(d)

	a.Equal(d, dir)
}

func TestCurrentDir(t *testing.T) {
	a := assert.New(t, false)

	dir, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(dir)

	a.Equal(CurrentDir(), dir)
}

func TestCurrentFile(t *testing.T) {
	a := assert.New(t, false)

	filename, err := filepath.Abs("./source_test.go")
	a.NotError(err).NotEmpty(filename)

	a.Equal(CurrentFile(), filepath.ToSlash(filename))
}

func TestCurrentFunction(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(CurrentFunction(), "TestCurrentFunction")
}

func TestCurrentLine(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(CurrentLine(), 54)
}

func TestCurrentLocation(t *testing.T) {
	a := assert.New(t, false)

	path, line := CurrentLocation()
	a.True(strings.HasSuffix(path, "source_test.go")).
		Equal(line, 60)
}

func TestStack(t *testing.T) {
	a := assert.New(t, false)

	str := Stack(1, true, "message", 12)
	t.Log(str)
	a.True(strings.HasPrefix(str, "message 12"))
	a.True(strings.Contains(str, "source_test.go:68"), str) // 依赖调用 Stack 的行号
}
