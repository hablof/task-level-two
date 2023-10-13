package service

import (
	"errors"
	"log"
	"time"

	"github.com/hablof/task-level-two/develop/dev11/internal/models"
)

type Repository interface {
	CreateEvent(event models.Event) (id int, err error)
	UpdateEvent(updateEvent models.UpdateEvent) error
	DeleteEvent(id uint64) error
	EventsInInterval(userID uint64, begin time.Time, end time.Time) ([]models.Event, error)
}

// returns errors
// on known errors
var (
	ErrUnableToSetID   = errors.New("cannot specifiy event id on create")
	ErrMustHaveUserID  = errors.New("new event must belong to some user")
	ErrEmptyTitle      = errors.New("new event must non-empty title")
	ErrNothingToUpdate = errors.New("nothing to update")
)

// returns errors
// on unpredicted repository error
var (
	ErrFailedToCreate = errors.New("failed to create event")
	ErrFailedToUpdate = errors.New("failed to update event")
	ErrFailedToDelete = errors.New("failed to delete event")
	ErrFailedToFetch  = errors.New("failed to fetch events")
)

// recives known repository errors
var (
	ErrNotFound = errors.New("not found")
)

type Service struct {
	r Repository
}

func (s Service) validateEventToCreate(event models.Event) error {
	var err error

	if event.ID != 0 {
		err = errors.Join(err, ErrUnableToSetID)
	}

	if event.UserID == 0 {
		err = errors.Join(err, ErrMustHaveUserID)
	}

	if event.Title == "" {
		err = errors.Join(err, ErrEmptyTitle)
	}

	return err
}

func (s Service) validateEventToUpdate(eventUpdate models.UpdateEvent) error {
	if eventUpdate.Date == nil &&
		eventUpdate.Notes == nil &&
		eventUpdate.Title == nil &&
		eventUpdate.UserID == nil {
		return ErrNothingToUpdate
	}

	return nil
}

func (s Service) CreateEvent(event models.Event) (int, error) {
	if err := s.validateEventToCreate(event); err != nil {
		return 0, err
	}

	id, err := s.r.CreateEvent(event)
	switch {
	case err != nil:
		log.Println(err)
		return 0, ErrFailedToCreate
	}

	return id, nil
}

func (s Service) UpdateEvent(eventUpdate models.UpdateEvent) error {
	if err := s.validateEventToUpdate(eventUpdate); err != nil {
		return err
	}

	switch err := s.r.UpdateEvent(eventUpdate); {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound

	case err != nil:
		log.Println(err)
		return ErrFailedToUpdate
	}

	return nil
}

func (s Service) DeleteEvent(id uint64) error {
	switch err := s.r.DeleteEvent(id); {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound

	case err != nil:
		log.Println(err)
		return ErrFailedToDelete
	}

	return nil
}

func (s Service) EventsForDay(userID uint64, beginDate time.Time) ([]models.Event, error) {
	endDate := beginDate.AddDate(0, 0, 1)
	return s.eventsForPerioD(userID, beginDate, endDate)
}

func (s Service) EventsForWeek(userID uint64, beginDate time.Time) ([]models.Event, error) {
	endDate := beginDate.AddDate(0, 0, 7)
	return s.eventsForPerioD(userID, beginDate, endDate)
}

func (s Service) EventsForMonth(userID uint64, beginDate time.Time) ([]models.Event, error) {
	endDate := beginDate.AddDate(0, 1, 0)
	return s.eventsForPerioD(userID, beginDate, endDate)
}

func (s Service) eventsForPerioD(userID uint64, beginDate time.Time, endDate time.Time) ([]models.Event, error) {
	events, err := s.r.EventsInInterval(userID, beginDate, endDate)
	switch {
	case err != nil:
		log.Println(err)
		return nil, ErrFailedToFetch
	}

	return events, nil
}
