package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CogitoNTNU/cogi-go/internal/config"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func InitDb(e *env.EnvConfig, logger *logrus.Entry, ctx context.Context) (db *sqlx.DB, err error) {
	var sqlDb *sql.DB

	switch e.Read("ENVRONMENT") {
	case env.DEV:
		cfg := config.NewLocalDbConfig(config.WithDbUser("postgres"), config.WithDbPassword("postgres"), config.WithDbName("master"))
		sqlDb, err = connectDb(ctx, logger, cfg)
	case env.PROD:
	// TODO: implement PROD db
	default:
		return nil, errors.New("unknown environment mode")
	}

	db = sqlx.NewDb(sqlDb, "sqlserver")
	return
}

func connectDb[T config.SQLConfig](ctx context.Context, logger *logrus.Entry, cfg T) (db *sql.DB, err error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", cfg.GetSQLHost(), cfg.GetSQLUser(), cfg.GetSQLPassword(), cfg.GetSQLPort(), cfg.GetSQLDatabase())
	db, err = sql.Open("sqlserver", connString)

	if err != nil {
		return nil, fmt.Errorf("failed to open: %w", err)
	}

	err = db.PingContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to ping local: %w", err)
	}

	logger.Infof("Successfully connected to %s", cfg.GetSQLHost())

	return
}
