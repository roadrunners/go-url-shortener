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
	Slug string `db:"slug" json:"slug"`
	URL  string `db:"url" json:"url"`
}

func (s ShortURL) String() string {
	return fmt.Sprintf("(%s, %s)", s.Slug, s.URL)
}

func (s *ShortURL) pull() {
	store := GetStore()
	store.Pull(s.Slug)
}

func ShortUrlBySlug(slug string) (*ShortURL, error) {
	s, err := db.Dbm.Get(ShortURL{}, slug)
	if err != nil || s == nil {
		return nil, err
	}
	return s.(*ShortURL), nil
}

func CachedShortUrlBySlug(slug string) (*ShortURL, error) {
	urlStore := GetStore()
	url, err := urlStore.Get(slug)
	if err != nil || url == nil {
		return nil, err
	}
	s := &ShortURL{slug, *url}
	return s, nil
}

func ShortURLCreate(url string) (*ShortURL, error) {
	count, err := shortURLCount()
	if err != nil {
		return nil, err
	}
	urlStore := GetStore()
	var slug string
	for {
		slug = genKey(int(count))
		s, err := urlStore.Get(slug)
		if err != nil {
			return nil, err
		}
		if s == nil {
			break
		}
		count++
	}

	s := &ShortURL{slug, url}
	if err = db.Dbm.Insert(s); err != nil {
		return nil, err
	}
	s.pull()
	return s, nil
}

func shortURLCount() (count int64, err error) {
	count, err = db.Dbm.SelectInt("select count(*) from short_urls")
	return
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

var keyChar = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genKey(n int) string {
	if n == 0 {
		return string(keyChar[0])
	}
	l := len(keyChar)
	s := make([]byte, 20) // FIXME: will overflow. eventually.
	i := len(s)
	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = keyChar[j]
	}
	return string(s[i:])
}
