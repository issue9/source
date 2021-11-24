// SPDX-License-Identifier: MIT

package source

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert"
)

func TestDumpGoFile(t *testing.T) {
	a := assert.New(t)

	a.NotError(DumpGoSource("./testdata/go.go", []byte("var x=1")))
	content, err := ioutil.ReadFile("./testdata/go.go")
	a.NotError(err)
	a.Equal(string(content), "var x = 1")
}

func TestCurrentPath(t *testing.T) {
	a := assert.New(t)

	dir, err := filepath.Abs("./file.go")
	a.NotError(err).NotEmpty(dir)

	d, err := filepath.Abs(CurrentPath("./file.go"))
	a.NotError(err).NotEmpty(d)

	a.Equal(d, dir)
}

func TestCurrentDir(t *testing.T) {
	a := assert.New(t)

	dir, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(dir)

	a.Equal(CurrentDir(), dir)
}

func TestCurrentFile(t *testing.T) {
	a := assert.New(t)

	filename, err := filepath.Abs("./source_test.go")
	a.NotError(err).NotEmpty(filename)

	a.Equal(CurrentFile(), filepath.ToSlash(filename))
}

func TestCurrentFunction(t *testing.T) {
	a := assert.New(t)

	a.Equal(CurrentFunction(), "TestCurrentFunction")
}

func TestCurrentLine(t *testing.T) {
	a := assert.New(t)

	a.Equal(CurrentLine(), 62)
}

func TestCurrentLocation(t *testing.T) {
	a:=assert.New(t)

	path,line :=CurrentLocation()
	a.True(strings.HasSuffix(path,"source_test.go")).
		Equal(line,68)
}

func TestTraceStack(t *testing.T) {
	a := assert.New(t)

	str, err := TraceStack(1, "message", 12)
	a.NotError(err)
	a.True(strings.HasPrefix(str, "message 12"))
	a.True(strings.Contains(str, "source_test.go")) // 肯定包含当前文件名
}
