package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePool(postgresHost string,
	postgresPort string,
	postgresUser string,
	postgresPassword string,
	postgresDb string) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDb)
	return pgxpool.Connect(context.Background(), connString)
}
