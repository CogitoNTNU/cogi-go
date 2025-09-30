package userRepository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type queries struct {
	Read  readQueries
	Write writeQueries
}

type readQueries struct {
	getUser  *sqlx.NamedStmt
	getUsers *sqlx.NamedStmt
}

type writeQueries struct {
	insertUser *sqlx.NamedStmt
	deleteUser *sqlx.NamedStmt
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

	if qs.getUser, err = db.PrepareNamed(qGetUser); err != nil {
		return qs, fmt.Errorf("Error preparing GetMember query: %s", err)
	}

	if qs.getUsers, err = db.PrepareNamed(qGetUsers); err != nil {
		return qs, fmt.Errorf("error preparing getmembers query: %s", err)
	}

	return qs, nil
}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.deleteUser, err = db.PrepareNamed(qDeleteUser); err != nil {
		return qs, fmt.Errorf("Error preparing DeleteMember query: %s", err)
	}

	if qs.insertUser, err = db.PrepareNamed(qInsertUser); err != nil {
		return qs, fmt.Errorf("error preparing InsertMember query: %s", err)
	}

	return qs, nil
}

func (q *queries) Close() (err error) {
	if err = q.Write.insertUser.Close(); err != nil {
		return fmt.Errorf("error closing InsertMember query: %w", err)
	}

	if err = q.Write.deleteUser.Close(); err != nil {
		return fmt.Errorf("error closing DeleteMember query: %w", err)
	}

	if err = q.Read.getUsers.Close(); err != nil {
		return fmt.Errorf("error closing GetMembers query: %w", err)
	}

	if err = q.Read.getUser.Close(); err != nil {
		return fmt.Errorf("error closing GetMember query: %w", err)
	}

	return nil
}
