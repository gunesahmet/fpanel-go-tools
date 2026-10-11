package Statisticsmng

import (
	"github.com/go-ini/ini"
	"strings"
)

type SystemConfig struct {
	ini *ini.File
}

func NewSystemConfig(config string) (*SystemConfig, error) {
	cfg, err := ini.Load(config)
	if err != nil {
		return nil, err
	}
	return &SystemConfig{cfg}, nil
}
func (c *SystemConfig) MainServerAddress(value string) error {
	c.ini.Section("server").Key("address").SetValue(value)
	return nil
}
func (c *SystemConfig) StatisticsRandom(value string) error {
	if strings.Compare(value, "on") != 0 {
		runningTime := "00:00"
		if c.ini.Section("statistics").HasKey("time") &&
			c.ini.Section("statistics").HasValue("time") &&
			len(c.ini.Section("statistics").Key("time").String()) > 0 {
			runningTime = c.ini.Section("statistics").Key("time").String()
		}
		c.ini.Section("statistics").Key("random").SetValue("0")
		c.ini.Section("statistics").Key("time").SetValue(runningTime)
	} else {
		c.ini.Section("statistics").Key("random").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) StatisticsTime(value string) error {
	c.ini.Section("statistics").Key("time").SetValue(value)
	return nil
}
func (c *SystemConfig) UpdateRandom(value string) error {
	if strings.Compare(value, "on") != 0 {
		runningTime := "00:00"
		if c.ini.Section("update").HasKey("time") &&
			c.ini.Section("update").HasValue("time") &&
			len(c.ini.Section("update").Key("time").String()) > 0 {
			runningTime = c.ini.Section("update").Key("time").String()
		}
		c.ini.Section("update").Key("random").SetValue("0")
		c.ini.Section("update").Key("time").SetValue(runningTime)
	} else {
		c.ini.Section("update").Key("random").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) UpdateAuto(value string) error {
	if strings.Compare(value, "on") != 0 {
		c.ini.Section("update").Key("auto").SetValue("0")
	} else {
		c.ini.Section("update").Key("auto").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) UpdateTime(value string) error {
	c.ini.Section("update").Key("time").SetValue(value)
	return nil
}
