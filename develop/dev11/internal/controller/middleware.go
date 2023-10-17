package controller

import (
	"log"
	"net/http"
)

func (c *Controller) logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Printf("request: %s %-4s %-16s %s", r.Host, r.Method, r.URL.Path, r.Form.Encode())
		f(w, r)
	}
}
