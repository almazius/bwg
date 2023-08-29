package utils

import (
	"bwg2/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func connectWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for i := 0; i < attempts; i++ {
		err = fn()
		if err != nil {
			log.Print(fmt.Sprintf("Try num: %d, connection failed", i))
			time.Sleep(delay)
			continue
		}
	}
	return err
}

// GetPool Creates a connection pool
func GetPool(ctx context.Context, conf *config.Config) (connect *pgxpool.Pool, err error) {
	err = connectWithTries(func() error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		connect, err = pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			conf.Postgres.User,
			conf.Postgres.Password,
			conf.Postgres.Host,
			conf.Postgres.Port,
			conf.Postgres.DbName,
			"disable"))
		return err
	}, 3, time.Second*3)
	if err != nil || connect.Ping(ctx) != nil {
		log.Fatal(err, " unable to connect to database.")
		return nil, err
	}
	return connect, nil
}
