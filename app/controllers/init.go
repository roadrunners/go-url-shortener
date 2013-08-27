package controllers

import (
	r "github.com/robfig/revel"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r.InterceptMethod((*XormController).Begin, r.BEFORE)
	r.InterceptMethod((*XormController).Commit, r.AFTER)
	r.InterceptMethod((*XormController).Rollback, r.FINALLY)
}
