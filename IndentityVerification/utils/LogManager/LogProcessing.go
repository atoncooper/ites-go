package LogManager

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
	"time"
)

type Info struct {
	InfoLogFile     *os.File
	InfoLogger      log.Logger
	WarnInfoFile    *os.File
	WarnInfoLogger  log.Logger
	ErrorInfoFile   *os.File
	ErrorInfoLogger log.Logger
}

func (i *Info) NewLogManager() {
	var err error

	// 打开 info.log 文件
	i.InfoLogFile, err = os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// 创建 info 日志记录器
	i.InfoLogger = log.NewLogfmtLogger(i.InfoLogFile)
	i.InfoLogger = log.NewSyncLogger(i.InfoLogger)
	i.InfoLogger = level.NewFilter(i.InfoLogger, level.AllowInfo())

	// 打开 warn.log 文件
	i.WarnInfoFile, err = os.OpenFile("warn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// 创建 warn 日志记录器
	i.WarnInfoLogger = log.NewLogfmtLogger(i.WarnInfoFile)
	i.WarnInfoLogger = log.NewSyncLogger(i.WarnInfoLogger)
	i.WarnInfoLogger = level.NewFilter(i.WarnInfoLogger, level.AllowWarn())

	// 打开 error.log 文件
	i.ErrorInfoFile, err = os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// 创建 error 日志记录器
	i.ErrorInfoLogger = log.NewLogfmtLogger(i.ErrorInfoFile)
	i.ErrorInfoLogger = log.NewSyncLogger(i.ErrorInfoLogger)
	i.ErrorInfoLogger = level.NewFilter(i.ErrorInfoLogger, level.AllowError())
}

func (i *Info) CloseLogFiles() {
	if i.InfoLogFile != nil {
		err := i.InfoLogFile.Close()
		if err != nil {
			return
		}
	}
	if i.WarnInfoFile != nil {
		err := i.WarnInfoFile.Close()
		if err != nil {
			return
		}
	}
	if i.ErrorInfoFile != nil {
		err := i.ErrorInfoFile.Close()
		if err != nil {
			return
		}
	}
}

func (i *Info) InfoLog(serverName string, detail error) {
	infoData := fmt.Sprintf("Server-Name:%s ||  Detail : %s ||  Time : %v",
		serverName,
		detail,
		time.Now(),
	)
	err := level.Info(i.InfoLogger).Log("InfoMessage", infoData)
	if err != nil {
		fmt.Println("Error writing info log:", err)
	}
}

func (i *Info) WarnLog(serverName string, detail error) {
	warnData := fmt.Sprintf("Server-Name:%s ||  Detail : %s ||  Time : %v",
		serverName,
		detail,
		time.Now(),
	)
	err := level.Warn(i.WarnInfoLogger).Log("WarnMessage", warnData)
	if err != nil {
		fmt.Println("Error writing warn log:", err)
	}
}

func (i *Info) ErrorLog(serverName string, detail error) {
	errorData := fmt.Sprintf("Server-Name:%s ||  Detail : %s ||  Time : %v",
		serverName,
		detail,
		time.Now(),
	)
	err := level.Error(i.ErrorInfoLogger).Log("ErrorMessage", errorData)
	if err != nil {
		fmt.Println("Error writing error log:", err)
	}
}
