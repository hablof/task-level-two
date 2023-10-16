package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	TSLayout      = "2006-01-02 15:04"
	TSLayoutShort = "2006-01-02"
)

var (
	ErrNoID = errors.New("update has not event id")
)

type TS struct {
	TS time.Time
}

func ParseTimestamp(str string) (time.Time, error) {
	t, err := time.Parse(TSLayout, str)
	if err != nil {
		tt, err := time.Parse(TSLayoutShort, str)
		if err != nil {
			return time.Time{}, err
		}
		t = tt
	}

	return t, nil
}

func (ts *TS) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `\"`)
	if s == "null" {
		ts.TS = time.Time{}
		return nil
	}

	t, err := ParseTimestamp(s)
	if err != nil {
		return err
	}
	ts.TS = t

	return nil
}

type Event struct {
	ID     uint64 `json:"id,omitempty"`
	Date   TS     `json:"date,omitempty"`
	UserID uint64 `json:"user_id,omitempty"`
	Title  string `json:"title,omitempty"`
	Notes  string `json:"notes,omitempty"`
}

func (e Event) MarshalJSON() ([]byte, error) {

	s := struct {
		ID     uint64 `json:"id,omitempty"`
		Date   string `json:"date,omitempty"`
		UserID uint64 `json:"user_id,omitempty"`
		Title  string `json:"title,omitempty"`
		Notes  string `json:"notes,omitempty"`
	}{
		ID:     e.ID,
		UserID: e.UserID,
		Title:  e.Title,
		Notes:  e.Notes,
	}

	// omit hh:mm if it's 00:00
	if e.Date.TS.Hour() == 0 && e.Date.TS.Minute() == 0 {
		s.Date = e.Date.TS.Format(TSLayoutShort)
	} else {
		s.Date = e.Date.TS.Format(TSLayout)
	}

	return json.Marshal(s)
}

type UpdateEvent struct {
	EventID uint64
	Date    *TS
	UserID  *uint64
	Title   *string
	Notes   *string
}

// patial json handling
func (u *UpdateEvent) UnmarshalJSON(b []byte) error {

	e := Event{}
	tempMap := make(map[string]interface{})

	if err := json.Unmarshal(b, &e); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &tempMap); err != nil {
		return err
	}

	if _, ok := tempMap["id"]; !ok {
		return ErrNoID
	} else {
		u.EventID = e.ID
	}

	u.Date = nil
	u.UserID = nil
	u.Title = nil
	u.Notes = nil

	for k := range tempMap {
		switch k {
		case "date":
			u.Date = &e.Date

		case "user_id":
			u.UserID = &e.UserID

		case "title":
			u.Title = &e.Title

		case "notes":
			u.Notes = &e.Notes
		}
	}

	return nil
}
