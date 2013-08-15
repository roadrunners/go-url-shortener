package shortener

import (
	r "github.com/robfig/revel"
	m "go-url-shortener/app/models"
	"go-url-shortener/app/store"
)

func Init() {
	shortURLGetter := func(slug string) (*string, error) {
		s, err := m.ShortUrlBySlug(slug)
		if err != nil || s == nil {
			return nil, err
		}
		return &s.URL, nil
	}

	store.NewStore(m.StoreName, store.GetterFunc(shortURLGetter))
}

func init() {
	r.OnAppStart(Init)
}
