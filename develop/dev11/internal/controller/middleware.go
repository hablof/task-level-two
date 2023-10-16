package controller

import (
	"log"
	"net/http"
)

func (c *Controller) logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %s %s %s", r.Host, r.Method, r.URL.String())
		f(w, r)
	}
}
