package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

type comms struct {
	command, amount int8
	path            string
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

func ConnectToTwitch(bridge chan comms) {
	client := twitch.NewClient("RedBot", "oauth:muc232tqsq9ethul0qxvqe5ev8b0xr")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)

		func() {
			if message.Message[0] != '>' {
				if message.Message[0] == '#' {
					bridge <- comms{0, 0, message.Message} //Send this data to RoboRio so it can actually be parsed
					//Maybe parse this first and then send, depending on the efficiency and reqs
				}
				return
			}

			inputs := strings.Split(message.Message[1:], " ")
			if len(inputs) < 2 {
				client.Say("redstoneagx", "Invalid Input: "+message.User.DisplayName) // Tell the user they submitted bad input
				return
			}
			amount, err := strconv.Atoi(inputs[1])
			if err != nil {
				client.Say("redstoneagx", "Invalid Input: "+message.User.DisplayName) // Tell the user invalid input
				return
			}

			var dir int8
			switch inputs[0] {
			case "u":
				dir = 0
			case "d":
				dir = 1
			case "ur":
				dir = 2
			case "ul":
				dir = 3
			case "dr":
				dir = 4
			case "dl":
				dir = 5
			}
			bridge <- comms{dir, int8(amount), ""}
		}()

		if message.Message == "!shutdown" && message.User.DisplayName == "RedstoneAGX" {
			client.Say("redstoneagx", "Chat Robot Input Deactivated...")
			bridge <- comms{path: "ex"}
			err := client.Disconnect()
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	client.OnConnect(func() {
		client.Say("redstoneagx", "Chat Robot Input Active...")
	})

	client.Join("redstoneagx")

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
	}
}

func ConnectToRoboRio(message chan comms) {
	// Obtain connection
	conn, err := net.Dial("tcp", "10.20.68.2:1735") // Change the port depending on which port is hosted on the RoboRio
	if err != nil {
		fmt.Println(err)
	}

	for {
		msg := <-message // Check for receiving Data && intepret
		if msg.path == "" {
			conn.Write([]byte{byte(msg.command), byte(msg.amount)}) // Call commands accordingly
		} else if msg.path == "ex" {
			break
		} else {
			conn.Write([]byte("Path:" + msg.path)) //Send to Path Planner
		}
	}
}
