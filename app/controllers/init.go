package controllers

import (
	"github.com/robfig/revel"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	revel.OnAppStart(Init)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
