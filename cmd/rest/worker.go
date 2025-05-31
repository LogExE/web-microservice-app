package main

import (
	"log"

	boxes "github.com/LogExE/web-microservice-app/internal"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
)

var newBoxesTopic = "newBoxes"
var likesTopic = "likeBoxes"

func outboxWorker(p *kafka.Producer, b *boxes.BoxRepo) {
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		outboxList, err := b.OutboxList()
		if err != nil {
			log.Println("ERROR: Error occured when getting outbox list from database")
			return
		}

		if len(outboxList) != 0 {
			for _, e := range outboxList {
				delivery_chan := make(chan kafka.Event, 100)

				var topic *string
				if e.Type == "boxNew" {
					topic = &newBoxesTopic
				} else if e.Type == "boxLike" {
					topic = &likesTopic
				}
				p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
					Value:          e.Payload,
				}, delivery_chan)

				ev := <-delivery_chan
				m := ev.(*kafka.Message)

				if m.TopicPartition.Error != nil {
					log.Printf("ERROR: Delivery failed: %v\n", m.TopicPartition.Error)
					return
				} else {
					log.Printf("INFO: Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				close(delivery_chan)

				err = b.OutboxClear(e.ID)
				if err != nil {
					log.Printf("ERROR: Error occured when event status updated with boxID -> %d\n", e.ID)
					return
				}
			}

		}

	})
	c.Start()
}
