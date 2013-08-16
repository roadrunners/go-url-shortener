package store

import (
	"fmt"
	gc "github.com/golang/groupcache"
	r "github.com/robfig/revel"
	"os"
	"os/signal"
	"sync"
)

const (
	cacheSize       = 64 << 20
	saveQueueLength = 1000
)

type Getter interface {
	Get(key string) (*string, error)
}

type GetterFunc func(key string) (*string, error)

func (f GetterFunc) Get(key string) (*string, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	stores = make(map[string]*Store)
)

func GetStore(name string) *Store {
	mu.RLock()
	defer mu.RUnlock()
	return stores[name]
}

func NewStore(name string, getter Getter) *Store {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	s := new(Store)
	s.name = name
	s.cacheGroup = gc.NewGroup(s.name, cacheSize, gc.GetterFunc(s.populate))
	s.keys = s.keysMonitor()
	s.getter = getter
	stores[name] = s
	return s
}

type Store struct {
	name       string
	cacheGroup *gc.Group
	keys       chan<- string
	getter     Getter
}

func (s *Store) Get(key string) (*string, error) {
	value, err := s.lookupKey(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	} else {
		return value, err
	}
}

func (s *Store) populate(ctx gc.Context, key string, dest gc.Sink) error {
	r.INFO.Printf("Asked for %s", key)
	item, err := s.getter.Get(key)
	if err != nil {
		return err
	}
	if item == nil {
		r.ERROR.Printf("Could not find key %s", key)
		return &CannotFindKeyError{key}
	}
	r.INFO.Printf("Found value %s for %s", *item, key)
	dest.SetBytes([]byte(*item))
	return nil
}

type CannotFindKeyError struct {
	key string
}

func (e *CannotFindKeyError) Error() string {
	return fmt.Sprintf("Cannot find key %s", e.key)
}

func (s *Store) Pull(key string) {
	go func() {
		s.keys <- key
	}()
}

func (s *Store) keysMonitor() chan<- string {
	updates := make(chan string, saveQueueLength)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit)
	go func() {
		for {
			select {
			case key := <-updates:
				r.INFO.Println("Hottening cache for", key)
				s.lookupKey(key)
			case <-quit:
				r.INFO.Printf("Stopping keys monitor for %s", s.name)
				close(updates)
				return
			}
		}
	}()
	return updates
}

func (s *Store) keyAlreadyTaken(key string) (bool, error) {
	v, err := s.lookupKey(key)
	if err != nil {
		return false, err
	}
	if v == nil {
		return false, nil
	}
	return true, nil
}

func (s *Store) lookupKey(key string) (*string, error) {
	var data []byte
	err := s.cacheGroup.Get(nil, key, gc.AllocatingByteSliceSink(&data))
	if err != nil {
		if _, ok := err.(*CannotFindKeyError); ok {
			return nil, nil
		}
		return nil, err
	}
	value := string(data)
	return &value, nil
}
