package Logger

import (
	"fmt"
	"os"
	"time"
)

var _panelLogDir string
var _debug bool

func InitOld(debug bool, panelLogDir string) {
	_debug = debug
	_panelLogDir = panelLogDir
}

func DebugOld(message interface{}, code string) {
	// [2020-04-24 22:59:00.174] [DEBUG][code]    Vhost updated for domain.com
	if _debug {
		dt := time.Now()
		currentDate := dt.Format("2006-01-02 15:04:05")
		str := fmt.Sprintf("%s", message)
		str = "[" + currentDate + "] [DEBUG][" + code + "]\t" + str + "\n"
		debugFile, err := os.OpenFile(_panelLogDir+"/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			os.Exit(0)
		}
		defer func() {
			err := debugFile.Close()
			if err != nil {
				os.Exit(0)
			}
		}()
		_, err = debugFile.WriteString(str)
		if err != nil {
			os.Exit(0)
		}
	}
}

func ErrorOld(message interface{}, code string) {
	dt := time.Now()
	currentDate := dt.Format("2006-01-02 15:04:05")
	str := fmt.Sprintf("%s", message)
	str = "[" + currentDate + "] [ERROR][" + code + "]\t" + str + "\n"
	fmt.Println(str)
	errorFile, err := os.OpenFile(_panelLogDir+"/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		os.Exit(0)
	}
	defer func() {
		err := errorFile.Close()
		if err != nil {
			os.Exit(0)
		}
	}()
	_, err = errorFile.WriteString(str)
	if err != nil {
		os.Exit(0)
	}
}

func SuccessOld(message interface{}, code string) {
	dt := time.Now()
	currentDate := dt.Format("2006-01-02 15:04:05")
	str := fmt.Sprintf("%v", message)
	str = "[" + currentDate + "] [SUCCESS][" + code + "]\t" + str + "\n"
	successFile, err := os.OpenFile(_panelLogDir+"/success.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		os.Exit(0)
	}
	defer func() {
		err := successFile.Close()
		if err != nil {
			os.Exit(0)
		}
	}()
	_, err = successFile.WriteString(str)
	if err != nil {
		os.Exit(0)
	}
}
