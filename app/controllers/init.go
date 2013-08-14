package controllers

import (
	r "github.com/robfig/revel"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r.OnAppStart(Init)
	r.InterceptMethod((*GorpController).Begin, r.BEFORE)
	r.InterceptMethod((*GorpController).Commit, r.AFTER)
	r.InterceptMethod((*GorpController).Rollback, r.FINALLY)
}
