package models

import (
	"github.com/coopernurse/gorp"
	"github.com/roadrunners/go-url-shortener/app/db"
	"github.com/robfig/revel"
)

func Init() {
	t := db.DbMap.AddTableWithName(ShortUrl{}, "short_url").SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"URL": 512,
	})

	db.DbMap.TraceOn("[gorp]", revel.INFO)
}

func setColumnSizes(t *gorp.TableMap, colSizes map[string]int) {
	for col, size := range colSizes {
		t.ColMap(col).MaxSize = size
	}
}

func init() {
	revel.OnAppStart(Init)
}
