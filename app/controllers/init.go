package controllers

import (
	r "github.com/robfig/revel"
	"runtime"
)

func init() {
	numCPU := runtime.NumCPU()
	gomaxprocs := runtime.GOMAXPROCS(numCPU)

	r.WARN.Printf("Total CPU detected %v, setting GOMAXPROCS to %v", numCPU, gomaxprocs)

	r.InterceptMethod((*XormController).Begin, r.BEFORE)
	r.InterceptMethod((*XormController).Commit, r.AFTER)
	r.InterceptMethod((*XormController).Rollback, r.FINALLY)
}
