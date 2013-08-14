package shortener

import (
	r "github.com/robfig/revel"
)

var store *urlStore

func Init() {
	store = newStore()
}

func Put(url string) (string, error) {
	return store.Put(url)
}

func Get(key string) (string, error) {
	return store.Get(key)
}

func init() {
	r.OnAppStart(Init)
}
