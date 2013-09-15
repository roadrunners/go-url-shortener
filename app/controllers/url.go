package controllers

import (
	"github.com/roadrunners/go-url-shortener/app/models"
	"github.com/robfig/revel"
	"net/http"
)

type URL struct {
	*revel.Controller
}

func (c URL) Create(url string) revel.Result {
	s, err := models.ShortUrlCreate(url)
	if err != nil {
		return c.RenderError(err)
	}
	c.Response.Status = http.StatusCreated
	return c.RenderJson(s)
}

func (c URL) Retrieve(slug string) revel.Result {
	s, err := models.CachedShortUrlBySlug(slug)
	if err != nil {
		return c.RenderError(err)
	}
	if s != nil {
		return c.RenderJson(s)
	}
	return c.NotFound("Short url not found for %s", slug)
}
