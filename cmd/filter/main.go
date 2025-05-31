package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/LogExE/web-microservice-app/internal"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type app struct {
	csmr *kafka.Consumer
}

func main() {
	c, err := boxes.MakeConsumer("boxesFilter")
	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}
	defer c.Close()
	err = c.Subscribe("newBoxes", nil)
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
				var ev boxes.NewBoxEvent
				err := json.Unmarshal(e.Value, &ev)
				if err != nil {
					fmt.Printf("%% Kafka value unprocessed: %v\n", e)
				}
				a.filterAct(ev)
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

func (a *app) filterAct(ev boxes.NewBoxEvent) {
	fmt.Printf("id %d, by %s\n", ev.BoxID, ev.Author)

	if strings.Contains(ev.Content, "http") {
		fmt.Printf("%d has spam, filtering\n", ev.BoxID)
	}
}
