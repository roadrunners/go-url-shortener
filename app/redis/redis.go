package redis

import (
	"github.com/hoisie/redis"
	"github.com/robfig/revel"
)

var (
	Client *redis.Client
)

func Init() {
	var found bool
	var addr string
	if addr, found = revel.Config.String("redis.addr"); !found {
		revel.ERROR.Fatal("No redis.addr found")
	}
	var db int
	if db, found = revel.Config.Int("redis.db"); !found {
		revel.ERROR.Fatal("No redis.db found")
	}
	revel.INFO.Printf("Connecting to redis db %v at %v", db, addr)
	Client = &redis.Client{Addr: addr, Db: db}
}

func init() {
	revel.OnAppStart(Init)
}
