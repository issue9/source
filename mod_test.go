// SPDX-License-Identifier: MIT

package source

import (
	"runtime"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestModFile(t *testing.T) {
	a := assert.New(t, false)

	f, err := ModFile("./")
	a.NotError(err).NotNil(f).Equal(f.Module.Mod.Path, "github.com/issue9/source")

	f, err = ModFile("./testdata")
	a.NotError(err).NotNil(f).Equal(f.Module.Mod.Path, "github.com/issue9/source")

	// NOTE: 可能不存在 c:\windows\system32 或是 c:\windows\system32 下正好存在一个 go.mod
	dir := "/windows/system32"
	if runtime.GOOS == "windows" {
		dir = "c:\\windows\\system32"
	}
	f, err = ModFile(dir)
	a.Error(err).Nil(f)
}
