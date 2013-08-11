package controllers

import (
	"github.com/robfig/revel"
	"go-url-shortener/app/shortener"
	"net/http"
)

type URL struct {
	*revel.Controller
}

type CreateResponse struct {
	Slug string
}

func (c URL) Create(url string) revel.Result {
	slug, err := shortener.Put(url)
	if err != nil {
		return c.RenderError(err)
	}

	c.Response.Status = http.StatusCreated
	return c.RenderJson(CreateResponse{slug})
}

type RetrieveResponse struct {
	URL string
}

func (c URL) Retrieve(slug string) revel.Result {
	url, err := shortener.Get(slug)
	if err != nil {
		if _, ok := err.(*shortener.CannotFindShortUrlError); ok {
			return c.NotFound("Short url not found for %s", slug)
		}

		return c.RenderError(err)
	}

	return c.RenderJson(RetrieveResponse{url})
}
