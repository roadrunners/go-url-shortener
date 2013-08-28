package models

import (
	"fmt"
	"github.com/lunny/xorm"
	"go-url-shortener/app/store"
)

const (
	StoreName = "ShortUrlStore"
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
	store := getStore()
	store.Pull(s.Slug)
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
	store := getStore()
	url, err := store.Get(slug)
	if err != nil || url == nil {
		return nil, err
	}
	s := ShortUrl{Slug: slug, URL: *url.(*string)}
	s.Id = genId(s.Slug)
	return &s, nil
}

func ShortUrlCreate(session *xorm.Session, url string) (*ShortUrl, error) {
	s := ShortUrl{URL: url}
	_, err := session.Insert(&s)
	if err != nil {
		return nil, err
	}
	s.Slug = genKey(s.Id)
	s.pull()
	return &s, nil
}

var (
	urlStore *store.Store
)

func getStore() *store.Store {
	return store.GetStore(StoreName)
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
