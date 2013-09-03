package db

import (
	"fmt"
	"github.com/lunny/xorm"
	r "github.com/robfig/revel"
)

var (
	Engine *xorm.Engine
)

func ensureOption(option string) string {
	value, found := r.Config.String(option)
	if !found {
		r.ERROR.Fatalf("Option %v not found", option)
	}
	return value
}

func calcSpec() string {
	host := ensureOption("db.host")
	name := ensureOption("db.name")
	username := ensureOption("db.username")
	password := r.Config.StringDefault("db.password", "")
	return fmt.Sprintf("%v:%v@tcp(%v)/%v", username, password, host, name)
}

func Init() {
	var found bool
	var driver string
	if driver, found = r.Config.String("db.driver"); !found {
		r.ERROR.Fatal("No db.driver found.")
	}
	var spec string
	if spec, found = r.Config.String("db.spec"); !found {
		spec = calcSpec()
	}
	r.INFO.Printf("Connecting to mysql at %v", spec)
	engine, err := xorm.NewEngine(driver, spec)
	if err != nil {
		panic(err)
	}
	Engine = engine
}

func init() {
	r.OnAppStart(Init)
}
