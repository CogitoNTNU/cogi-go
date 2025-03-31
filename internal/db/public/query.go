package public

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type queries struct {
  Read readQueries
  Write writeQueries
}

type readQueries struct {
	getMembers          *sqlx.NamedStmt
  getMember           *sqlx.NamedStmt
}

type writeQueries struct {
  insertMember        *sqlx.NamedStmt
  updateMember        *sqlx.NamedStmt
  deleteMember        *sqlx.NamedStmt
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

	if qs.getMember, err = db.PrepareNamed(qGetMember); err != nil {
    return qs, fmt.Errorf("Error preparing GetMember query: %s", err)
	}

  if qs.getMembers, err = db.PrepareNamed(qGetMembers); err != nil {
    return qs, fmt.Errorf("error preparing getmembers query: %s", err)
  }

	return qs, nil
}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.deleteMember, err = db.PrepareNamed(qDeleteMember); err != nil {
    return qs, fmt.Errorf("Error preparing DeleteMember query: %s", err)
	}

  if qs.insertMember, err = db.PrepareNamed(qInsertMember); err != nil {
    return qs, fmt.Errorf("error preparing InsertMember query: %s", err)
  }

  if qs.updateMember, err = db.PrepareNamed(qUpdateMember); err != nil {
    return qs, fmt.Errorf("error preparing InsertMember query: %s", err)
  }

	return qs, nil
}

func (q *queries) Close() (err error) {
	if err = q.Write.updateMember.Close(); err != nil {
		return fmt.Errorf("error closing UpdateMember query: %w", err)
	}

  if err = q.Write.insertMember.Close(); err != nil {
    return fmt.Errorf("error closing InsertMember query: %w", err)
  }

  if err = q.Write.deleteMember.Close(); err != nil {
    return fmt.Errorf("error closing DeleteMember query: %w", err)
  }

  if err = q.Read.getMembers.Close(); err != nil {
    return fmt.Errorf("error closing GetMembers query: %w", err)
  }

  if err = q.Read.getMember.Close(); err != nil {
    return fmt.Errorf("error closing GetMember query: %w", err)
  }

	return nil
}


const (
  qGetMember = ``
  qGetMembers = ``
  qInsertMember = ``
  qUpdateMember = ``
  qDeleteMember = ``
)
