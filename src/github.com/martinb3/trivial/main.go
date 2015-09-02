package main

import (
	"os"
	"fmt"
	"time"
	"github.com/nlopes/slack"
)

func main() {
	SLACK_API_TOKEN := os.Getenv("SLACK_API_TOKEN")
	if SLACK_API_TOKEN == "" {
		fmt.Println("You must set SLACK_API_TOKEN environment variable")
		os.Exit(1)
	}
	api := slack.New(SLACK_API_TOKEN)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
						rtm.SendMessage(rtm.NewOutgoingMessage("Tick", "C02HQ9XT8"))
        }
    }()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				// Replace #general with your Channel ID
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C02HQ9XT8"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
