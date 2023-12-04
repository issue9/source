// SPDX-License-Identifier: MIT

// Package source 提供与 Go 源码相关的一些操作
package source

import (
	"bytes"
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/issue9/errwrap"
)

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

// CurrentLocation 获取`调用者`当前的位置信息
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

// Stack 返回调用堆栈信息
//
// skip 需要忽略的内容。
//
//   - 1 表示 Stack 自身；
//   - 2 表示 Stack 的调用者，以此类推；
//
// msg 表示需要输出的额外信息；
func Stack(skip int, msg ...interface{}) string {
	buf := &bytes.Buffer{}
	DumpStack(buf, skip+1, msg...)
	return buf.String()
}

// DumpStack 将调用的堆栈信息写入 w
//
// skip 需要忽略的内容。
//
//   - 1 表示 Stack 自身；
//   - 2 表示 Stack 的调用者，以此类推；
//
// msg 表示需要输出的额外信息；
func DumpStack(w io.Writer, skip int, msg ...interface{}) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	buf := errwrap.Writer{Writer: w}
	buf.Println(msg...)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		if strings.Contains(frame.File, "runtime/") { // 忽略 runtime 下的系统调用
			continue
		}

		buf.WString(frame.Function).WByte('\n').
			WByte('\t').WString(frame.File).WByte(':').WString(strconv.Itoa(frame.Line)).WByte('\n')
	}
}
