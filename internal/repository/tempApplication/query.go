package tempApplicationRepository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type queries struct {
	Read  readQueries
	Write writeQueries
}

type readQueries struct {
	getTempApplication     *sqlx.NamedStmt
	getTempApplicationByEmail *sqlx.NamedStmt
	getAllTempApplications    *sqlx.NamedStmt
}

type writeQueries struct {
	insertTempApplication  *sqlx.NamedStmt
	deleteTempApplication  *sqlx.NamedStmt
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

	if qs.getTempApplication, err = db.PrepareNamed(qGetTempApplication); err != nil {
		return qs, fmt.Errorf("error preparing GetTempApplication query: %s", err)
	}

	if qs.getTempApplicationByEmail, err = db.PrepareNamed(qGetTempApplicationByEmail); err != nil {
		return qs, fmt.Errorf("error preparing GetTempApplicationByEmail query: %s", err)
	}

	if qs.getAllTempApplications, err = db.PrepareNamed(qGetAllTempApplications); err != nil {
		return qs, fmt.Errorf("error preparing GetAllTempApplications query: %s", err)
	}

	return qs, nil
}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.insertTempApplication, err = db.PrepareNamed(qInsertTempApplication); err != nil {
		return qs, fmt.Errorf("error preparing InsertTempApplication query: %s", err)
	}

	if qs.deleteTempApplication, err = db.PrepareNamed(qDeleteTempApplication); err != nil {
		return qs, fmt.Errorf("error preparing DeleteTempApplication query: %s", err)
	}

	return qs, nil
}

func (q *queries) Close() (err error) {
	if q.Read.getTempApplication.Close(); err != nil {
		return fmt.Errorf("error closing GetTempApplication statement: %s", err)
	}
	if q.Read.getTempApplicationByEmail.Close(); err != nil {
		return fmt.Errorf("error closing GetTempApplicationByEmail statement: %s", err)
	}
	if q.Read.getAllTempApplications.Close(); err != nil {
		return fmt.Errorf("error closing GetAllTempApplications statement: %s", err)
	}
	if q.Write.insertTempApplication.Close(); err != nil {
		return fmt.Errorf("error closing InsertTempApplication statement: %s", err)
	}
	if q.Write.deleteTempApplication.Close(); err != nil {
		return fmt.Errorf("error closing DeleteTempApplication statement: %s", err)
	}
	return nil
}