package SystemConf

import (
	"github.com/go-ini/ini"
	"strings"
)

type SystemConfig struct {
	Ini *ini.File
}

func NewSystemConfig(config string) (*SystemConfig, error) {
	cfg, err := ini.Load(config)
	if err != nil {
		return nil, err
	}
	return &SystemConfig{cfg}, nil
}
func (c *SystemConfig) MainServerAddress(value string) error {
	c.Ini.Section("server").Key("address").SetValue(value)
	return nil
}
func (c *SystemConfig) StatisticsRandom(value string) error {
	if strings.Compare(value, "on") != 0 {
		runningTime := "00:00"
		if c.Ini.Section("statistics").HasKey("time") &&
			c.Ini.Section("statistics").HasValue("time") &&
			len(c.Ini.Section("statistics").Key("time").String()) > 0 {
			runningTime = c.Ini.Section("statistics").Key("time").String()
		}
		c.Ini.Section("statistics").Key("random").SetValue("0")
		c.Ini.Section("statistics").Key("time").SetValue(runningTime)
	} else {
		c.Ini.Section("statistics").Key("random").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) StatisticsTime(value string) error {
	c.Ini.Section("statistics").Key("time").SetValue(value)
	return nil
}
func (c *SystemConfig) UpdateRandom(value string) error {
	if strings.Compare(value, "on") != 0 {
		runningTime := "00:00"
		if c.Ini.Section("update").HasKey("time") &&
			c.Ini.Section("update").HasValue("time") &&
			len(c.Ini.Section("update").Key("time").String()) > 0 {
			runningTime = c.Ini.Section("update").Key("time").String()
		}
		c.Ini.Section("update").Key("random").SetValue("0")
		c.Ini.Section("update").Key("time").SetValue(runningTime)
	} else {
		c.Ini.Section("update").Key("random").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) UpdateAuto(value string) error {
	if strings.Compare(value, "on") != 0 {
		c.Ini.Section("update").Key("auto").SetValue("0")
	} else {
		c.Ini.Section("update").Key("auto").SetValue("1")
	}
	return nil
}
func (c *SystemConfig) UpdateTime(value string) error {
	c.Ini.Section("update").Key("time").SetValue(value)
	return nil
}
