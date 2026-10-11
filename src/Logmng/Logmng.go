package Logmng

import (
	"Environments"
	"Helpers"
	"strconv"
)

var Env *Environments.Env

func CheckFlag() {
	if len(Message) == 0 || len(Code) == 0 {
		Helpers.HandleError(nil, 1016, "Missed parameters")
	}
}

func Error() {
	code, _ := strconv.Atoi(Code)
	Env.Logger.Error(code, Message)
}

func Success() {
	code, _ := strconv.Atoi(Code)
	Env.Logger.Success(code, Message)
}

func Debug() {
	code, _ := strconv.Atoi(Code)
	Env.Logger.Debug(code, Message)
}
