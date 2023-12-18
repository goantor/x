package x

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// IConfig todo 暂时 简单封装，后续业务需要时扩展
type IConfig interface {
	FlagFile(name string, value string, usage string)
	AddFlagSet(name string, newSet *flag.FlagSet)
	GivePath(path string) IConfig
	TakePath() string

	Parse() (err error)
	Get(key string) any
	Bind(configs any) (err error)
}

type defaultConfig struct {
	set  *flag.FlagSet
	name string
	path string
}

func NewConfig() IConfig {
	return &defaultConfig{
		set: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
}

func (c *defaultConfig) FlagFile(name string, value string, usage string) {
	c.set.String(name, value, usage)
	c.name = name
}

func (c *defaultConfig) parse() error {
	pflag.CommandLine.AddGoFlagSet(c.set)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}

	return c.readFile()
}

func (c *defaultConfig) Parse() (err error) {
	if !c.checkPath() {
		return c.parse()
	}

	return c.readFile()
}

func (c *defaultConfig) checkPath() bool {
	return c.path != ""
}

func (c *defaultConfig) makeConfigPath() (err error) {
	var (
		ok bool
	)

	if ok = c.checkPath(); ok {
		return
	}

	fileAny := viper.Get(c.name)
	if c.path, ok = fileAny.(string); ok {
		return
	}

	return fmt.Errorf("%s must be a string var", c.name)
}

func (c *defaultConfig) readFile() (err error) {
	fmt.Println("do read")
	if err = c.makeConfigPath(); err != nil {
		return
	}

	fmt.Printf("path: %s\n", c.path)
	viper.SetConfigFile(c.path)
	return viper.ReadInConfig()
}

func (c *defaultConfig) AddFlagSet(name string, set *flag.FlagSet) {
	c.name, c.set = name, set
}

func (c *defaultConfig) GivePath(path string) IConfig {
	c.path = path
	return c
}

func (c *defaultConfig) TakePath() string {
	return c.path
}

func (c *defaultConfig) Get(key string) any {
	return viper.Get(key)
}

func (c *defaultConfig) Bind(configs any) (err error) {

	return viper.Unmarshal(&configs)
}
