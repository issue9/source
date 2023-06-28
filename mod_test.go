// SPDX-License-Identifier: MIT

package source

import (
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
	f, err = ModFile("c:\\windows\\system32")
	a.Error(err).Nil(f)
}
