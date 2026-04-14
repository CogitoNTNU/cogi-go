package eventRepository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type queries struct {
	Read  readQueries
	Write writeQueries
}

type readQueries struct {
	getEvent  *sqlx.NamedStmt
	getEvents *sqlx.NamedStmt
}

type writeQueries struct {
	createEvent *sqlx.NamedStmt
	deleteEvent *sqlx.NamedStmt
}

func PreparedQueries(db *sqlx.DB) (qs queries, err error) {
	readQueries, err := InitRead(db)

	if err != nil {
		return queries{}, err
	}

	writeQueries, err := InitWrite(db)

	if err != nil {
		return queries{}, err
	}

	queries := queries{
		Read:  readQueries,
		Write: writeQueries,
	}

	return queries, nil
}

func InitRead(db *sqlx.DB) (readQueries, error) {
	qs := readQueries{}
	var err error

	if qs.getEvents, err = db.PrepareNamed(qGetEvents); err != nil {
		return qs, fmt.Errorf("Error preparing GetEvents query: %s", err)
	}

	if qs.getEvent, err = db.PrepareNamed(qGetEvent); err != nil {
		return qs, fmt.Errorf("Error preparing GetEvent query: %s", err)
	}

	return qs, nil

}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.createEvent, err = db.PrepareNamed(qCreateEvent); err != nil {
		return qs, fmt.Errorf("Error preparing CreateEvent query: %s", err)
	}

	if qs.deleteEvent, err = db.PrepareNamed(qDeleteEvent); err != nil {
		return qs, fmt.Errorf("Error preparing DeleteEvent query: %s", err)
	}

	return qs, nil

}

func (q *queries) Close() (err error) {
	if err = q.Write.createEvent.Close(); err != nil {
		return fmt.Errorf("Error closing CreateEvent query: %w", err)
	}
	if err = q.Write.deleteEvent.Close(); err != nil {
		return fmt.Errorf("Error closing DeleteEvent query: %w", err)
	}
	if err = q.Read.getEvent.Close(); err != nil {
		return fmt.Errorf("Error closing GetEvent query: %w", err)
	}
	if err = q.Read.getEvents.Close(); err != nil {
		return fmt.Errorf("Error closing GetEvents query: %w", err)
	}
	return nil
}
