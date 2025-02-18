package intiator

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func InitDB(url string, log *zap.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	// if idleConnTimeout == 0 {
	// 	idleConnTimeout = 4 * time.Minute
	// }

	config.MaxConnIdleTime = 4 * time.Minute
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if _, err := conn.Exec(context.Background(), "show tables"); err != nil {
		log.Fatal(fmt.Sprintf("Failed to ping database: %v", err))
	}

	return conn
}
