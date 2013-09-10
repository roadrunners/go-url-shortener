package controllers

import (
	m "github.com/roadrunners/go-url-shortener/app/models"
	r "github.com/robfig/revel"
	"net/http"
)

type URL struct {
	Application
}

func (c URL) Create(url string) r.Result {
	s, err := m.ShortUrlCreate(c.XormSession, url)
	if err != nil {
		return c.RenderError(err)
	}
	c.Response.Status = http.StatusCreated
	return c.RenderJson(s)
}

func (c URL) Retrieve(slug string) r.Result {
	s, err := m.CachedShortUrlBySlug(c.XormSession, slug)
	if err != nil {
		return c.RenderError(err)
	}
	if s != nil {
		return c.RenderJson(s)
	}
	return c.NotFound("Short url not found for %s", slug)
}
