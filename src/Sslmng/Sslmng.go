package Sslmng

import (
	"Environments"
	"Helpers"
	"os"
)

var Env *Environments.Env

const DIRECTORYSEPERATOR = "/"

func CreateCheckFlag() {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 7605, "Missing parameters")
		os.Exit(7601)
	}
}
func RenewCheckFlag() {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 7605, "Missing parameters")
		os.Exit(7601)
	}
}

func Create() {

}
func Renew() {

}
