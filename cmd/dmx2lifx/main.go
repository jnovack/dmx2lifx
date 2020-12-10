package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Hundemeier/go-sacn/sacn"
)

func main() {
	ifi, err := net.InterfaceByName("en0") //this name depends on your machine!
	if err != nil {
		log.Fatal(err)
	}
	recv, err := sacn.NewReceiverSocket("", ifi)
	if err != nil {
		log.Fatal(err)
	}
	recv.SetOnChangeCallback(func(old sacn.DataPacket, newD sacn.DataPacket) {
		fmt.Println("data changed on", newD.Universe())
		fmt.Println(newD.Data())
	})
	recv.SetTimeoutCallback(func(univ uint16) {
		fmt.Println("timeout on", univ)
	})
	recv.Start()
	recv.JoinUniverse(1)
	time.Sleep(10 * time.Second) //only join for 10 seconds, just for testing
	recv.LeaveUniverse(1)
	fmt.Println("Leaved")
	select {} //only that our program does not exit. Exit with Ctrl+C
}
