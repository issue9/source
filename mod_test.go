// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package source

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestPkgSourceDir(t *testing.T) {
	a := assert.New(t, false)

	// std
	dir, err := PkgSourceDir("encoding/json", "./", false)
	a.NotError(err).FileExists(dir)

	dir, err = PkgSourceDir("note-exists", "./", false)
	a.NotError(err).FileNotExists(dir)

	// require

	dir, err = PkgSourceDir("github.com/issue9/assert/v4", "./", false)
	a.NotError(err).FileExists(dir).
		True(strings.HasSuffix(filepath.ToSlash(dir), "assert/v4@v4.3.0"))

	dir, err = PkgSourceDir("github.com/issue9/assert/v4/rest", "./", false)
	a.NotError(err).FileExists(dir).
		True(strings.HasSuffix(filepath.ToSlash(dir), "assert/v4@v4.3.0/rest"))

	// replace

	dir, err = PkgSourceDir("github.com/issue9/web/v2", "./testdata/go.mod", true)
	a.NotError(err).NotEmpty(dir) // 此处 dir 可能不存在，因为 go.mod 关于 web 的包是随便指定的

	dir, err = PkgSourceDir("github.com/issue9/source", "./testdata/go.mod", true)
	a.NotError(err).FileExists(dir)

	// not exist

	dir, err = PkgSourceDir("github.com/issue9/not-exists", "./testdata/go.mod", true)
	a.ErrorIs(err, fs.ErrNotExist).Empty(dir)
}

func TestModFile(t *testing.T) {
	a := assert.New(t, false)

	p, f, err := ModFile("./")
	a.NotError(err).NotNil(f).Equal(f.Module.Mod.Path, "github.com/issue9/source").NotEmpty(p)

	p, f, err = ModFile("./testdata")
	a.NotError(err).NotNil(f).Equal(f.Module.Mod.Path, "github.com/issue9/source").NotEmpty(p)

	// NOTE: 可能不存在 c:\windows\system32 或是 c:\windows\system32 下正好存在一个 go.mod
	dir := "/windows/system32"
	if runtime.GOOS == "windows" {
		dir = "c:\\windows\\system32"
	}
	p, f, err = ModFile(dir)
	a.Error(err).Nil(f).Empty(p)
}

func TestModDir(t *testing.T) {
	a := assert.New(t, false)

	d, err := ModDir("./")
	a.NotError(err).NotEmpty(d)
}

func TestPkgPath(t *testing.T) {
	a := assert.New(t, false)

	p, err := PkgPath("./")
	a.NotError(err).Equal(p, "github.com/issue9/source")

	p, err = PkgPath("./testdata")
	a.NotError(err).Equal(p, "github.com/issue9/source/testdata")

	p, err = PkgPath("./testdata/go.mod/sub/")
	a.NotError(err).Equal(p, "github.com/issue9/source/mod/sub")

	p, err = PkgPath("./testdata/go.mod/sub/sub.go")
	a.NotError(err).Equal(p, "github.com/issue9/source/mod/sub")

	p, err = PkgPath("./testdata/go.mod/go.mod")
	a.NotError(err).Equal(p, "github.com/issue9/source/mod")
}
