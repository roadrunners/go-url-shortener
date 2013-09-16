package models

import (
	"fmt"
	"github.com/roadrunners/go-url-shortener/app/db"
	"github.com/roadrunners/go-url-shortener/app/models/key"
	"github.com/roadrunners/go-url-shortener/app/redis"
	"github.com/robfig/revel"
)

type ShortUrl struct {
	Id   int64  `db:"id" json:"-"`
	Slug string `db:"-" json:"slug"`
	URL  string `db:"url" json:"url"`
}

func (s ShortUrl) String() string {
	return fmt.Sprintf("(%s, %s)", s.Slug, s.URL)
}

func (s *ShortUrl) pushToRedis() {
	k := fmt.Sprintf("shorturl:%s:url", s.Slug)
	err := redis.Client.Set(k, []byte(s.URL))
	if err != nil {
		revel.ERROR.Fatal("Could not push short url to redis")
	}
}

func ShortUrlById(id int64) (*ShortUrl, error) {
	v, err := db.DbMap.Get(ShortUrl{}, id)
	if err != nil || v == nil {
		return nil, err
	}
	s := v.(*ShortUrl)
	s.Slug = key.GenKey(s.Id)
	return s, nil
}

func ShortUrlBySlug(slug string) (*ShortUrl, error) {
	id := key.GenId(slug)
	return ShortUrlById(id)
}

func CachedShortUrlBySlug(slug string) (*ShortUrl, error) {
	k := fmt.Sprintf("shorturl:%s:url", slug)
	data, err := redis.Client.Get(k)
	if err == nil {
		s := ShortUrl{Slug: slug, URL: string(data)}
		return &s, nil
	}
	id := key.GenId(slug)
	revel.WARN.Printf("Missed cache for slug %v (id %v, key %v)", slug, id, k)
	s, err := ShortUrlById(id)
	if s != nil && err == nil {
		go s.pushToRedis()
	}
	return s, err
}

func ShortUrlCreate(url string) (*ShortUrl, error) {
	s := &ShortUrl{URL: url}
	if err := db.DbMap.Insert(s); err != nil {
		return nil, err
	}
	s.Slug = key.GenKey(s.Id)
	go s.pushToRedis()
	return s, nil
}
