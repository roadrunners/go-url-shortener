package app

import (
	r "github.com/robfig/revel"
	"runtime"
)

func Init() {
	numCPU := runtime.NumCPU()
	gomaxprocs := runtime.GOMAXPROCS(numCPU)

	r.WARN.Printf("GOMAXPROCS was %v, Total CPU detected %v, set as new GOMAXPROCS", gomaxprocs, numCPU)
}

func init() {
	r.Filters = []r.Filter{
		r.PanicFilter,             // Recover from panics and display an error page instead.
		r.RouterFilter,            // Use the routing table to select the right Action
		r.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		r.ParamsFilter,            // Parse parameters into Controller.Params.
		r.SessionFilter,           // Restore and write the session cookie.
		r.FlashFilter,             // Restore and write the flash cookie.
		r.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		r.I18nFilter,              // Resolve the requested language
		r.InterceptorFilter,       // Run interceptors around the action.
		r.ActionInvoker,           // Invoke the action.
	}

	r.OnAppStart(Init)
}
