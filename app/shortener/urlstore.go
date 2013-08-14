package shortener

import (
	"fmt"
	gc "github.com/golang/groupcache"
	r "github.com/robfig/revel"
	"go-url-shortener/app/models"
	"os"
	"os/signal"
)

type urlStore struct {
	cacheGroup *gc.Group
	keys       chan<- string
}

const (
	cacheName       = "URLCache"
	cacheSize       = 64 << 20
	saveQueueLength = 1000
)

func newStore() *urlStore {
	s := new(urlStore)
	s.cacheGroup = gc.NewGroup(cacheName, cacheSize, gc.GetterFunc(s.getter))
	s.keys = s.keysMonitor()
	return s
}

type CannotFindShortUrlError struct {
	key string
}

func (e *CannotFindShortUrlError) Error() string {
	return fmt.Sprintf("Cannot find short url %s", e.key)
}

func (s *urlStore) getter(ctx gc.Context, key string, dest gc.Sink) (err error) {
	r.INFO.Printf("Asked for %s", key)
	shortURL, err := models.ShortUrlBySlug(key)
	if err != nil {
		return
	}

	if shortURL == nil {
		return &CannotFindShortUrlError{key}
	}

	r.INFO.Printf("Found value %s", shortURL)
	dest.SetBytes([]byte(shortURL.URL))
	return nil
}

func (s *urlStore) Get(key string) (string, error) {
	url, err := s.lookupKey(key)
	if err != nil {
		if _, ok := err.(*CannotFindShortUrlError); ok {
			r.ERROR.Print(err)
			return "", err
		}
	}

	r.INFO.Printf("Retrieved short url %s for %s", url, key)

	return url, err
}

func (s *urlStore) Put(url string) (string, error) {
	var key string

	for {
		count, err := models.ShortURLCount()
		if err != nil {
			panic(err)
		}
		key = genKey(int(count))
		err = s.set(key, url)
		if err == nil {
			break
		}
		if _, ok := err.(*CannotFindShortUrlError); ok {
			break
		}
	}

	s.keys <- key

	r.INFO.Printf("Created short url %s for %s", key, url)

	return key, nil
}

type KeyAlreadyPresentError struct {
	key string
}

func (e *KeyAlreadyPresentError) Error() string {
	return fmt.Sprintf("Short url %s already present", e.key)
}

func (s *urlStore) set(key, url string) (err error) {
	taken, err := s.keyAlreadyTaken(key)
	if err != nil {
		return
	}
	if taken {
		return &KeyAlreadyPresentError{key}
	}

	_, err = models.ShortURLCreate(key, url)
	return
}

func (s *urlStore) keysMonitor() chan<- string {
	updates := make(chan string)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		for {
			select {
			case key := <-updates:
				r.INFO.Println("Hottening cache for", key)
				s.lookupKey(key)
			case <-quit:
				r.INFO.Println("Stopping keys monitor")
				return
			}
		}
	}()
	return updates
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
