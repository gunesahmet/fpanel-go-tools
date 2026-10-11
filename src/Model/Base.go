package Model

import (
	"Environments"
	"gorm.io/gorm"
)

var Env *Environments.Env

type Connection struct {
	DB *gorm.DB
}

var Conn Connection

const BACKUP_WAITING = "0"
const BACKUP_FINISHED = "1"
const BACKUP_RESTORING = "2"
const BACKUP_READING = "3"
const BACKUP_FAILED = "4"
const RESTORE_FAILED = "5"
