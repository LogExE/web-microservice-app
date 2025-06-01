package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/LogExE/web-microservice-app/config"
	"github.com/LogExE/web-microservice-app/internal"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type app struct {
	csmr *kafka.Consumer
}

func main() {
	cfg := config.New()

	c, err := boxes.MakeConsumer(cfg, "boxesNotifications")
	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}
	defer c.Close()
	err = c.Subscribe("likeBoxes", nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic: ", err)
	}

	a := app{
		csmr: c,
	}

	a.filterLoop()
}

func (a *app) filterLoop() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	run := true
	for run == true {
		select {
		case <-sigchan:
			// SIGINT occured, gracefully stopping
			run = false
		default:
			ev := a.csmr.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				var ev boxes.BoxLikedEvent
				err := json.Unmarshal(e.Value, &ev)
				if err != nil {
					fmt.Printf("%% Kafka value unprocessed: %v\n", e)
				}
				a.noteAct(ev)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				//fmt.Printf("Ignored %v\n", e)
			}
		}
	}

}

func (a *app) noteAct(ev boxes.BoxLikedEvent) {
	fmt.Printf("id %d, likes %d\n", ev.BoxID, ev.Likes)
	if ev.Likes%10 == 0 {
		fmt.Printf("Sending notification to %s\n", ev.Author)
	}
}
