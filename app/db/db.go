package db

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/robfig/revel"
)

var (
	DbMap *gorp.DbMap
)

func ensureOption(option string) string {
	value, found := revel.Config.String(option)
	if !found {
		revel.ERROR.Fatalf("Option %v not found", option)
	}
	return value
}

func calcSpec() string {
	host := ensureOption("db.host")
	name := ensureOption("db.name")
	username := ensureOption("db.username")
	password := revel.Config.StringDefault("db.password", "")
	return fmt.Sprintf("%v:%v@tcp(%v)/%v", username, password, host, name)
}

func Init() {
	var found bool
	var driver string
	if driver, found = revel.Config.String("db.driver"); !found {
		revel.ERROR.Fatal("No db.driver found.")
	}
	var spec string
	if spec, found = revel.Config.String("db.spec"); !found {
		spec = calcSpec()
	}
	revel.INFO.Printf("Connecting to mysql at %v", spec)
	db, err := sql.Open(driver, spec)
	if err != nil {
		panic(err)
	}
	DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}

func init() {
	revel.OnAppStart(Init)
}
