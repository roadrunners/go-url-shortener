package controllers

import (
	r "github.com/robfig/revel"
)

func init() {
	r.InterceptMethod((*XormController).Before, r.BEFORE)
	r.InterceptMethod((*XormController).After, r.AFTER)
	r.InterceptMethod((*XormController).After, r.FINALLY)
}
