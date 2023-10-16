package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"

	"github.com/hablof/task-level-two/develop/dev11/internal/models"
	"github.com/hablof/task-level-two/develop/dev11/internal/service"
)

const (
	defaultRepoTimeout = 15 * time.Second
)

type Repository struct {
	db            *sql.DB
	initStatement sq.StatementBuilderType
}

func NewRepository() (*Repository, error) {

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=rwc", "data.db")) //)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS events(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TIMESTAMP NOT NULL,
		user_id INTEGER NOT NULL,
		title VARCHAR(64) NOT NULL,
		notes VARCHAR(512) NOT NULL);

	CREATE INDEX IF NOT EXISTS idx_user_date ON events(user_id, date);`)
	if err != nil {
		return nil, err
	}

	r := &Repository{
		db:            db,
		initStatement: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}

	return r, nil
}
func (r *Repository) CreateEvent(ctx context.Context, event models.Event) (id int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRepoTimeout)
	defer cancel()

	q, args, err := r.initStatement.Insert("events").
		Columns("date", "user_id", "title", "notes").
		Values(event.Date.TS, event.UserID, event.Title, event.Notes).ToSql()

	if err != nil {
		return 0, err
	}

	result, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	i, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, updateEvent models.UpdateEvent) error {
	ctx, cancel := context.WithTimeout(ctx, defaultRepoTimeout)
	defer cancel()

	ub := r.initStatement.Update("events")

	if updateEvent.Date != nil {
		ub = ub.Set("date", updateEvent.Date.TS)
	}

	if updateEvent.Notes != nil {
		ub = ub.Set("notes", *updateEvent.Notes)
	}

	if updateEvent.Title != nil {
		ub = ub.Set("title", *updateEvent.Title)
	}

	if updateEvent.UserID != nil {
		ub = ub.Set("user_id", *updateEvent.UserID)
	}

	query, args, err := ub.Where(sq.Eq{"id": updateEvent.EventID}).ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, defaultRepoTimeout)
	defer cancel()

	query, args, err := r.initStatement.Delete("events").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (r *Repository) EventsInInterval(ctx context.Context, userID uint64, begin time.Time, end time.Time) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRepoTimeout)
	defer cancel()

	query, args, err := r.initStatement.
		Select("id", "date", "user_id", "title", "notes").
		From("events").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.GtOrEq{"date": begin}).
		Where(sq.LtOrEq{"date": end}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {
		oneEvent := models.Event{}
		if err := rows.Scan(
			&oneEvent.ID,
			&oneEvent.Date.TS,
			&oneEvent.UserID,
			&oneEvent.Title,
			&oneEvent.Notes); err != nil {
			return nil, err
		}

		events = append(events, oneEvent)
	}

	return events, nil
}
