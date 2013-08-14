package models

import (
	"fmt"
	"go-url-shortener/app/db"
)

type ShortURL struct {
	Id   int    `db:"id"`
	Slug string `db:"slug"`
	URL  string `db:"url"`
}

func (s ShortURL) String() string {
	return fmt.Sprintf("(%s, %s)", s.Slug, s.URL)
}

func ShortUrlBySlug(slug string) (s *ShortURL, err error) {
	var list []*ShortURL
	_, err = db.Dbm.Select(&list, "select id, slug, url from short_urls where slug = ?", slug)
	if err != nil {
		return
	}
	if len(list) < 1 {
		return
	}
	return list[0], nil
}

func ShortURLCreate(slug, url string) (s ShortURL, err error) {
	shortURL := ShortURL{Slug: slug, URL: url}
	if err = db.Dbm.Insert(&shortURL); err != nil {
		return
	}

	s = shortURL
	return
}

func ShortURLCount() (count int64, err error) {
	count, err = db.Dbm.SelectInt("select count(*) from short_urls")
	return
}
