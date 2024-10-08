// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package source

import (
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const modFile = "go.mod"

var (
	pkgSource = filepath.Join(build.Default.GOPATH, "pkg", "mod")
	stdSource = filepath.Join(build.Default.GOROOT, "src")
)

// PkgSourceDir 查找包 pkgPath 的源码目录
//
// 如果 pkgPath 是标准库的名称，如 encoding/json 等，则返回当前使用的 Go 版本对应的标准库地址。
// 其它情况则从 modDir 指向的 go.mod 中查找 require 或是 replace 字段的定义，
// 并根据这些定义找到其指向的源码路径。
// 如果 modDir 中不存在 go.mod 会尝试向上一级目录查找。
//
// pkgPath 需要查找的包路径，如果指向的是模块下的包级别的导出路径，则会尝试使用 [strings.HasPrefix] 与 require 指令进行对比；
// modDir go.mod 所在的目录，将在该文件中查找 pkgPath 指定的目录；
// replace 是否考虑 go.mod 中的 replace 指令的影响；
//
// 如果找不到，会返回 [fs.ErrNotExist]
//
// NOTE: 这并不会检测 dir 指向目录是否真实且准确。
func PkgSourceDir(pkgPath, modDir string, replace bool) (dir string, err error) {
	if strings.IndexByte(pkgPath, '.') < 0 {
		return filepath.Join(stdSource, pkgPath), nil
	}

	_, mod, err := ModFile(modDir)
	if err != nil {
		return "", err
	}

	// 保证长的在前面，这样在碰到 xxx.com/pkg/v2 与 xxx.com/pkg 两个包同时出现时，v2 会出现在前面，拥有优先匹配的权利。
	slices.SortFunc(mod.Require, func(a, b *modfile.Require) int { return len(b.Mod.Path) - len(a.Mod.Path) })

	for _, pkg := range mod.Require {
		if !strings.HasPrefix(pkgPath, pkg.Mod.Path) {
			continue
		}
		suffix := strings.TrimPrefix(pkgPath, pkg.Mod.Path)
		// github.com/issue9/web 与 github.com/issue9/webuse 也匹配，但不是同一个包
		// 只有 github.com/issue9/web 与 github.com/issue9/web/v4 这种才行
		if suffix != "" && suffix[0] != '/' {
			continue
		}

		index := slices.IndexFunc(mod.Replace, func(r *modfile.Replace) bool { return r.Old.Path == pkg.Mod.Path })
		if !replace || index < 0 {
			p, err := escapePath(pkg.Mod.Path, pkg.Mod.Version, suffix)
			if err != nil {
				return "", err
			}
			return filepath.Join(pkgSource, p), nil
		}

		p := mod.Replace[index].New.Path
		if p != "" && (p[0] == '.' || p[0] == '/') { // 指向本地
			if !filepath.IsAbs(p) {
				p = filepath.Join(modDir, p)
			}
			return filepath.Abs(p)
		}
		return PkgSourceDir(p, modDir, false)
	}

	return "", fs.ErrNotExist
}

func escapePath(p, v, s string) (path string, err error) {
	p, err = module.EscapePath(p)
	if err != nil {
		return "", err
	}

	v, err = module.EscapeVersion(v)
	if err != nil {
		return "", err
	}

	// 不对 s 作转码，同一个项目下应该不至于有不同大小写的同名包。

	return p + "@" + v + s, nil
}

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

// PkgPath 文件或目录 p 的导出路径
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
