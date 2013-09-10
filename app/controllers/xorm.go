package controllers

import (
	"database/sql"
	"github.com/lunny/xorm"
	db "github.com/roadrunners/go-url-shortener/app/db"
	r "github.com/robfig/revel"
)

type XormController struct {
	*r.Controller
	XormSession *xorm.Session
}

func (c *XormController) Begin() r.Result {
	session := db.Engine.NewSession()
	err := session.Begin()
	if err != nil {
		panic(err)
	}
	c.XormSession = session
	return nil
}

func (c *XormController) Commit() r.Result {
	if c.XormSession == nil {
		return nil
	}
	if err := c.XormSession.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.XormSession.Close()
	c.XormSession = nil
	return nil
}

func (c *XormController) Rollback() r.Result {
	if c.XormSession == nil {
		return nil
	}
	if err := c.XormSession.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.XormSession.Close()
	c.XormSession = nil
	return nil
}
