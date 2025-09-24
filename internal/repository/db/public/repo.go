package public

import (
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	db                *sqlx.DB
	queries           queries
	queryTimeoutLimit time.Duration
	log               logrus.FieldLogger
}
