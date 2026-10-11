package Statisticsmng

import (
	Crypt "CrypterWin"
	"Environments"
	"Gorm"
	"Helpers"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

func calculateDatabases(domain Gorm.Domains) int64 {
	var totalMicrosoftSqlSize int64 = 0
	var totalMySqlSize int64 = 0
	var totalPostgreSqlSize int64 = 0
	for _, db := range domain.Databases {
		var databaseServer Gorm.DatabaseServer
		var databaseUser Gorm.DatabaseUsers
		res := DbConn.DB.Find(&databaseServer, "id = ?", db.DatabaseServerId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Database server not found or has an error: %v", res.Error))
			continue
		}
		resDbUser := DbConn.DB.Find(&databaseUser, "id = ?", db.DatabaseUserId)
		if resDbUser.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Database user not found or has an error: %v", res.Error))
			continue
		}
		defaultCrypto, _ := Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
		dbName := fmt.Sprintf("%s_%s", domain.Prefix, db.Name)
		dbPassword, _ := defaultCrypto.Decode(databaseServer.Password)
		if db.Type == Environments.MSSQL {
			//Env.Logger.Debug(9999, "Checking Microsoft SQL for: %s:%s", domain.Name, db.Name)
			dbSize := calculateMicrosoftSQL(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
			totalMicrosoftSqlSize += dbSize
		} else if db.Type == Environments.MYSQL {
			//Env.Logger.Debug(9999, "Checking MySQL for: %s:%s", domain.Name, db.Name)
			dbSize := calculateMySQL(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
			totalMySqlSize += dbSize
		} else if db.Type == Environments.PGSQL {
			//Env.Logger.Debug(9999, "Checking Postgre SQL for: %s:%s", domain.Name, db.Name)
			dbSize := calculatePostgreSQL(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
			totalPostgreSqlSize += dbSize
		}
	}
	return totalMicrosoftSqlSize + totalMySqlSize + totalPostgreSqlSize
}

func calculateMicrosoftSQL(username string, password string, host string, port int, dbName string) int64 {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", host, username, password, port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	stmt, err := conn.Prepare(fmt.Sprintf("SELECT " +
		"database_name = DB_NAME(database_id), " +
		"log_size = CAST(SUM(CASE WHEN type_desc = 'LOG' THEN size END) * 8 AS DECIMAL(8,2)), " +
		"row_size = CAST(SUM(CASE WHEN type_desc = 'ROWS' THEN size END) * 8 AS DECIMAL(8,2)), " +
		"total_size = CAST(SUM(size) * 8 AS DECIMAL(8,2)) " +
		"FROM sys.master_files WITH(NOWAIT) " +
		"WHERE database_id = DB_ID('%s') " +
		"GROUP BY database_id;", dbName))
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(stmt)

	row := stmt.QueryRow()
	var databaseName string
	var logSize int64
	var rowSize int64
	var totalSize int64
	err = row.Scan(&databaseName, &logSize, &rowSize, &totalSize)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
	return totalSize
}

func calculateMySQL(username string, password string, host string, port int, dbName string) int64 {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	query := fmt.Sprintf("SELECT " +
		"table_schema \"dbname\", " +
		"SUM(data_length + index_length) \"size\" " +
		"FROM information_schema.tables where table_schema = '%s' GROUP BY table_schema;", dbName)


	rows, err := conn.Query(query)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
		}
	}(rows)

	var databaseName string
	var totalSize int64 = 0
	for rows.Next() {
		err = rows.Scan(&databaseName, &totalSize)
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
		}
	}


	//Env.Logger.Debug(9999, "Total MySQL DB Size for:%s => %d", dbName, totalSize)
	return totalSize
}

