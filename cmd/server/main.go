package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer connection.Close()
	fmt.Printf("Connection to RabbidMQ at %s was succesfull\n", connectionString)
	gamelogic.PrintServerHelp()

	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for {
		input := gamelogic.GetInput()
		if len(input) < 1 {
			continue
		} else {
			cmd := input[0]
			switch cmd {
			case "pause":
				if err := pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
					fmt.Printf("PublishJSON fail: %s\n", err.Error())
				}
			case "resume":
				if err := pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false}); err != nil {
					fmt.Printf("PublishJSON fail: %s\n", err.Error())
				}
			case "quit":
				fmt.Println("Gamelogic shutting down...")
				return
			default:
				fmt.Printf("Unknown commands.. please stick to following commands structure")
				gamelogic.PrintServerHelp()
			}
			continue
		}
	}

	//chEnd := make(chan os.Signal, 1)
	//signal.Notify(chEnd, os.Interrupt)
	//<-chEnd
	//fmt.Println("CRTL+C pushed, shutting down Peril server...")
}
