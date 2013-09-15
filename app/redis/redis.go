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
	poolSize := revel.Config.IntDefault("redis.poolsize", 0)
	revel.INFO.Printf("Connecting to redis db %v at %v with poolsize %v (0 = use driver default)", db, addr, poolSize)
	Client = &redis.Client{Addr: addr, Db: db, MaxPoolSize: poolSize}
}

func init() {
	revel.OnAppStart(Init)
}
