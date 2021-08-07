package common

import (
	"log"
	"strings"

	// "github.com/chanyipiaomiao/hltool"
	"fmt"
	"node-monitor/libs/logs"
	"os"
	"path"
)

// InitLog 初始化日志
// func InitLog() {
// 	var logpath string
// 	var logname string
// 	var logpathname string
// 	var Workdir string

// 	logpath = "log"
// 	_, err := os.Stat(logpath)
// 	if err != nil {
// 		os.Mkdir(logpath, 0755)
// 	}

// 	logname = Cfg.Name + ".log"

// 	// 设定日志级别
// 	level := Cfg.Log.LogLevel
// 	beego.SetLevel(level)

// 	//打印行号
// 	beego.SetLogFuncCall(true)

// 	//删除终端打印
// 	beego.BeeLogger.DelLogger("console")

// 	Workdir = getWorkDir()
// 	logpathname = path.Join(Workdir, logpath, logname)
// 	maxdays := Cfg.Log.Maxdays
// 	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s","maxdays":%d}`, logpathname, maxdays))

// }

// InitLog 初始化日志-
func InitLogNew() {
	var logpath string
	var logname string
	var logpathname string
	var Workdir string

	logpath = "log"
	_, err := os.Stat(logpath)
	if err != nil {
		os.Mkdir(logpath, 0755)
	}

	logname = Cfg.Name + ".log"

	// 设定日志级别
	level := Cfg.Log.LogLevel
	logs.SetLevel(level)

	//打印行号
	logs.SetLogFuncCall(true)
	// logs.EnableFuncCallDepth(true)
	// logs.SetLogFuncCallDepth(3)

	//异步输出,提升性能
	logs.Async()

	//删除终端打印
	// logs.DelLogger("console")

	Workdir = getWorkDir()
	logpathname = path.Join(Workdir, logpath, logname)
	maxdays := Cfg.Log.Maxdays
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","maxdays":%d,"color":true,"daily":true}`, logpathname, maxdays))

}

func getWorkDir() (workDir string) {
	execPath, err := getExecPath()
	if err != nil {
		log.Fatal("Fail to get work directory: %v", err)
	}
	workDir = path.Dir(strings.Replace(execPath, "\\", "/", -1))
	return
}
