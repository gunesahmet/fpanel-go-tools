package Sqlmng

import (
	"Environments"
	"Helpers"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"os"
	"strconv"
	"strings"
)

var Env *Environments.Env

const DIRECTORYSEPERATOR = "/"
const DUMPEXTENTION = ".dump"

var postgreAppDataDir = "C:/Users/Administrator/AppData/Roaming/postgresql/"

func backupMySql() {
	var dumpPath string
	if len(CloudServerIp) > 0 {
		dumpPath = Env.PanelTmpDir + "/mysql_" + DbName + DUMPEXTENTION
	} else {
		dumpPath = BackupPath + "/mysql_" + DbName + DUMPEXTENTION
		if !Helpers.DirExists(BackupPath) {
			Helpers.HandleError(nil, 7701, fmt.Sprintf("MySQL backup path doesnt exists %s", dumpPath))
			os.Exit(1)
		}
	}

	var err error
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		err, _ = Helpers.Exec("/bin/sh", "-c", "mysqldump --opt --host=127.0.0.1 --port="+Port+" --user="+DbUser+" --password="+DbPass+" "+DbName+" > "+dumpPath)
	} else {
		err, _ = Helpers.Exec("cmd", "/c", Env.PanelDir+DIRECTORYSEPERATOR+"mysql/bin/mysqldump.exe --opt --host=127.0.0.1 --port="+Port+" --user="+DbUser+" --password="+DbPass+" "+DbName+" > "+dumpPath)
	}

	if err != nil {
		Helpers.HandleError(nil, 7702, fmt.Sprintf("MySQL Error: %v", err))
	} else {
		SendToMainServer(CloudServerIp, CloudKey, CloudIv, "mysql_"+DbName+DUMPEXTENTION, dumpPath, BackupPath)
	}
	_ = Helpers.RemoveFile(dumpPath)
}
func backupSqlServer() {
	var dumpPath string
	if len(CloudServerIp) > 0 {
		dumpPath = Env.PanelTmpDir + "/mssql_" + DbName + DUMPEXTENTION
	} else {
		dumpPath = BackupPath + "/mssql_" + DbName + DUMPEXTENTION
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s", "127.0.0.1", "sa", Env.FpanelAdminPasswd, Port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		Helpers.HandleError(err, 7703, fmt.Sprintf("Microsoft SQL connection failed %s", connString))
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()

	rows, err := conn.Query(fmt.Sprintf("use master; ALTER DATABASE %s SET RECOVERY FULL; BACKUP DATABASE [%s] TO DISK = N'%s' WITH NOFORMAT, NOINIT, NAME = N'%s-Full Database Backup', SKIP, NOREWIND, NOUNLOAD, STATS = 10;",
		DbName, DbName, dumpPath, DbName))
	if err != nil {
		Helpers.HandleError(err, 7704, fmt.Sprintf("Backup SQL failed %v", err))
	}
	for rows.Next() {
	}
	SendToMainServer(CloudServerIp, CloudKey, CloudIv, "/mssql_"+DbName+DUMPEXTENTION, dumpPath, BackupPath)
	_ = Helpers.RemoveFile(dumpPath)
}
func backupPostgreSql() {
	var dumpPath string
	if len(CloudServerIp) > 0 {
		dumpPath = Env.PanelTmpDir + "/pgsql_" + DbName + DUMPEXTENTION
	} else {
		dumpPath = BackupPath + "/pgsql_" + DbName + DUMPEXTENTION
		if !Helpers.DirExists(BackupPath) {
			Helpers.HandleError(nil, 7705, fmt.Sprintf("PostgreSql backup path doesnt exists %s", dumpPath))
			os.Exit(1)
		}
	}

	var err error
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		err, _ = Helpers.Exec("/bin/sh", "-c", "PGPASSWORD="+DbPass+" pg_dump"+" -h 127.0.0.1 -p "+Port+" -U "+DbUser+" -O -F c -b -f "+dumpPath+" "+DbName)
	} else {
		writePgPassFile(Port, DbName, DbUser, DbPass)
		err, _ = Helpers.Exec("cmd", "/c", "c:/fpanel/postgresql/bin/pg_dump.exe -h 127.0.0.1 -p "+Port+" -U "+DbUser+" -O -F c -b -f "+dumpPath+" "+DbName)
		_ = Helpers.RemoveFile(postgreAppDataDir + "pgpass.conf")
	}
	if err != nil {
		Helpers.HandleError(err, 7706, fmt.Sprintf("PostgreSql Error: %v", err))
	} else {
		SendToMainServer(CloudServerIp, CloudKey, CloudIv, "pgsql_"+DbName+DUMPEXTENTION, dumpPath, BackupPath)
	}
	_ = Helpers.RemoveFile(dumpPath)
}

