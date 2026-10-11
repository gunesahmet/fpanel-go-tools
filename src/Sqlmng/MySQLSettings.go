package Sqlmng

import (
	"Helpers"
	"github.com/go-ini/ini"
	"path/filepath"
	"strings"
)

type MySQLSettings struct {
	ini *ini.File
}

func NewMySQLSettings(config string) (*MySQLSettings, error) {
	cfg, err := ini.Load(config)
	if err != nil {
		return nil, err
	}

	return &MySQLSettings{cfg}, nil
}
func (c *MySQLSettings) MySQLRemoteConnection(value string) error {
	if strings.Compare(value, "on") == 0 {
		c.ini.Section("mysqld").Key("bind-address").SetValue("0.0.0.0")
	} else {
		c.ini.Section("mysqld").Key("bind-address").SetValue("127.0.0.1")
	}
	return nil
}
func (c *MySQLSettings) MySQLBufferPoolSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("innodb_buffer_pool_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("innodb_buffer_pool_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLLogFileSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("innodb_log_file_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("innodb_log_file_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLFilePerTable(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("innodb_file_per_table").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("innodb_file_per_table")
	}
	return nil
}
func (c *MySQLSettings) MySQLFlushMethod(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("innodb_flush_method").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("innodb_flush_method")
	}
	return nil
}
func (c *MySQLSettings) MySQLLogError(value string) error {
	if strings.Compare(value, "default") != 0 {
		dir := filepath.Dir(value)
		if !Helpers.DirExists(dir) {
			Helpers.CreateDirectory(dir, true)
		_:
			Helpers.Exec("chown", "-R", "mysql:mysql", dir)
		}
		c.ini.Section("mysqld").Key("log-error").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("log-error")
	}
	return nil
}
func (c *MySQLSettings) MySQLSlowQueryLog(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("slow_query_log").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("slow_query_log")
	}
	return nil
}
func (c *MySQLSettings) MySQLSlowQueryLogFile(value string) error {
	if !Helpers.DirExists(value) {
		Helpers.CreateDirectory(value, true)
	}
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("slow_query_log_file").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("slow_query_log_file")
	}
	return nil
}
func (c *MySQLSettings) MySQLTmpTableSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("tmp_table_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("tmp_table_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLMaxHeapTableSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("max_heap_table_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("max_heap_table_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLQueryCacheType(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("query_cache_type").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("query_cache_type")
	}
	return nil
}
func (c *MySQLSettings) MySQLQueryCacheSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("query_cache_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("query_cache_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLMaxConnections(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("max_connections").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("max_connections")
	}
	return nil
}
func (c *MySQLSettings) MySQLThreadCacheSize(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("thread_cache_size").SetValue(value + "M")
	} else {
		c.ini.Section("mysqld").DeleteKey("thread_cache_size")
	}
	return nil
}
func (c *MySQLSettings) MySQLTableOpenCache(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("table_open_cache").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("table_open_cache")
	}
	return nil
}
func (c *MySQLSettings) MySQLOpenFilesLimit(value string) error {
	if strings.Compare(value, "default") != 0 {
		c.ini.Section("mysqld").Key("open_files_limit").SetValue(value)
	} else {
		c.ini.Section("mysqld").DeleteKey("open_files_limit")
	}
	return nil
}
