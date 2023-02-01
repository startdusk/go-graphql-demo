package postgres

import (
	"context"
	"embed"
	"log"
	"os"

	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	zerologadapter "github.com/jackc/pgx-zerolog"
)

type DB struct {
	Pool *pgxpool.Pool
	URL  string
}

// Example databaseURL
// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
func New(ctx context.Context, databaseURL string, timeout time.Duration) (*DB, error) {
	conf, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a pgxpool")
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// output.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }
	// output.FormatMessage = func(i interface{}) string {
	// 	return fmt.Sprintf("***%s****", i)
	// }
	// output.FormatFieldName = func(i interface{}) string {
	// 	return fmt.Sprintf("%s:", i)
	// }
	// output.FormatFieldValue = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("%s", i))
	// }

	zlogger := zerolog.New(output).With().Timestamp().Logger()
	logger := zerologadapter.NewLogger(zlogger)
	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: tracelog.LogLevelDebug,
	}
	pool, err := pgxpool.NewWithConfig(ctx, conf)
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
