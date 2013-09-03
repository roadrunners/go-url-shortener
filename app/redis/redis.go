package redis

import (
	"github.com/hoisie/redis"
	r "github.com/robfig/revel"
)

var (
	Client *redis.Client
)

func Init() {
	var found bool
	var addr string
	if addr, found = r.Config.String("redis.addr"); !found {
		r.ERROR.Fatal("No redis.addr found")
	}
	var db int
	if db, found = r.Config.Int("redis.db"); !found {
		r.ERROR.Fatal("No redis.db found")
	}
	r.INFO.Printf("Connecting to redis db %v at %v", db, addr)
	Client = &redis.Client{Addr: addr, Db: db}
}

func init() {
	r.OnAppStart(Init)
}
