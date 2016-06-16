package assert

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

func Cond(cond bool, msg string) {
	if !cond {
		panic(getErr(msg))
	}
}

func NilErr(err error) {
	if err != nil {
		panic(getErr(err.Error()))
	}
}
