// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package source

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"slices"

	"golang.org/x/mod/modfile"
)

const modFile = "go.mod"

// ModFile 文件或目录 p 所在模块的 go.mod 内容
//
// 从当前目录开始依次向上查找  go.mod，从其中获取 go.mod 文件位置，以及文件内容的解析。
func ModFile(p string) (string, *modfile.File, error) {
	path, err := modDir(p)
	if err != nil {
		return "", nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	mod, err := modfile.Parse(path, data, nil)
	if err != nil {
		return "", nil, err
	}

	return path, mod, nil
}

// ModDir 向上查找 p 所在的目录的 go.mod
func ModDir(p string) (string, error) {
	dir, err := modDir(p)
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

func modDir(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

LOOP:
	for {
		path := filepath.Join(abs, modFile)
		stat, err := os.Stat(path)
		switch {
		case err == nil:
			if stat.IsDir() {
				abs = filepath.Dir(abs)
				continue LOOP
			}

			return path, nil
		case errors.Is(err, os.ErrNotExist):
			abs1 := filepath.Dir(abs)
			if abs1 == abs {
				return "", os.ErrNotExist
			}
			abs = abs1
			continue LOOP
		default: // 文件存在，但是出错。
			return "", err
		}
	}
}

// PkgPath 文件或目录 p 所在 Go 文件的导出路径
//
// 会向上查找 go.mod，根据 go.mod 中的 module 结合当前目录组成当前目录的导出路径。
func PkgPath(p string) (string, error) {
	abs, err := filepath.Abs(p)
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
		p := filepath.Join(abs, modFile)
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
			slices.Reverse(pkgNames)
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
