package main

import (
	"fmt"
	"net"
    "math"
	"strconv"
	"strings"
    "encoding/binary"
    "github.com/gempir/go-twitch-irc/v3"
)

type comms struct {
	angle    uint16
    distance float
	path     string
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

            distance, err := strconv.ParseFloat(inputs[1], 8)
            if err != nil {
				client.Say("redstoneagx", "Invalid Input: "+message.User.DisplayName) // Tell the user invalid input
				return
            }
            angle, err := strconv.ParseUint(inputs[0], 10, 16)
            if err != nil {
				client.Say("redstoneagx", "Invalid Input: "+message.User.DisplayName) // Tell the user invalid input
				return
            }

            bridge <- comms{angle, distance, ""}
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
	// Obtain connection [Maybe use USB later]
	conn, err := net.Dial("tcp", "10.20.68.2:1735") // Change the port depending on which port is hosted on the RoboRio
	if err != nil {
		fmt.Println(err)
	}
    msg_buffer := make([]byte, 6)

	for {
		msg := <-message // Check for receiving Data && intepret
		if msg.path == "" {
            binary.BigEndian.Uint16(msg_buffer[:2], msg.angle)
            binary.BigEndian.Uint32(msg_Buffer[2:], math.Float32bits(msg.distance))
            conn.Write(buffer) // Call commands accordingly
		} else if msg.path == "ex" {
			break
		} else {
			conn.Write([]byte("Path:" + msg.path)) //Send to Path Planner
		}
	}
}
