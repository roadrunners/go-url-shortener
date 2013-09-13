package models

import (
	"fmt"
	"github.com/lunny/xorm"
	k "github.com/roadrunners/go-url-shortener/app/models/key"
	"github.com/roadrunners/go-url-shortener/app/redis"
	r "github.com/robfig/revel"
	"os"
	"os/signal"
)

const totalBufferedPulls = 1000

var (
	pulls chan *ShortUrl
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
	r.INFO.Printf("Populating cache %v: %v", key, s.URL)
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
	s.Slug = k.GenKey(s.Id)
	return &s, nil
}

func ShortUrlBySlug(session *xorm.Session, slug string) (*ShortUrl, error) {
	id := k.GenId(slug)
	return ShortUrlById(session, id)
}

func CachedShortUrlBySlug(session *xorm.Session, slug string) (*ShortUrl, error) {
	id := k.GenId(slug)
	key := fmt.Sprintf("shorturl:%d:url", id)
	data, err := redis.Client.Get(key)
	if err == nil {
		s := ShortUrl{Id: id, Slug: slug, URL: string(data)}
		return &s, nil
	}
	r.WARN.Printf("Missed cache for slug %v (id %v, key %v)", slug, id, key)
	s, err := ShortUrlById(session, id)
	if s != nil && err != nil {
		pulls <- s
	}
	return s, err
}

func ShortUrlCreate(session *xorm.Session, url string) (*ShortUrl, error) {
	s := ShortUrl{URL: url}
	_, err := session.Insert(&s)
	if err != nil {
		return nil, err
	}
	s.Slug = k.GenKey(s.Id)
	pulls <- &s
	return &s, nil
}

func pullMonitor() chan *ShortUrl {
	pulls := make(chan *ShortUrl, totalBufferedPulls)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		for {
			select {
			case s := <-pulls:
				s.pull()
			case <-quit:
				r.WARN.Print("Stopping pull monitor")
				close(pulls)
				return
			}
		}
	}()
	return pulls
}

func Init() {
	pulls = pullMonitor()
}

func init() {
	r.OnAppStart(Init)
}
