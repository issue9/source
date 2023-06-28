// SPDX-License-Identifier: MIT

package source

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ModFile 查找 dir 所在模块的 go.mod 内容
//
// 从当前目录开始依次向上查找  go.mod，从其中获取 module 变量的值。
func ModFile(dir string) (*modfile.File, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

LOOP:
	for {
		path := filepath.Join(abs, "go.mod")
		stat, err := os.Stat(path)
		switch {
		case err == nil:
			if stat.IsDir() {
				abs = filepath.Dir(abs)
				continue LOOP
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return modfile.Parse(path, data, nil)
		case errors.Is(err, os.ErrNotExist):
			abs1 := filepath.Dir(abs)
			if abs1 == abs {
				return nil, os.ErrNotExist
			}
			abs = abs1
			continue LOOP
		default: // 文件存在，但是出错。
			return nil, err
		}
	}
}
