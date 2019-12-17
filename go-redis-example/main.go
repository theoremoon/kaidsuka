package main

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/go-redis/redis/v7"
)

func newClient() {
	c := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	name := randomdata.SillyName()

	pubsub := c.Subscribe("greeting")
	c.Publish("greeting", "ドーモ、皆サン。"+name+"デス")

	sub := pubsub.Channel()
	for {
		select {
		case m := <-sub:
			fmt.Println(m.Payload)
		}
	}
}

func main() {
	newClient()
}