func calculatePostgreSQL(username string, password string, host string, port int, dbName string) int64 {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, username, password)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	query := fmt.Sprintf("SELECT " +
		"t1.datname AS dbname, " +
		"pg_database_size(t1.datname) as size " +
		"from pg_database t1 where t1.datname = '%s' " +
		"order by pg_database_size(t1.datname) desc;", dbName)

	row := conn.QueryRow(query)
	var databaseName string
	var totalSize int64
	err = row.Scan(&databaseName, &totalSize)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
	//Env.Logger.Debug(9999, "Total PostgreSQL DB Size for:%s => %d", dbName, totalSize)
	return totalSize
}

func lockAllDatabases(domain Gorm.Domains)  {
	for _, db := range domain.Databases {
		var databaseServer Gorm.DatabaseServer
		res := DbConn.DB.Find(&databaseServer, "id = ?", db.DatabaseServerId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Database server not found or has an error: %v", res.Error))
			continue
		}
		defaultCrypto, _ := Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
		dbName := fmt.Sprintf("%s_%s", domain.Prefix, db.Name)
		dbPassword, _ := defaultCrypto.Decode(databaseServer.Password)
		if db.Type == Environments.MYSQL {
			//Env.Logger.Debug(9999, "Locking MySQL for: %s:%s", domain.Name, db.Name)
			lockAllMySQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
		} else if db.Type == Environments.MSSQL {
			//Env.Logger.Debug(9999, "Locking Microsoft SQL for: %s:%s", domain.Name, db.Name)
			lockAllMicrosoftSQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
		} else if db.Type == Environments.PGSQL {
			//Env.Logger.Debug(9999, "Locking Postgre SQL for: %s:%s", domain.Name, db.Name)
			lockAllPostgreSQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName, db.DatabaseUsers.Login)
		}
	}
}

func unlockAllDatabases(domain Gorm.Domains)  {
	for _, db := range domain.Databases {
		var databaseServer Gorm.DatabaseServer
		res := DbConn.DB.Find(&databaseServer, "id = ?", db.DatabaseServerId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Database server not found or has an error: %v", res.Error))
			continue
		}
		defaultCrypto, _ := Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
		dbName := fmt.Sprintf("%s_%s", domain.Prefix, db.Name)
		dbPassword, _ := defaultCrypto.Decode(databaseServer.Password)
		if db.Type == Environments.MYSQL {
			//Env.Logger.Debug(9999, "UnLocking MySQL for: %s:%s", domain.Name, db.Name)
			unlockAllMySQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName)
		} else if db.Type == Environments.MSSQL {
			//Env.Logger.Debug(9999, "UnLocking Microsoft SQL for: %s:%s", domain.Name, db.Name)
			unlockAllMicrosoftSQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName, domain.DomainStatistics.DbMax)
		} else if db.Type == Environments.PGSQL {
			//Env.Logger.Debug(9999, "UnLocking Postgre SQL for: %s:%s", domain.Name, db.Name)
			unlockAllPostgreSQLPrivileges(databaseServer.Username, dbPassword, databaseServer.Host, databaseServer.Port, dbName, db.DatabaseUsers.Login)
		}
	}
}

