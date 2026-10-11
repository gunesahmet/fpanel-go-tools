package main

import (
	"Environments"
	"Helpers"
	"Mailqueuemng"
	"os"
	"strings"
)

func main() {

	Mailqueuemng.ArgInit()
	Environments.Init()
	Mailqueuemng.MqEnv = Environments.GetConfig()

	if Mailqueuemng.MqEnv.Debug {
		Mailqueuemng.MqEnv.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Mailqueuemng.MainProcess {
	case "queue":
		Mailqueuemng.QueueListCheckFlag()
		Mailqueuemng.QueueList()
	case "clear":
		Mailqueuemng.ClearList()
	case "flush":
		Mailqueuemng.FlushQueue()
	case "remove":
		Mailqueuemng.RemoveQueueCheckFlag()
		Mailqueuemng.RemoveQueue()
	default:
	}
}
