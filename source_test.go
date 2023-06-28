// SPDX-License-Identifier: MIT

package source

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestDumpGoFile(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(DumpGoSource("./testdata/go.go", []byte("var x=1")))
	content, err := os.ReadFile("./testdata/go.go")
	a.NotError(err)
	a.Equal(string(content), "var x = 1")

	os.Remove("./testdata/go.go")
}

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

	a.Equal(CurrentLine(), 64)
}

func TestCurrentLocation(t *testing.T) {
	a := assert.New(t, false)

	path, line := CurrentLocation()
	a.True(strings.HasSuffix(path, "source_test.go")).
		Equal(line, 70)
}

func TestStack(t *testing.T) {
	a := assert.New(t, false)

	str := Stack(1, "message", 12)
	t.Log(str)
	a.True(strings.HasPrefix(str, "message 12"))
	a.True(strings.Contains(str, "source_test.go:78"), str) // 依赖调用 Stack 的行号
}
