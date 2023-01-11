package x

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	Configs = &DefaultConfigs{}
)

type DefaultConfigs struct {
	Http  HttpConfig `json:"http" toml:"http"`
	Log   LogConfig  `json:"log" toml:"log"`
	Mysql Mysql      `json:"mysql" toml:"mysql"`
	Redis Redis      `json:"redis" toml:"redis"`
}

func (c DefaultConfigs) Boot() {
	c.Mysql.Init()
}

type TomlConfig struct {
	Path string
}

func NewTomlConfig(dir string) *TomlConfig {
	path := fmt.Sprintf("%s/%s.toml", dir, Env)

	t := &TomlConfig{Path: path}
	if exists, err := t.fileExists(); !exists {
		if err == nil {
			panic("toml config file not exists")
		}

		panic(err)
	}

	return t
}

func (c *TomlConfig) Parse(config *DefaultConfigs) (err error) {
	_, err = toml.DecodeFile(c.Path, &config)
	return
}

func (c *TomlConfig) fileExists() (bool, error) {
	// 获取文件信息
	_, err := os.Stat(c.Path)
	if err == nil {
		return true, nil
	}

	// 如果出现错误, 判断错误是不是文件不存在
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
