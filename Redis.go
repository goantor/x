package x

import (
	"fmt"
	"github.com/go-redis/redis"
)

var (
	RedisDB *redis.Client
)

type Redis struct {
	Host     string `json:"host" toml:"host"`
	Port     int    `json:"port" toml:"port"`
	Password string `json:"password" toml:"password"`
	Database int    `json:"database" toml:"database"`
}

func (r *Redis) Link() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *Redis) Init() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     r.Link(),
		Password: r.Password,
		DB:       int(r.Database),
	})

	if _, err := RedisDB.Ping().Result(); err != nil {
		panic(fmt.Sprintf("redis connect failed: %+v\n", err))
	}
}

func (r *Redis) Close() error {
	if RedisDB == nil {
		return nil
	}
	return RedisDB.Close()
}
