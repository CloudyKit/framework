// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ensure

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type Unexpected struct {
	Msg, PackageName, FileName, FuncName string
	Line                                 int
}

func (u *Unexpected) Error() string {
	return fmt.Sprintf("(%s).%s => %s on file %s line %d", u.PackageName, u.FuncName, u.Msg, u.FileName, u.Line)
}

func getErr(msg string) *Unexpected {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &Unexpected{
		Msg:         msg,
		PackageName: packageName,
		FileName:    fileName,
		FuncName:    funcName,
		Line:        line,
	}
}

func Condition(cond bool, msg string) {
	if !cond {
		panic(getErr(msg))
	}
}

func Nil(value interface{}, msg string) {
	if value != nil {
		panic(getErr(msg))
	}
}

func NilErr(err error) {
	if err != nil {
		panic(getErr(err.Error()))
	}
}
