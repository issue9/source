source
[![Go](https://github.com/issue9/source/workflows/Go/badge.svg)](https://github.com/issue9/source/actions?query=workflow%3AGo)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/issue9/source/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/source)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/source)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/source)](https://pkg.go.dev/github.com/issue9/source)
======

source 模块提供了一些与源码相关的功能

- DumpGoSource 输出并格式化 Go 的源代码；
- CurrentFile 相当于部分语言的 `__FILE__`；
- CurrentDir 相当于部分语言的 `__DIR__`；
- CurrentLine 相当于部分语言的 `__LINE__`；
- CurrentFunction 相当于部分语言的 `__FUNCTION__`；
- Stack 返回调用者的堆栈信息；

安装
----

```shell
go get github.com/issue9/source
```

文档
----

[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/source)](https://pkg.go.dev/github.com/issue9/source)

版权
----

本项目采用 [MIT](http://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
