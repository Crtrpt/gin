// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"html/template"
	"runtime"
	"strconv"
	"strings"
)

const ginSupportMinGoVer = 16

// 返回true表示debug模式
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		if DebugPrintRouteFunc == nil {
			debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
		} else {
			DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
		}
	}
}

func debugPrintLoadTemplate(tmpl *template.Template) {
	if IsDebugging() {
		var buf strings.Builder
		for _, tmpl := range tmpl.Templates() {
			buf.WriteString("\t- ")
			buf.WriteString(tmpl.Name())
			buf.WriteString("\n")
		}
		debugPrint("加载 html 模板 (%d): \n%s\n", len(tmpl.Templates()), buf.String())
	}
}

// 调试打印输出
func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[调试信息] "+format, values...)
	}
}

// 获取小版本号
func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	//检查runtime 版本 和 gin 最小支持版本进行对比
	if v, e := getMinVer(runtime.Version()); e == nil && v < ginSupportMinGoVer {
		debugPrint(`[警告] gin 最低版本 Go 1.16+.

`)
	}
	debugPrint(`[警告] 创建一个 Engine 实例 开启 Logger 和 Recovery 中间件.

`)
}

func debugPrintWARNINGNew() {
	debugPrint(`[警告] 当前工作在调试模式  在生产环境需要切换到 release 模式.
 - 使用环境变量 :	export GIN_MODE=release
 - 使用代码 :	gin.SetMode(gin.ReleaseMode)

`)
}

func debugPrintWARNINGSetHTMLTemplate() {
	debugPrint(`[警告] SetHTMLTemplate() 不是线程安全的 请在注册路由和 监听socket之前设置
	router := gin.Default()
	router.SetHTMLTemplate(template) // << 好的写法

`)
}

func debugPrintError(err error) {
	if err != nil && IsDebugging() {
		fmt.Fprintf(DefaultErrorWriter, "[GIN-调试] [错误] %v\n", err)
	}
}
