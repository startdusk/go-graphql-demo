package postgres

import (
	"context"
	"embed"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type DB struct {
	Pool *pgxpool.Pool
	URL  string
}

// Example databaseURL
// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
func New(ctx context.Context, databaseURL string, timeout time.Duration) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a pgxpool")
	}
	ready := make(chan struct{})
	go func() {
		for {
			if err := pool.Ping(ctx); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return &DB{
			Pool: pool,
			URL:  databaseURL,
		}, nil
	case <-time.After(timeout * time.Second):
		return nil, errors.New("database not ready")
	}
}

//go:embed migrations/*.sql
var migrationSource embed.FS

func (db *DB) Migrate() error {
	d, err := iofs.New(migrationSource, "migrations") // Get migrations from migrations folder
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithSourceInstance("iofs", d, db.URL)
	if err != nil {
		return errors.Wrap(err, "creating migrator")
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "executing migration")
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		return errors.Wrap(err, "getting migration version")
	}

	log.Printf("postgres migrator version: %d, dirty: %v\n", version, dirty)

	return nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
