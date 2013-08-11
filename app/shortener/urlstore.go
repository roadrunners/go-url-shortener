package shortener

import (
	"fmt"
	gc "github.com/golang/groupcache"
	"github.com/robfig/revel"
)

var (
	db = make(map[string]string)
)

type urlStore struct {
	cacheGroup *gc.Group
	keys       chan string
}

const (
	cacheName       = "URLCache"
	cacheSize       = 64 << 20
	saveQueueLength = 1000
)

func newStore() *urlStore {
	s := new(urlStore)
	s.cacheGroup = gc.NewGroup(cacheName, cacheSize, gc.GetterFunc(s.getter))
	s.keys = make(chan string, saveQueueLength)
	return s
}

type CannotFindShortUrlError struct {
	key string
}

func (e *CannotFindShortUrlError) Error() string {
	return fmt.Sprintf("Cannot find short url %s", e.key)
}

func (s *urlStore) getter(ctx gc.Context, key string, dest gc.Sink) error {
	revel.INFO.Printf("Asked for %s", key)
	result, present := db[key]
	if present {
		revel.INFO.Printf("Found value %s for %s", result, key)
		dest.SetBytes([]byte(result))
		return nil
	}

	return &CannotFindShortUrlError{key}
}

func (s *urlStore) Get(key string) (string, error) {
	url, err := s.lookupKey(key)
	if err != nil {
		if _, ok := err.(*CannotFindShortUrlError); ok {
			revel.ERROR.Print(err)
			return "", err
		}
	}

	revel.INFO.Printf("Retrieved short url %s for %s", url, key)

	return url, err
}

func (s *urlStore) Put(url string) (string, error) {
	var key string

	for {
		key = genKey(len(db))
		err := s.set(key, url)
		if err == nil {
			break
		}
		if _, ok := err.(*CannotFindShortUrlError); ok {
			break
		}
	}

	s.keys <- key

	revel.INFO.Printf("Created short url %s for %s", key, url)

	return key, nil
}

type KeyAlreadyPresentError struct {
	key string
}

func (e *KeyAlreadyPresentError) Error() string {
	return fmt.Sprintf("Short url %s already present", e.key)
}

func (s *urlStore) set(key, url string) error {
	taken, err := s.keyAlreadyTaken(key)
	if err != nil {
		return err
	}
	if taken {
		return &KeyAlreadyPresentError{key}
	}

	db[key] = url
	return nil
}

func (s *urlStore) cacheKeys() {
	for {
		key := <-s.keys
		revel.INFO.Println("Hottening cache for", key)
		s.lookupKey(key)
	}
}

func (s *urlStore) keyAlreadyTaken(key string) (taken bool, err error) {
	_, err = s.lookupKey(key)
	if err != nil {
		if _, ok := err.(*CannotFindShortUrlError); ok {
			err = nil
		}

		return
	}

	taken = true
	return
}

func (s *urlStore) lookupKey(key string) (url string, err error) {
	var data []byte
	err = s.cacheGroup.Get(nil, key, gc.AllocatingByteSliceSink(&data))
	if err == nil {
		url = string(data)
	}

	return
}
