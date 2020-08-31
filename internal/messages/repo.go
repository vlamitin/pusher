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

func (mr *PgMessageRepo) GetByTimeFrame(from time.Time, to time.Time) ([]Message, error) {
	rows, err := mr.pool.Query(context.Background(), "select * from messages where send_time >= $1 and send_time <= $2 order by send_time desc", from, to)
	if err != nil {
		return []Message{}, fmt.Errorf("error when exec select query: %v\n", err)
	}

	var resultItems []Message

	for rows.Next() {
		var id int
		var sendStatus int
		var sendTime time.Time

		err := rows.Scan(&id, &sendStatus, &sendTime)
		if err != nil {
			return []Message{}, fmt.Errorf("error when read from select query row: %v\n", err)
		}

		resultItems = append(resultItems, *NewMessage(id, sendStatus, sendTime))
	}

	err = rows.Err()
	if err != nil {
		return []Message{}, fmt.Errorf("error after read from select query row: %v\n", err)
	}

	return resultItems, nil
}
