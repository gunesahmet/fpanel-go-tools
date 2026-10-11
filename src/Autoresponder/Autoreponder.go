package Autoresponder

import (
	"Environments"
	"Helpers"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Env *Environments.Env

const DIRECTORYSEPERATOR = "/"

var RESPONSE_DIR string
var SENDMAIL_BIN string

// var config *ini.File

func isRegularFile(name string) bool {
	st, err := os.Lstat(name)
	if err != nil || !st.Mode().IsRegular() {
		return false
	}

	return true
}

func sendMail(from string, to string, msg string) error {
	cmd := exec.Command(SENDMAIL_BIN, "-i", "-f", from, to)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	stdin.Write([]byte(msg))
	stdin.Close()

	sentBytes, _ := ioutil.ReadAll(stdout)
	err = cmd.Wait()
	if err != nil {
		return err
	}

	Env.Logger.Debug(1401, "Send Command Output %s", string(sentBytes))
	return nil
}

type responseSQL struct {
	subject string
	content string
}

func getDB() (*sql.DB, *sql.DB) {
	dbM, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", Env.MailDb))
	if err != nil {
		Env.Logger.Error(1401, "Can not open: %s", Env.MailDb)
		os.Exit(1)
	}
	dbA, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", Env.AutoRespondDb))
	if err != nil {
		Env.Logger.Error(1402, "Can not open: %s", Env.AutoRespondDb)
		os.Exit(1)
	}
	return dbM, dbA
}

func getResponseFromDb() (string, error) {
	dbM, dbA := getDB()

	queryMail := fmt.Sprintf("SELECT autoresponder FROM mailbox WHERE mailbox.mailaddr='%v';", Recipient)
	row := dbM.QueryRow(queryMail)

	var autoresponder int
	err := row.Scan(&autoresponder)
	if err != nil {
		Env.Logger.Debug(1402, "Query %s", queryMail)
		Env.Logger.Error(1403, "Error: %v", err)
		os.Exit(1)
	}
	if autoresponder != 1 {
		Env.Logger.Error(1404, "User has not autoresponder: %s", Recipient)
		os.Exit(1)
	}
	err = dbM.Close()

	currentTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	queryAutoresponder := fmt.Sprintf("SELECT subject, content FROM autoresponse WHERE autoresponse.mailaddr='%v' "+
		"AND ((from_time IS NULL AND to_time IS NULL) OR (%v >= from_time AND %v <= to_time));",
		Recipient,
		currentTime,
		currentTime)

	Env.Logger.Debug(1403, "Query %s", queryAutoresponder)

	rows, err := dbA.Query(queryAutoresponder)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var key responseSQL
	for rows.Next() {
		err = rows.Scan(&key.subject, &key.content)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()

	if err != nil {
		return "", fmt.Errorf("no entries for recipient %s", Recipient)
	}

	err = dbA.Close()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`From: %v
To: %v
Content-Type: text/plain; charset=UTF-8
Subject: %v
%v`, Recipient, Sender, key.subject, key.content), nil
}

// Forward email using supplied arguments and stdin (email body)
func forwardEmailAndAutoresponse(responseRate uint) error {
	recipientRateLog := filepath.Join(RESPONSE_DIR, Recipient)
	recipientSenderRateLog := filepath.Join(RESPONSE_DIR, Recipient, Sender)

	response, err := getResponseFromDb()
	if err == nil {
		// Check rate log
		sendResponse := true
		if isRegularFile(recipientSenderRateLog) {
			curTime := time.Now()
			st, err := os.Stat(recipientSenderRateLog)
			if err != nil {
				return err
			}
			modTime := st.ModTime()

			if int64(curTime.Sub(modTime))/int64(time.Second) < int64(responseRate) {
				sendResponse = false
				Env.Logger.Info(1401, "Autoresponse has already been sent from %v to %v within last %v seconds", Recipient, Sender, responseRate)
			}
		}

		if sendResponse && strings.ToLower(Recipient) != strings.ToLower(Sender) {
			Env.Logger.Info(1402, "Sending Response from %s to %s", Recipient, Sender)

			sendMail(Recipient, Sender, response)
			if !Helpers.DirExists(recipientRateLog) {
				err = os.MkdirAll(recipientRateLog, 0770)
				if err != nil {
					return err
				}
			}

			// Touch rate log file
			fl, err := os.Create(recipientSenderRateLog)
			if err != nil {
				return err
			}
			fl.Close()
			Env.Logger.Info(1403, "Autoresponse sent from %s to %s", Recipient, Sender)
		}
	} else {
		Env.Logger.Info(1404, "No response found for %v", Sender)
	}

	// Now resend original mail
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		if fi.Size() == 0 {
			return errors.New("Stdin is empty")
		}
	}

	text, err := ioutil.ReadAll(os.Stdin)
	err = sendMail(Sender, Recipient, string(text))
	return nil
}

func Send() {
	// responseRatePtr := flag.Uint("t", 18000, "Response rate in seconds (0 - send each time)")

	RESPONSE_DIR = Env.MailDir + DIRECTORYSEPERATOR + "var/spool/autoresponse"
	SENDMAIL_BIN = Env.MailDir + DIRECTORYSEPERATOR + "sbin/sendmail"

	Env.Logger.Debug(1404, "Recipient: %s, Sender: %s, Response rate: %v", Recipient, Sender, ResponseRate)

	if Sender == "" || Recipient == "" {
		Env.Logger.Debug(1405, "Sender and recipient can not be empty")
		os.Exit(1)
	}

	Env.Logger.Info(1405, "Requested email forward from %s, to %s", Sender, Recipient)

	err := forwardEmailAndAutoresponse(ResponseRate)
	if err != nil {
		Env.Logger.Error(1405, err.Error())
		os.Exit(1)
	}
}