func BackupCheckFlag() {
	if len(Type) == 0 || len(DbPass) == 0 || len(DbName) == 0 || len(DbUser) == 0 || len(Host) == 0 || len(Port) == 0 ||
		len(Date) == 0 || len(BackupPath) == 0 {
		Helpers.HandleError(nil, 7707, "Missing parameters")
		os.Exit(1)
	}
}
func Backup() {
	switch Type {
	case "mysql":
		backupMySql()
		break
	case "mssql":
		backupSqlServer()
		break
	case "pgsql":
		backupPostgreSql()
		break
	}
}

func restoreMySql() {
	dumpPath := BackupPath + "/mysql_" + OldDbName + DUMPEXTENTION
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		dumpPath = Helpers.ReplaceAll(dumpPath, "c:/fpanel", Env.PanelDir)
	} else {
		dumpPath = Helpers.ReplaceAll(dumpPath, "/usr/local/fpanel", Env.PanelDir)
	}
	if !Helpers.FileExists(dumpPath) {
		Helpers.HandleError(nil, 7708, fmt.Sprintf("File: %s does not exists", dumpPath))
		os.Exit(1)
	}
	var err error

	if strings.Compare(Env.SystemOS, "linux") == 0 {
		host := Host
		if strings.Compare(Env.SystemType, "centos") == 0 {
			host = "127.0.0.1"
		}
		err, _ = Helpers.Exec("/bin/sh", "-c",
			"mysql --host="+host+
				" --port="+Port+
				" --user="+DbUser+
				" --password="+DbPass+
				" "+DbName+
				" < "+dumpPath)
	} else {
		err, _ = Helpers.Exec("cmd", "/c", Env.PanelDir+"/mysql/bin/mysql.exe"+
			"--host="+Host+
			" --port="+Port+
			" --user="+DbUser+
			" --password="+DbPass+
			" "+DbName+
			" < "+dumpPath)
	}

	if err != nil {
		Helpers.HandleError(err, 7709, fmt.Sprintf("MySQL restore has an error %s. Error: %v", dumpPath, err))
	}
	//_ = Helpers.RemoveFile(dumpPath)
}
func restoreSqlServer() {
	dumpPath := BackupPath + "/mssql_" + OldDbName + DUMPEXTENTION
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		dumpPath = Helpers.ReplaceAll(dumpPath, "c:/fpanel", Env.PanelDir)
	} else {
		dumpPath = Helpers.ReplaceAll(dumpPath, "/usr/local/fpanel", Env.PanelDir)
	}
	if !Helpers.FileExists(dumpPath) {
		Helpers.HandleError(nil, 7710, fmt.Sprintf("File: %s does not exists", dumpPath))
		os.Exit(1)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s", "127.0.0.1", "sa", Env.FpanelAdminPasswd, Port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		Helpers.HandleError(err, 7711, fmt.Sprintf("Microsoft SQL connection failed %s", connString))
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()

	////////////////
	query := fmt.Sprintf("use master; ALTER DATABASE %s SET RECOVERY FULL; BACKUP DATABASE [%s] TO DISK = N'%s' WITH NOFORMAT, NOINIT, NAME = N'%s-Full Database Backup', SKIP, NOREWIND, NOUNLOAD, STATS = 10;",
		DbName, DbName, dumpPath, DbName)
	rows, err := conn.Query(query)
	if err != nil {
		Helpers.HandleError(err, 7712, fmt.Sprintf("Backup SQL failed %v", err))
	}
	for rows.Next() {
	}
	////////////////
	_ = Helpers.RemoveFile(dumpPath)
}
func restorePostgreSql() {
	dumpPath := BackupPath + "/pgsql_" + OldDbName + DUMPEXTENTION
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		dumpPath = Helpers.ReplaceAll(dumpPath, "c:/fpanel", Env.PanelDir)
	} else {
		dumpPath = Helpers.ReplaceAll(dumpPath, "/usr/local/fpanel", Env.PanelDir)
	}

	if !Helpers.FileExists(dumpPath) {
		Helpers.HandleError(nil, 7713, fmt.Sprintf("File: %s does not exists", dumpPath))
		os.Exit(1)
	}

	var err error

	if strings.Compare(Env.SystemOS, "linux") == 0 {
		err, _ = Helpers.Exec("/bin/sh", "-c", "PGPASSWORD="+DbPass+" pg_restore"+" -h 127.0.0.1 -p "+
			Port+" -U "+DbUser+" -d "+DbName+" "+dumpPath)
	} else {
		writePgPassFile(Port, DbName, DbUser, DbPass)
		err, _ = Helpers.Exec("cmd", "/c", "c:/fpanel/postgresql/bin/pg_restore -h 127.0.0.1 -p "+
			Port+" -U "+DbUser+" -d "+DbName+" "+dumpPath)
		_ = Helpers.RemoveFile(postgreAppDataDir + "pgpass.conf")
	}

	if err != nil {
		Helpers.HandleError(err, 7714, fmt.Sprintf("PostgreSQL restore has an error %s", dumpPath))
	}
	_ = Helpers.RemoveFile(dumpPath)
}

