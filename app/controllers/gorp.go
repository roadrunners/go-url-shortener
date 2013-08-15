package controllers

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	r "github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
	. "go-url-shortener/app/db"
	"go-url-shortener/app/models"
)

func Init() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{}}

	t := Dbm.AddTableWithName(models.ShortURL{}, "short_urls").SetKeys(false, "Slug")
	setColumnSizes(t, map[string]int{
		"Slug": 20,
		"URL":  512,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
}

func setColumnSizes(t *gorp.TableMap, colSizes map[string]int) {
	for col, size := range colSizes {
		t.ColMap(col).MaxSize = size
	}
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
