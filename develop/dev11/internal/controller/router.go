package controller

import (
	"net/http"
	"strings"
)

func NewController(s Service) Controller {
	return Controller{
		s: s,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.logging(c.handleRequest)(w, r)
}

func (c *Controller) handleRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		c.handleGet(w, r)
		return

	case http.MethodPost:
		c.handlePost(w, r)
		return
	}

	c.sendErr(w, http.StatusNotFound, "unknown request")
}

func (c *Controller) handleGet(w http.ResponseWriter, r *http.Request) {

	switch strings.Trim(r.URL.Path, "/") {
	case "events_for_day", "events_for_week", "events_for_month":
		c.eventsForInterval(w, r)
		return
	}

	c.sendErr(w, http.StatusNotFound, "unknown request")
}

func (c *Controller) handlePost(w http.ResponseWriter, r *http.Request) {

	switch strings.Trim(r.URL.Path, "/") {
	case "create_event":
		c.createEvent(w, r)
		return

	case "update_event":
		c.updateEvent(w, r)
		return

	case "delete_event":
		c.deleteEvent(w, r)
		return
	}

	c.sendErr(w, http.StatusNotFound, "unknown request")
}
