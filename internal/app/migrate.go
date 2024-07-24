package app

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"time"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func init() {
	dbUrl, ok := os.LookupEnv("PG_URL")
	if !ok {
		log.Fatal("migration error: PG url is not declared")
	}
	dbUrl += "?sslmode=disable"
	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)
	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbUrl)
		if err == nil {
			break
		}
		log.Printf("migration trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}
	if err != nil {
		log.Fatalf("migration db connect error: %s", err)
	}
	err = m.Up()

	defer func() { _, _ = m.Close() }()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration up error: %s", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrations without change")
		return
	}
	log.Printf("migration up success")

}
