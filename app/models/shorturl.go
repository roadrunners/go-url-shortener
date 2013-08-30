package models

import (
	"fmt"
	"github.com/lunny/xorm"
	r "github.com/robfig/revel"
	"go-url-shortener/app/redis"
)

type ShortUrl struct {
	Id   int64  `xorm:"id pk not null autoincr" json:"-"`
	Slug string `xorm:"-" json:"slug"`
	URL  string `xorm:"url unique" json:"url"`
}

func (s ShortUrl) String() string {
	return fmt.Sprintf("(%s, %s)", s.Slug, s.URL)
}

func (s *ShortUrl) pull() {
	key := fmt.Sprintf("shorturl:%d:url", s.Id)
	err := redis.Client.Set(key, []byte(s.URL))
	if err != nil {
		r.ERROR.Fatal("Could not push short url to redis")
	}
}

func ShortUrlById(session *xorm.Session, id int64) (*ShortUrl, error) {
	var s ShortUrl
	has, err := session.Id(id).Get(&s)
	if err != nil || !has {
		return nil, err
	}
	s.Slug = genKey(s.Id)
	return &s, nil
}

func ShortUrlBySlug(session *xorm.Session, slug string) (*ShortUrl, error) {
	id := genId(slug)
	return ShortUrlById(session, id)
}

func CachedShortUrlBySlug(session *xorm.Session, slug string) (*ShortUrl, error) {
	id := genId(slug)
	key := fmt.Sprintf("shorturl:%d:url", id)
	data, err := redis.Client.Get(key)
	if err == nil {
		s := ShortUrl{Id: id, Slug: slug, URL: string(data)}
		return &s, nil
	}
	s, err := ShortUrlById(session, id)
	if s != nil {
		go s.pull()
	}
	return s, err
}

func ShortUrlCreate(session *xorm.Session, url string) (*ShortUrl, error) {
	s := ShortUrl{URL: url}
	_, err := session.Insert(&s)
	if err != nil {
		return nil, err
	}
	s.Slug = genKey(s.Id)
	go s.pull()
	return &s, nil
}

var (
	keyChar   = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	decodeMap = makeDecodeMap()
)

func makeDecodeMap() map[byte]int64 {
	m := make(map[byte]int64)
	for i, b := range keyChar {
		m[b] = int64(i)
	}
	return m
}

func genKey(n int64) string {
	if n == 0 {
		return string(keyChar[0])
	}
	l := int64(len(keyChar))
	s := make([]byte, 20)
	i := int64(len(s))
	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = keyChar[j]
	}
	return string(s[i:])
}

func genId(key string) int64 {
	l := int64(len(keyChar))
	n := int64(0)
	for _, b := range key {
		n *= l
		n += decodeMap[byte(b)]
	}
	return n
}
