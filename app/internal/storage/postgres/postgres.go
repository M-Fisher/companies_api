package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	pgxv4 "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/config"
)

type DB struct {
	pool    *pgxpool.Pool
	Queries Queriable
	Log     *zap.Logger
}

type txConn interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Queries struct {
	builder *sq.StatementBuilderType
	tx      txConn
}

type Queriable interface {
	CompaniesQueries
}

func NewPostgres(conf *config.DB, log *zap.Logger) (*DB, error) {
	storage := &DB{
		Log: log,
	}
	dsn := formDbURI(conf)
	log.Debug("creating database connection", zap.String("dsn", conf.Host))
	db, err := checkDB(dsn, conf)
	if err != nil {
		log.Error("database ping error", zap.Error(err))
		return nil, err
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	storage.Queries = &Queries{
		builder: &builder,
		tx:      db,
	}
	storage.pool = db

	return storage, nil
}

func (d *DB) Close() {
	d.Log.Info("Closing DB connection")
	d.pool.Close()
}

func (db *DB) ExecTx(ctx context.Context, f func(q *Queries) error) (err error) {
	tx, err := db.pool.BeginTx(ctx, pgxv4.TxOptions{IsoLevel: pgxv4.ReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = tx.Rollback(ctx)
			}
		} else {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				err = fmt.Errorf("rollback failed %s; %w", txErr.Error(), err)
			}
		}
	}()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := &Queries{tx: tx, builder: &builder}

	err = f(q)

	return err
}

// checkDB Check the DB connection via a given DSN string.
func checkDB(dsn string, conf *config.DB) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(conf.MaxConns)
	dbConfig.MaxConnLifetime = conf.ConnMaxLifetime
	dbConfig.MaxConnIdleTime = conf.ConnMaxLifetime
	return pgxpool.ConnectConfig(context.Background(), dbConfig)
}

func formDbURI(conf *config.DB) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&connect_timeout=10",
		conf.User, conf.Password, conf.Host, conf.Database,
	)
}

func (p *DB) GetStatus() error {
	_, err := p.pool.Exec(context.Background(), `SELECT 1`)
	return err
}
