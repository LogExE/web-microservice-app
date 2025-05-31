package boxes

import (
	"database/sql"
	_ "modernc.org/sqlite"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "boxes.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS box (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			content TEXT NOT NULL, 
			author TEXT NOT NULL,
                        likes INTEGER NOT NULL
		);
                CREATE TABLE IF NOT EXISTS outbox (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
                        processed INTEGER NOT NULL,
                        type TEXT NOT NULL,
                        payload BLOB NOT NULL
		)`,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MakeConsumer(group string) (*kafka.Consumer, error) {
	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          group,
		"auto.offset.reset": "smallest",
	})

	if err != nil {
		return nil, err
	}
	return p, nil
}

func MakeProducer(client string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         client,
		"acks":              "all"})

	if err != nil {
		return nil, err
	}
	return p, nil
}
