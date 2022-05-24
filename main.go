package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

type comms struct {
	command uint8
	path    string
}

// PATH is if string is != nil
// Up or Right is 1
// Down or Left is 0

// Set up Switch connection on a seperate thread
// Set up RoboRio Connection on a seperate thread
// Have a channel travel through both and be send using both

func main() {
	bridge := make(chan comms)
	go ConnectToTwitch(bridge)
	ConnectToRoboRio(bridge)
}

func ConnectToRoboRio(message chan comms) {
	// Obtain connection
	conn, err := net.Dial("tcp", "10.20.68.2:1735") // Change the port depending on which port is hosted on the RoboRio
	if err != nil {
		fmt.Println(err)
	}

	// Check for receiving Data && intepret
	msg := <-message
	if msg.path != "" {
		//Send to Path Planner
		conn.Write([]byte("Path:" + msg.path))
	}

	// Call commands accordingly
}

func ConnectToTwitch(bridge chan comms) {
	client := twitch.NewClient("RedBot", "oauth:muc232tqsq9ethul0qxvqe5ev8b0xr")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)

		func() {
			if message.Message[0] != '>' {
				return
			}

			if message.Message[0] == '#' {
				bridge <- comms{0, message.Message} //Send this data to RoboRio so it can actually be parsed
				//Maybe parse this first and then send, depending on the efficiency and reqs
				return
			}

			inputs := strings.Split(message.Message[1:], " ")
			switch inputs[0] {
			case "u":
			case "d":
			case "r":
			case "l":
			}
		}()

		if message.Message == "!shutdown" && message.User.DisplayName == "RedstoneAGX" {
			client.Say("redstoneagx", "Chat Robot Input Deactivated...")
			err := client.Disconnect()
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	client.OnConnect(func() {
		client.Say("redstoneagx", "Chat Robot Input Active")
	})

	client.Join("redstoneagx")

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
	}
}
