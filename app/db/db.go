package db

import (
	"github.com/lunny/xorm"
	r "github.com/robfig/revel"
)

var (
	Engine *xorm.Engine
)

func Init() {
	var found bool
	var driver string
	if driver, found = r.Config.String("db.driver"); !found {
		r.ERROR.Fatal("No db.driver found.")
	}
	var spec string
	if spec, found = r.Config.String("db.spec"); !found {
		r.ERROR.Fatal("No db.spec found.")
	}
	engine, err := xorm.NewEngine(driver, spec)
	if err != nil {
		panic(err)
	}
	Engine = engine
}

func init() {
	r.OnAppStart(Init)
}
