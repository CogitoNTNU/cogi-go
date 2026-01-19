package sponsorRepository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type queries struct {
	Read  readQueries
	Write writeQueries
}

type readQueries struct {
	getSponsor  *sqlx.NamedStmt
	getSponsors *sqlx.NamedStmt
}

type writeQueries struct {
	insertSponsor *sqlx.NamedStmt
	deleteSponsor *sqlx.NamedStmt
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

	if qs.getSponsor, err = db.PrepareNamed(qGetSponsor); err != nil {
		return qs, fmt.Errorf("error preparing GetSponsor query: %s", err)
	}

	if qs.getSponsors, err = db.PrepareNamed(qGetSponsors); err != nil {
		return qs, fmt.Errorf("error preparing GetSponsors query: %s", err)
	}

	return qs, nil
}

func InitWrite(db *sqlx.DB) (writeQueries, error) {
	qs := writeQueries{}
	var err error

	if qs.deleteSponsor, err = db.PrepareNamed(qDeleteSponsor); err != nil {
		return qs, fmt.Errorf("error preparing DeleteSponsor query: %s", err)
	}

	if qs.insertSponsor, err = db.PrepareNamed(qInsertSponsor); err != nil {
		return qs, fmt.Errorf("error preparing InsertSponsor query: %s", err)
	}

	return qs, nil
}

func (q *queries) Close() (err error) {
	if err = q.Read.getSponsor.Close(); err != nil {
		return fmt.Errorf("error closing getSponsor statement: %s", err)
	}
	if err = q.Read.getSponsors.Close(); err != nil {
		return fmt.Errorf("error closing getSponsors statement: %s", err)
	}
	if err = q.Write.insertSponsor.Close(); err != nil {
		return fmt.Errorf("error closing insertSponsor statement: %s", err)
	}
	if err = q.Write.deleteSponsor.Close(); err != nil {
		return fmt.Errorf("error closing deleteSponsor statement: %s", err)
	}

	return nil
}
