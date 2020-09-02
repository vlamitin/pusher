package messages

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

type PgMessageRepo struct {
	pool *pgxpool.Pool
}

func NewPgMessageRepo(pool *pgxpool.Pool) *PgMessageRepo {
	return &PgMessageRepo{pool: pool}
}

func (mr *PgMessageRepo) Save(sendStatus int, sendTime time.Time) error {
	log.Infof("Saving sendStatus %d and sendTime %v", sendStatus, sendTime)
	_, err := mr.pool.Exec(context.Background(), "insert into messages (status, send_time) values ($1, $2)", sendStatus, sendTime)
	if err != nil {
		return fmt.Errorf("error when exec insert query: %v\n", err)
	}
	log.Infof("Successfully saved sendStatus %d and sendTime %v", sendStatus, sendTime)
	return nil
}

func (mr *PgMessageRepo) GetCountsByTimeFrame(from time.Time, to time.Time) ([]StatusCount, error) {
	rows, err := mr.pool.Query(context.Background(), "select status, count(status) from messages where send_time >= $1 and send_time <= $2 group by status order by status", from, to)
	if err != nil {
		return []StatusCount{}, fmt.Errorf("error when exec select query: %v\n", err)
	}

	var resultItems []StatusCount

	for rows.Next() {
		var sendStatus int
		var count int

		err := rows.Scan(&sendStatus, &count)
		if err != nil {
			return []StatusCount{}, fmt.Errorf("error when read from select query row: %v\n", err)
		}

		resultItems = append(resultItems, *NewStatusCount(sendStatus, count))
	}

	err = rows.Err()
	if err != nil {
		return []StatusCount{}, fmt.Errorf("error after read from select query row: %v\n", err)
	}

	return resultItems, nil
}
