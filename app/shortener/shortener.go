package shortener

import (
	r "github.com/robfig/revel"
	"go-url-shortener/app/db"
	m "go-url-shortener/app/models"
	"go-url-shortener/app/store"
)

func Init() {
	urlGetter := func(slug string) (interface{}, error) {
		session := db.Engine.NewSession()
		defer session.Close()
		s, err := m.ShortUrlBySlug(session, slug)
		if err != nil || s == nil {
			return nil, err
		}
		return &s.URL, nil
	}

	store.NewStore(m.StoreName, "", store.GetterFunc(urlGetter))
}

func init() {
	r.OnAppStart(Init)
}