func RestoreCheckFlag() {
	if len(Type) == 0 || len(DbPass) == 0 || len(DbName) == 0 ||
		len(DbUser) == 0 || len(Host) == 0 || len(Port) == 0 ||
		len(OldDbName) == 0 || len(BackupPath) == 0 {
		Helpers.HandleError(nil, 7715, "Missing parameters")
		os.Exit(1)
	}
}
func Restore() {
	switch Type {
	case "mysql":
		restoreMySql()
		break
	case "mssql":
		restoreSqlServer()
		break
	case "pgsql":
		restorePostgreSql()
		break
	}
}

func SettingsCheckFlag() {

}
func Settings() {
	var myCnf string
	if (strings.Compare(Env.SystemType, "debian") == 0) || (strings.Compare(Env.SystemType, "ubuntu") == 0) {
		if Helpers.FileExists("/etc/mysql/mysql.conf.d/mysqld.cnf") {
			myCnf = "/etc/mysql/mysql.conf.d/mysqld.cnf"
		} else if Helpers.FileExists("/etc/mysql/my.cnf") {
			myCnf = "/etc/mysql/my.cnf"
		}
	} else if strings.Compare(Env.SystemType, "centos") == 0 {
		myCnf = "/etc/my.cnf"
	} else if strings.Compare(Env.SystemType, "windows") == 0 {
		myCnf = "c:/fpanel/mysql/etc/my.ini"
	}

	err := Helpers.CopyFile(myCnf, myCnf+"_"+strconv.FormatInt(Helpers.CurrentTimestamp(), 10)+".backup")
	if err != nil {
		Helpers.HandleError(err, 7716, "Setting MySQL error")
		os.Exit(0)
	}

	my, err := NewMySQLSettings(myCnf)
	if err != nil {
		Helpers.HandleError(err, 7717, "Setting MySQL error")
		os.Exit(0)
	}

	if len(RemoteConnection) > 0 {
		err = my.MySQLRemoteConnection(RemoteConnection)
		if err != nil {
			Helpers.HandleError(err, 7718, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(BufferPoolSize) > 0 {
		err = my.MySQLBufferPoolSize(BufferPoolSize)
		if err != nil {
			Helpers.HandleError(err, 7719, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(LogFileSize) > 0 {
		err = my.MySQLLogFileSize(LogFileSize)
		if err != nil {
			Helpers.HandleError(err, 7720, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(FilePerTable) > 0 {
		err = my.MySQLFilePerTable(FilePerTable)
		if err != nil {
			Helpers.HandleError(err, 7721, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(FlushMethod) > 0 {
		err = my.MySQLFlushMethod(FlushMethod)
		if err != nil {
			Helpers.HandleError(err, 7722, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(LogError) > 0 {
		err = my.MySQLLogError(LogError)
		if err != nil {
			Helpers.HandleError(err, 7723, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(SlowQueryLog) > 0 {
		err = my.MySQLSlowQueryLog(SlowQueryLog)
		if err != nil {
			Helpers.HandleError(err, 7724, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(SlowQueryLogFile) > 0 {
		err = my.MySQLSlowQueryLogFile(SlowQueryLogFile)
		if err != nil {
			Helpers.HandleError(err, 7725, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(TmpTableSize) > 0 {
		err = my.MySQLTmpTableSize(TmpTableSize)
		if err != nil {
			Helpers.HandleError(err, 7726, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(MaxHeapTableSize) > 0 {
		err = my.MySQLMaxHeapTableSize(MaxHeapTableSize)
		if err != nil {
			Helpers.HandleError(err, 7727, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(QueryCacheType) > 0 {
		err = my.MySQLQueryCacheType(QueryCacheType)
		if err != nil {
			Helpers.HandleError(err, 7728, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(QueryCacheSize) > 0 {
		err = my.MySQLQueryCacheSize(QueryCacheSize)
		if err != nil {
			Helpers.HandleError(err, 7729, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(MaxConnections) > 0 {
		err = my.MySQLMaxConnections(MaxConnections)
		if err != nil {
			Helpers.HandleError(err, 7730, "Setting MySQL error")
			os.Exit(0)
		}
	}

	//fmt.Println(ThreadCacheSize)
	if len(ThreadCacheSize) > 0 {
		err = my.MySQLThreadCacheSize(ThreadCacheSize)
		if err != nil {
			Helpers.HandleError(err, 7731, "Setting MySQL error")
			os.Exit(0)
		}
	}
	if len(TableOpenCache) > 0 {
		err = my.MySQLTableOpenCache(TableOpenCache)
		if err != nil {
			Helpers.HandleError(err, 7732, "Setting MySQL error")
			os.Exit(0)
		}
	}

	if len(OpenFilesLimit) > 0 {
		err = my.MySQLOpenFilesLimit(OpenFilesLimit)
		if err != nil {
			Helpers.HandleError(err, 7733, "Setting MySQL error")
			os.Exit(0)
		}
	}

	err = my.ini.SaveTo(myCnf)
	if err != nil {
		Helpers.HandleError(err, 7734, "Setting MySQL error")
		os.Exit(0)
	}
}

func SendToMainServer(cloudServerIp string, cloudKey string, cloudIv string, fileName string, sourcePath string, filePath string) {
	if len(cloudServerIp) != 0 && len(cloudKey) != 0 && len(cloudIv) != 0 {
		// --ipaddr=192.168.139.53 --cloud-key=ffaaed55add60e5e139d2d1e6caa13ca --cloud-iv=afe88787a5f2f48c13aa09dce9083fec --os=linux --filename=test.conf --filepath=/root/panel.conf.backup --filetype=backup
		err, _ := Helpers.Exec(Env.Cloudmng,
			"--send-file",
			fmt.Sprintf("--ipaddr=%s", cloudServerIp),
			fmt.Sprintf("--cloud-key=%s", cloudKey),
			fmt.Sprintf("--cloud-iv=%s", cloudIv),
			fmt.Sprintf("--os=%s", strings.ToLower(Env.SystemOS)),
			fmt.Sprintf("--filename=%s", fileName),
			fmt.Sprintf("--filepath=%s", filePath),
			fmt.Sprintf("--source=%s", sourcePath),
			"--filetype=sql-backup")

		if err != nil {
			return
		}
	}
}

func Test() {
	//dbName := "test"
	//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", "127.0.0.1", "sa", Env.FpanelAdminPasswd, 1433)
	//conn, err := sql.Open("mssql", connString)
	//if err != nil {
	//	Logger.Error(fmt.Sprintf("Microsoft SQL connection failed %s", connString), "6901")
	//	os.Exit(1)
	//}
	//defer func() {
	//	_ = conn.Close()
	//}()
	//
	//stmt, err := conn.Prepare("DECLARE @filepath NVARCHAR(260)\nEXEC master.dbo.xp_instance_regread \nN'HKEY_LOCAL_MACHINE', N'Software\\Microsoft\\MSSQLServer\\MSSQLServer', \nN'DefaultData', \n@filepath output, 'no_output' \nSELECT @filepath as path")
	//
	//if err != nil {
	//	Logger.Error(fmt.Sprintf("Microsoft SQL STMT failed %v", err), "6901")
	//	os.Exit(1)
	//}
	//defer func() {
	//	_ = stmt.Close()
	//}()
	//
	//row := stmt.QueryRow()
	//var path string
	//err = row.Scan(&path)
	//if err != nil {
	//	fmt.Println("Scan failed:", err.Error())
	//}
	//
	//if len(path) > 0 {
	//	pathArr := Helpers.Split(path, "\\")
	//	pathArr = pathArr[:len(pathArr)-2]
	//	newPath := strings.Join(pathArr, "\\")
	//	newPath = newPath + "\\backup"
	//	newPathQuery := fmt.Sprintf("declare\n @dirPath varchar(1024)\nset @dirPath = '%s';\nexec xp_create_subdir @dirPath;", newPath)
	//	_, err = conn.Query(newPathQuery)
	//	if err != nil {
	//		fmt.Println("Declare Query failed:", err.Error())
	//	}
	//	backupPath := newPath + "\\" + dbName
	//	rows, err := conn.Query(fmt.Sprintf("use master; ALTER DATABASE %s SET RECOVERY FULL; BACKUP DATABASE [%s] TO DISK = N'%s.dump' WITH NOFORMAT, NOINIT, NAME = N'%s-Full Database Backup', SKIP, NOREWIND, NOUNLOAD, STATS = 10;",
	//		dbName, dbName, backupPath, dbName))
	//	if err != nil {
	//		fmt.Println("Backup Query failed:", err.Error())
	//	}
	//	for rows.Next() {
	//
	//	}
	//}
	//
	//fmt.Printf("path:%s\n", path)
}

func writePgPassFile(port string, dbName string, dbUser string, dbPass string) {
	pgPassFileContent := fmt.Sprintf("%s:%s:%s:%s:%s", "127.0.0.1", port, dbName, dbUser, dbPass)
	if !Helpers.DirExists(postgreAppDataDir) {
		Helpers.CreateDirectory(postgreAppDataDir, true)
	}
	err := Helpers.WriteFile(pgPassFileContent, postgreAppDataDir+"pgpass.conf")
	if err != nil {
		Helpers.HandleError(err, 7735, fmt.Sprintf("pgpass.conf write failed %v", err))
	}
}
