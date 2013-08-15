package controllers

import (
	r "github.com/robfig/revel"
	m "go-url-shortener/app/models"
	"net/http"
)

type URL struct {
	*r.Controller
}

func (c URL) Create(url string) r.Result {
	s, err := m.ShortURLCreate(url)
	if err != nil {
		return c.RenderError(err)
	}
	c.Response.Status = http.StatusCreated
	return c.RenderJson(s)
}

func (c URL) Retrieve(slug string) r.Result {
	s, err := m.CachedShortUrlBySlug(slug)
	if err != nil {
		return c.RenderError(err)
	}
	if s != nil {
		return c.RenderJson(s)
	}
	return c.NotFound("Short url not found for %s", slug)
}
