// SPDX-License-Identifier: MIT

package source

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/issue9/sliceutil"
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

// ModPath 目录 dir 所在 Go 文件的导出路径
func ModPath(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		abs = filepath.Dir(abs)
	}

	pkgNames := make([]string, 0, 10)
LOOP:
	for {
		p := filepath.Join(abs, "go.mod")
		stat, err := os.Stat(p)
		switch {
		case err == nil:
			if stat.IsDir() { // 名为 go.mod 的目录
				pkgNames = append(pkgNames, filepath.Base(abs))
				abs = filepath.Dir(abs)
				continue LOOP
			}

			data, err := os.ReadFile(p)
			if err != nil {
				return "", err
			}
			mod, err := modfile.Parse(p, data, nil)
			if err != nil {
				return "", err
			}

			pkgNames = append(pkgNames, mod.Module.Mod.Path)
			sliceutil.Reverse(pkgNames)
			return path.Join(pkgNames...), nil
		case errors.Is(err, os.ErrNotExist):
			// 这两行不能用 filepath.Split 代替，split 会为 abs1 留下最后的分隔符，
			// 导致下一次的 filepath.Split 返回空的 file 值。
			base := filepath.Base(abs)
			abs1 := filepath.Dir(abs)

			if abs1 == abs { // 到根目录了
				return "", os.ErrNotExist
			}

			abs = abs1
			pkgNames = append(pkgNames, base)
			continue LOOP
		default: // 文件存在，但是出错。
			return "", err
		}
	}
}
