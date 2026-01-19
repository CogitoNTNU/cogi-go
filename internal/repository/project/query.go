package projectRepository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type queries struct {
	Read  readQueries
	Write writeQueries
}

type readQueries struct {
	getProject  *sqlx.NamedStmt
	getProjects *sqlx.NamedStmt
}

type writeQueries struct {
	insertProject *sqlx.NamedStmt
	deleteProject *sqlx.NamedStmt
}

func PrepareQueries(db *sqlx.DB) (qs queries, err error) {
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

	if qs.getProject, err = db.PrepareNamed(qGetProject); err != nil {
		return qs, fmt.Errorf("Error preparing GetProject query: %s", err)
	}

	if qs.getProjects, err = db.PrepareNamed(qGetProjects); err != nil {
		return qs, fmt.Errorf("error preparing GetProjects query: %s", err)
	}

	return qs, nil
}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.deleteProject, err = db.PrepareNamed(qDeleteProject); err != nil {
		return qs, fmt.Errorf("Error preparing DeleteProject query: %s", err)
	}

	if qs.insertProject, err = db.PrepareNamed(qInsertProject); err != nil {
		return qs, fmt.Errorf("error preparing InsertProject query: %s", err)
	}

	return qs, nil
}

func (q *queries) Close() (err error) {
	if err = q.Write.insertProject.Close(); err != nil {
		return fmt.Errorf("error closing InsertProject query: %w", err)
	}

	if err = q.Write.deleteProject.Close(); err != nil {
		return fmt.Errorf("error closing DeleteProject query: %w", err)
	}

	if err = q.Read.getProjects.Close(); err != nil {
		return fmt.Errorf("error closing GetProjects query: %w", err)
	}

	if err = q.Read.getProject.Close(); err != nil {
		return fmt.Errorf("error closing GetProject query: %w", err)
	}

	return nil
}
