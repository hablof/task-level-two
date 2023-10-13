package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateEvent_UnmarshalJSON(t *testing.T) {
	var (
		testString string = "nazvanie"
	)
	testTS, err := time.Parse(TSLayout, "2021-07-19 15:20")
	if err != nil {
		t.FailNow()
		return
	}

	tests := []struct {
		name    string
		b       []byte
		wantUE  UpdateEvent
		wantErr error
	}{
		{
			name:    "error no id",
			b:       []byte(`{"title":"nazvanie"}`),
			wantUE:  UpdateEvent{},
			wantErr: ErrNoID,
		},
		{
			name: "default ok",
			b:    []byte(`{"id":15,"title":"nazvanie"}`),
			wantUE: UpdateEvent{
				EventID: 15,
				Title:   &testString,
			},
			wantErr: nil,
		},
		{
			name: "test parse datetime",
			b:    []byte(`{"id":15,"title":"nazvanie","date":"2021-07-19 15:20"}`),
			wantUE: UpdateEvent{
				EventID: 15,
				Date:    &TS{testTS},
				Title:   &testString,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UpdateEvent{}
			err := json.Unmarshal(tt.b, &u)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUE, u)
		})
	}
}
