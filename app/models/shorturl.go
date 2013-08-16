package models

import (
	"fmt"
	"go-url-shortener/app/db"
	"go-url-shortener/app/store"
)

const (
	StoreName = "URLStore"
)

type ShortURL struct {
	Id   int    `db:"id" json:"-"`
	Slug string `db:"-" json:"slug"`
	URL  string `db:"url" json:"url"`
}

func (s ShortURL) String() string {
	return fmt.Sprintf("(%s, %s)", s.Slug, s.URL)
}

func (s *ShortURL) pull() {
	store := GetStore()
	store.Pull(s.Slug)
}

func ShortUrlById(id int) (*ShortURL, error) {
	v, err := db.Dbm.Get(ShortURL{}, id)
	if err != nil || v == nil {
		return nil, err
	}
	s := v.(*ShortURL)
	s.Slug = genKey(s.Id)
	return s, nil
}

func ShortUrlBySlug(slug string) (*ShortURL, error) {
	id := genId(slug)
	return ShortUrlById(id)
}

func CachedShortUrlBySlug(slug string) (*ShortURL, error) {
	urlStore := GetStore()
	url, err := urlStore.Get(slug)
	if err != nil || url == nil {
		return nil, err
	}
	s := &ShortURL{Slug: slug, URL: *url}
	s.Id = genId(s.Slug)
	return s, nil
}

func ShortURLCreate(url string) (*ShortURL, error) {
	s := &ShortURL{URL: url}
	if err := db.Dbm.Insert(s); err != nil {
		return nil, err
	}

	s.Slug = genKey(s.Id)
	s.pull()
	return s, nil
}

var (
	urlStore *store.Store
)

func GetStore() *store.Store {
	if urlStore == nil {
		urlStore = store.GetStore(StoreName)
	}
	return urlStore
}

var (
	keyChar   = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	decodeMap = makeDecodeMap()
)

func makeDecodeMap() map[byte]int {
	m := make(map[byte]int)
	for i, b := range keyChar {
		m[b] = i
	}
	return m
}

func genKey(n int) string {
	if n == 0 {
		return string(keyChar[0])
	}
	l := len(keyChar)
	s := make([]byte, 20)
	i := len(s)
	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = keyChar[j]
	}
	return string(s[i:])
}

func genId(key string) int {
	l := len(keyChar)
	n := 0
	for _, b := range key {
		n *= l
		n += decodeMap[byte(b)]
	}
	return n
}
