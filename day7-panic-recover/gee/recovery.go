package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//在 gee 中添加一个非常简单的错误处理机制，
//即在此类错误发生时，向用户返回 Internal Server Error，
//并且在日志中打印必要的错误信息，方便进行错误定位。
//实现中间件 Recovery

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	//Callers 用来返回调用栈的程序计数器,
	//第 0 个 Caller 是 Callers 本身，
	//第 1 个是上一层 trace，
	//第 2 个是再上一层的 defer func。
	//因此，为了日志简洁一点，我们跳过了前 3 个 Caller。
	//一个堆栈
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		//通过 runtime.FuncForPC(pc) 获取对应的函数，
		//在通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中。
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

//使用 defer 挂载上错误恢复的函数，在这个函数中调用 recover()，
//捕获 panic，并且将堆栈信息打印在日志中，向用户返回 Internal Server Error。
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
