// SPDX-License-Identifier: MIT

// Package source 提供与 Go 源码相关的一些操作
package source

import (
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/issue9/errwrap"
)

// DumpGoSource 输出 Go 源码到 path
//
// 会对源代码作格式化。
func DumpGoSource(path string, content []byte) error {
	src, err := format.Source(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, src, os.ModePerm)
}

// CurrentPath 获取`调用者`所在目录的路径
//
// 类似于部分语言的的 __DIR__ + "/" + path
func CurrentPath(path string) string {
	_, fi, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(fi), path)
}

// CurrentDir 获取`调用者`所在的目录
//
// 相当于部分语言的 __DIR__
func CurrentDir() string {
	_, fi, _, _ := runtime.Caller(1)
	return filepath.Dir(fi)
}

// CurrentFile 获取`调用者`所在的文件
//
// 相当于部分语言的 __FILE__
func CurrentFile() string {
	_, fi, _, _ := runtime.Caller(1)
	return fi
}

// CurrentLine 获取`调用者`所在的行
//
// 相当于部分语言的 __LINE__
func CurrentLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// CurrentLocation 获取调用者当前的位置信息
func CurrentLocation() (path string, line int) {
	_, path, line, _ = runtime.Caller(1)
	return path, line
}

// CurrentFunction 获取`调用者`所在的函数名
//
// 相当于部分语言的 __FUNCTION__
func CurrentFunction() string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	index := strings.LastIndexByte(name, '.')
	if index > 0 {
		name = name[index+1:]
	}

	return name
}

// Stack 返回调用者的堆栈信息
//
// skip 需要忽略的内容。
// 1 表示 Stack 自身， 2 表示 TraceStack 的调用者，以此类推；
// msg 表示需要输出的额外信息；
func Stack(skip int, msg ...interface{}) (string, error) {
	var w errwrap.StringBuilder

	if len(msg) > 0 {
		w.Println(msg...)
	}

	for i := skip; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		w.WString(fn.Name()).WByte('\n')
		w.WByte('\t').WString(file).WByte(':').WString(strconv.Itoa(line)).WByte('\n')
	}

	if w.Err != nil {
		return "", w.Err
	}
	return w.String(), nil
}
