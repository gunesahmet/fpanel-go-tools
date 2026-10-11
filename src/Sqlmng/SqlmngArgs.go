package Sqlmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Host          string `arg:"--host"`
	Port          string `arg:"--port"`
	DbName        string `arg:"--dbname"`
	BackupPath    string `arg:"--backup-path"`
	DumpPath      string `arg:"--dump-path"`
	OldDbName     string `arg:"--old-dbname"`
	DbUser        string `arg:"--dbuser"`
	DbPass        string `arg:"--dbpass"`
	Date          string `arg:"--date"`
	Type          string `arg:"--type"`
	CloudServerIp string `arg:"--cloud-server-ip"`
	CloudKey      string `arg:"--cloud-key"`
	CloudIv       string `arg:"--cloud-iv"`

	RemoteConnection string `arg:"--remote-connection"`
	BufferPoolSize   string `arg:"--buffer-pool-size"`
	LogFileSize      string `arg:"--log-file-size"`
	FilePerTable     string `arg:"--file-per-table"`
	FlushMethod      string `arg:"--flush-method"`
	LogError         string `arg:"--log-error"`
	SlowQueryLog     string `arg:"--slow-query-log"`
	SlowQueryLogFile string `arg:"--slow-query-log-file"`
	TmpTableSize     string `arg:"--tmp-table-size"`
	MaxHeapTableSize string `arg:"--max-heap-table-size"`
	QueryCacheType   string `arg:"--query-cache-type"`
	QueryCacheSize   string `arg:"--query-cache-size"`
	MaxConnections   string `arg:"--max-connection"`
	ThreadCacheSize  string `arg:"--thread-cache-size"`
	TableOpenCache   string `arg:"--table-open-cache"`
	OpenFilesLimit   string `arg:"--open-file-limit"`

	Backup   bool `arg:"--backup"`
	Restore  bool `arg:"--restore"`
	Settings bool `arg:"--settings"`
	Test     bool `arg:"--test"`
}

var Host string
var Port string
var DbName string
var OldDbName string
var DbUser string
var DbPass string
var Date string
var BackupPath string
var DumpPath string
var Type string

var CloudServerIp string
var CloudKey string
var CloudIv string

var RemoteConnection string
var BufferPoolSize string
var LogFileSize string
var FilePerTable string
var FlushMethod string
var LogError string
var SlowQueryLog string
var SlowQueryLogFile string
var TmpTableSize string
var MaxHeapTableSize string
var QueryCacheType string
var QueryCacheSize string
var MaxConnections string
var ThreadCacheSize string
var TableOpenCache string
var OpenFilesLimit string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	Host = args.Host
	Port = args.Port
	DbName = args.DbName
	OldDbName = args.OldDbName
	BackupPath = args.BackupPath
	DumpPath = args.DumpPath
	DbUser = args.DbUser
	DbPass = args.DbPass
	Date = args.Date
	Type = args.Type

	CloudServerIp = args.CloudServerIp
	CloudKey = args.CloudKey
	CloudIv = args.CloudIv

	RemoteConnection = args.RemoteConnection
	BufferPoolSize = args.BufferPoolSize
	LogFileSize = args.LogFileSize
	FilePerTable = args.FilePerTable
	FlushMethod = args.FlushMethod
	LogError = args.LogError
	SlowQueryLog = args.SlowQueryLog
	SlowQueryLogFile = args.SlowQueryLogFile
	TmpTableSize = args.TmpTableSize
	MaxHeapTableSize = args.MaxHeapTableSize
	QueryCacheType = args.QueryCacheType
	QueryCacheSize = args.QueryCacheSize
	MaxConnections = args.MaxConnections
	ThreadCacheSize = args.ThreadCacheSize
	TableOpenCache = args.TableOpenCache
	OpenFilesLimit = args.OpenFilesLimit

	if args.Backup {
		MainProcess = "backup"
	} else if args.Restore {
		MainProcess = "restore"
	} else if args.Settings {
		MainProcess = "settings"
	} else if args.Test {
		MainProcess = "test"
	}
}
