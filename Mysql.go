package x

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Database interface {
	Link() string
	Initialize() error
	Close() error
}

var (
	MysqlDB *gorm.DB
)

type Mysql struct {
	Debug       bool   `json:"debug" toml:"debug"`
	Host        string `json:"host" toml:"host"`
	Port        int    `json:"port" toml:"port"`
	User        string `json:"user" toml:"user"`
	Password    string `json:"password" toml:"password"`
	Database    string `json:"database" toml:"database"`
	Charset     string `json:"charset" toml:"charset"`
	MaxIdleConn int    `json:"max_idle_conn" toml:"max_idle_conn"`
	MaxOpenConn int    `json:"max_open_conn" toml:"max_open_conn"`
	MaxLifeTime int    `json:"max_life_time" toml:"max_life_time"`
}

func NewMysql() *Mysql {
	return &Mysql{}
}

func (m Mysql) Link() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func (m Mysql) Init() {
	var err error

	MysqlDB, err = gorm.Open(mysql.Open(m.Link()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic(fmt.Sprintf("mysql connect failed: %+v\n", err))
	}

	if m.Debug {
		MysqlDB.Debug()
	}

	pool, err := MysqlDB.DB()
	if err != nil {
		panic(fmt.Sprintf("mysql get pool failed: %+v\n", err))
	}

	pool.SetMaxIdleConns(m.MaxIdleConn)
	pool.SetMaxOpenConns(m.MaxOpenConn)
	pool.SetConnMaxLifetime(time.Duration(m.MaxLifeTime) * time.Second)
}

func (m Mysql) Close() error {
	return nil
}
