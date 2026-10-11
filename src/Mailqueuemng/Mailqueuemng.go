package Mailqueuemng

import (
	"Environments"
	"Helpers"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type QueueStruct struct {
	QueueName    string `json:"queue_name"`
	QueueID      string `json:"queue_id"`
	ArrivalTime  int    `json:"arrival_time"`
	MessageSize  int    `json:"message_size"`
	ForcedExpire bool   `json:"forced_expire"`
	Sender       string `json:"sender"`
	Subject      string `json:"subject,omitempty"`
	Recipients   []struct {
		Address     string `json:"address"`
		DelayReason string `json:"delay_reason"`
	} `json:"recipients"`
}

var MqEnv *Environments.Env

func QueueListCheckFlag() {
	if len(Start) == 0 || len(End) == 0 {
		Helpers.HandleError(nil, 7501, "Missed parameters")
		os.Exit(1)
	}
}
func QueueList() {
	command := MqEnv.MailDir + "/sbin/" + "postqueue"
	_, response := Helpers.Exec(command, "-c", MqEnv.MailDir+"/etc", "-j")

	scanner := bufio.NewScanner(strings.NewReader(response))
	var list []interface{}
	for scanner.Scan() {
		var queueStruct QueueStruct
		_ = json.Unmarshal([]byte(scanner.Text()), &queueStruct)
		contentCommand := MqEnv.MailDir + "/sbin/" + "postcat"
		_, responseContent := Helpers.Exec(contentCommand, "-c", MqEnv.MailDir+"/etc", "-q", queueStruct.QueueID)
		scannerContent := bufio.NewScanner(strings.NewReader(responseContent))
		for scannerContent.Scan() {
			if matched, _ := regexp.MatchString("^Subject:", scannerContent.Text()); matched {
				subject := Helpers.RemoveAll(scannerContent.Text(), "Subject:")
				queueStruct.Subject = subject
				break
			}
		}

		list = append(list, queueStruct)
	}
	result, _ := json.Marshal(list)
	fmt.Println(string(result))
}

func FlushQueue() {
	command := MqEnv.MailDir + "/sbin/" + "postqueue"
	Helpers.ExecWithoutOutPut(command, "-c", MqEnv.MailDir+"/etc", "-f")
}

func ClearList() {
	command := MqEnv.MailDir + "/sbin/" + "postsuper"
	Helpers.ExecWithoutOutPut(command, "-c", MqEnv.MailDir+"/etc", "-d", "ALL")
}

func RemoveQueueCheckFlag() {
	if len(MailId) == 0 {
		Helpers.HandleError(nil, 7502, "Missed parameters")
		os.Exit(1)
	}
}
func RemoveQueue() {
	command := MqEnv.MailDir + "/sbin/" + "postsuper"
	Helpers.ExecWithoutOutPut(command, "-c", MqEnv.MailDir+"/etc", "-d", MailId)
}
