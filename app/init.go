package app

import (
	"github.com/robfig/revel"
	"runtime"
)

func Init() {
	numCPU := runtime.NumCPU()
	gomaxprocs := runtime.GOMAXPROCS(numCPU)

	revel.WARN.Printf("GOMAXPROCS was %v, Total CPU detected %v, set as new GOMAXPROCS", gomaxprocs, numCPU)
}

func init() {
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(Init)
}
