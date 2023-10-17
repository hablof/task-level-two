package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hablof/task-level-two/develop/dev11/internal/models"
	"github.com/hablof/task-level-two/develop/dev11/internal/service"
)

type result struct {
	Res    []models.Event `json:"result,omitempty"`
	ErrMsg string         `json:"error,omitempty"`
}

type Service interface {
	CreateEvent(ctx context.Context, event models.Event) (int64, error)
	UpdateEvent(ctx context.Context, eventUpdate models.UpdateEvent) error
	DeleteEvent(ctx context.Context, id uint64) error
	EventsForInterval(ctx context.Context, userID uint64, beginDate time.Time, interval service.Interval) ([]models.Event, error)
}

type Controller struct {
	s Service
}

func (c *Controller) sendErr(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": errMsg}); err != nil {
		log.Printf("failed to encode err: %v", err)
	}
}

func (c *Controller) createEvent(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := strconv.ParseUint(r.PostForm.Get("user_id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse user id")

		return
	}

	timeStamp, err := models.ParseTimestamp(r.PostForm.Get("date"))
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse date")

		return
	}

	newEvent := models.Event{
		Date: models.TS{
			TS: timeStamp,
		},
		UserID: userID,
		Title:  r.PostForm.Get("title"),
		Notes:  r.PostForm.Get("notes"),
	}

	newEventID, err := c.s.CreateEvent(r.Context(), newEvent)
	if err != nil {
		log.Println(err)

		switch {
		case errors.Is(err, service.ErrUnableToSetID) ||
			errors.Is(err, service.ErrMustHaveUserID) ||
			errors.Is(err, service.ErrEmptyTitle):

			c.sendErr(w, http.StatusServiceUnavailable, err.Error())

		default:
			c.sendErr(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	result := struct {
		NewEventId int64 `json:"new_event_id"`
	}{
		NewEventId: newEventID,
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"result": result}); err != nil {
		log.Printf("failed to encode result: %v", err)
	}
	log.Println("success")
}

func (c *Controller) updateEvent(w http.ResponseWriter, r *http.Request) {

	ue := models.UpdateEvent{}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, err.Error())

		return
	}

	updatingEventID, err := strconv.ParseUint(r.PostForm.Get("id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse event id")

		return
	}

	ue.EventID = updatingEventID

	if r.PostForm.Has("date") {
		timeStamp, err := models.ParseTimestamp(r.PostForm.Get("date"))
		if err != nil {
			log.Println(err)
			c.sendErr(w, http.StatusBadRequest, "cannot parse date")

			return
		}

		ue.Date = &models.TS{
			TS: timeStamp,
		}
	}

	if r.PostForm.Has("user_id") {
		userId, err := strconv.ParseUint(r.PostForm.Get("user_id"), 10, 64)
		if err != nil {
			log.Println(err)
			c.sendErr(w, http.StatusBadRequest, "cannot parse user id")

			return
		}

		ue.UserID = &userId
	}

	if r.PostForm.Has("title") {
		newTitle := r.PostForm.Get("title")
		ue.Title = &newTitle
	}

	if r.PostForm.Has("notes") {
		newNotes := r.PostForm.Get("notes")
		ue.Notes = &newNotes
	}

	if err := c.s.UpdateEvent(r.Context(), ue); err != nil {
		log.Println(err)

		switch err {
		case service.ErrNothingToUpdate:
			c.sendErr(w, http.StatusServiceUnavailable, err.Error())

		case service.ErrNotFound:
			c.sendErr(w, http.StatusNotFound, err.Error())

		default:
			c.sendErr(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string]string{"result": "updated"}); err != nil {
		log.Printf("failed to encode result: %v", err)
	}
	log.Println("success")
}

func (c *Controller) deleteEvent(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, err.Error())

		return
	}

	deletingEventID, err := strconv.ParseUint(r.PostForm.Get("id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse event id")

		return
	}

	if err := c.s.DeleteEvent(r.Context(), deletingEventID); err != nil {
		log.Println(err)

		switch err {
		case service.ErrNotFound:
			c.sendErr(w, http.StatusNotFound, err.Error())

		default:
			c.sendErr(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string]string{"result": "deleted"}); err != nil {
		log.Printf("failed to encode result: %v", err)
	}
	log.Println("success")
}

func (c *Controller) eventsForInterval(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := strconv.ParseUint(r.Form.Get("user_id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse user id")

		return
	}

	timeStamp, err := models.ParseTimestamp(r.Form.Get("date"))
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusBadRequest, "cannot parse date")

		return
	}

	var interval service.Interval
	switch strings.Trim(r.URL.Path, "/") {
	case "events_for_day":
		interval = service.Day

	case "events_for_week":
		interval = service.Week

	case "events_for_month":
		interval = service.Month
	}

	events, err := c.s.EventsForInterval(r.Context(), userID, timeStamp, interval)
	if err != nil {
		log.Println(err)
		c.sendErr(w, http.StatusInternalServerError, err.Error())

		return
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string][]models.Event{"result": events}); err != nil {
		log.Printf("failed to encode result: %v", err)
		return
	}
	log.Println("success")
}
