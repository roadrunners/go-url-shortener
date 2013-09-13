package controllers

import (
	"github.com/lunny/xorm"
	"github.com/roadrunners/go-url-shortener/app/db"
	r "github.com/robfig/revel"
)

type XormController struct {
	*r.Controller
	XormSession *xorm.Session
}

func (c *XormController) Before() r.Result {
	c.XormSession = db.Engine.NewSession()
	return nil
}

func (c *XormController) After() r.Result {
	if c.XormSession == nil {
		return nil
	}
	c.XormSession.Close()
	c.XormSession = nil
	return nil
}