func lockAllMySQLPrivileges(username string, password string, host string, port int, dbName string)  {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	query := fmt.Sprintf("use mysql; UPDATE db SET Insert_priv='N', Create_priv='N', Update_priv='N' WHERE Db='%s'; FLUSH PRIVILEGES;", dbName)
	_, err = conn.Exec(query)

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL locking table failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

func lockAllPostgreSQLPrivileges(username string, password string, host string, port int, dbName string, dbUsername string)  {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, username, password)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	query := fmt.Sprintf("use %s; SELECT " +
		"(table_schema || '.' || table_name) as tnames " +
		"FROM " +
		"information_schema.tables " +
		"WHERE " +
		"table_schema NOT IN ('pg_catalog', 'information_schema');", dbName)

	row := conn.QueryRow(query)
	var tnames []string
	err = row.Scan(&tnames)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
	tableNames := strings.Join(tnames, ",")

	queryUpdate := fmt.Sprintf("use %[1]s; " +
		"REVOKE INSERT ON %[2]s FROM %[1]s; "+
		"REVOKE UPDATE ON %[2]s FROM %[1]s; "+
		"GRANT DELETE ON %[2]s TO  %[1]s; ", dbUsername, tableNames)

	_, err = conn.Exec(queryUpdate)

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL locking table failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

func lockAllMicrosoftSQLPrivileges(username string, password string, host string, port int, dbName string)  {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", host, username, password, port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	stmt, err := conn.Prepare(fmt.Sprintf("SELECT " +
		"database_name = DB_NAME(database_id), " +
		"log_size = CAST(SUM(CASE WHEN type_desc = 'LOG' THEN size END) * 8 AS DECIMAL(8,2)), " +
		"row_size = CAST(SUM(CASE WHEN type_desc = 'ROWS' THEN size END) * 8 AS DECIMAL(8,2)), " +
		"total_size = CAST(SUM(size) * 8 AS DECIMAL(8,2)) " +
		"FROM sys.master_files WITH(NOWAIT) " +
		"WHERE database_id = DB_ID('%s') " +
		"GROUP BY database_id;", dbName))

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL sql query failed for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(stmt)

	row := stmt.QueryRow()
	var totalSize int64
	err = row.Scan(&totalSize)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}

	_, err = conn.Exec(fmt.Sprintf("USE master; ALTER DATABASE [%[1]s] MODIFY FILE ( NAME = N'%[1]s', MAXSIZE = %dKB );", dbName, totalSize))

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

func unlockAllMySQLPrivileges(username string, password string, host string, port int, dbName string)  {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	query := fmt.Sprintf("use mysql; UPDATE db SET Insert_priv='Y', Create_priv='Y', Update_priv='Y' WHERE Db='%s'; FLUSH PRIVILEGES;", dbName)

	_, err = conn.Exec(query)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MySQL unlocking table failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

func unlockAllPostgreSQLPrivileges(username string, password string, host string, port int, dbName string, dbUsername string)  {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, username, password)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	query := fmt.Sprintf("use %s; SELECT " +
		"(table_schema || '.' || table_name) as tnames " +
		"FROM " +
		"information_schema.tables " +
		"WHERE " +
		"table_schema NOT IN ('pg_catalog', 'information_schema');", dbName)

	row := conn.QueryRow(query)
	var tnames []string
	err = row.Scan(&tnames)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL getting result failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
	tableNames := strings.Join(tnames, ",")

	queryUpdate := fmt.Sprintf("use %[1]s; " +
		"GRANT INSERT ON %[2]s FROM %[1]s; "+
		"GRANT UPDATE ON %[2]s FROM %[1]s; "+
		"GRANT DELETE ON %[2]s TO  %[1]s; ", dbUsername, tableNames)

	_, err = conn.Exec(queryUpdate)

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("PostgreSQL unlocking table failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

func unlockAllMicrosoftSQLPrivileges(username string, password string, host string, port int, dbName string, totalSize string)  {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", host, username, password, port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL connection has an error for domain: %s:%d Database Name: %s", host, port, dbName))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL connection close has an error for domain: %s:%d Database Name: %s", host, port, dbName))
		}
	}(conn)

	var dTotalSize int64
	if strings.Compare(totalSize, Environments.UNLIMITED) == 0 {
		dTotalSize = 1024*1024*1024
	} else {
		tTotalSize, _ := strconv.Atoi(totalSize)
		dTotalSize = int64(tTotalSize)
	}
	_, err = conn.Exec(fmt.Sprintf("USE master; ALTER DATABASE [%[1]s] MODIFY FILE ( NAME = N'%[1]s', MAXSIZE = %dKB );", dbName, dTotalSize/1024))

	if err != nil {
		Helpers.HandleError(err, 9001, fmt.Sprintf("MicrosoftSQL exec query failed for host: %s:%d Database Name: %s", host, port, dbName))
	}
}

