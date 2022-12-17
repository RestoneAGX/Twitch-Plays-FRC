package test

import (
	"testing"
)

func TestTwitch(t *testing.T) {

}

func TestUsbTransfer(t *testing.T){
	//Create usb connection
	//Send data to the usb
	//ask for data back (also send null data or a large buffer)
	//Check the actual data
}

func TestQueue(t *testing.T){
	//Send Multiple data
	//Check Reponse time
	//Ask for data back in a queue
	//Ask for the time it took to execute each
}

func TestKill(t *testing.T){
	//Start USB connection
	// Send Commands to be queued
	// Send Kill Command
	// Ask for report on the queue
	// Check if the queued data was executed
}