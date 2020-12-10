package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/jnovack/dmx2lifx/internal/lifx"
)

var waitGroup sync.WaitGroup
var dmxChannel chan []byte

func main() {
	dmxChannel = make(chan []byte)

	waitGroup.Add(1)
	go worker()

	startListener()
	waitGroup.Wait()
}

func startListener() {
	ifi, err := net.InterfaceByName("en0") // TODO Variablize
	if err != nil {
		log.Fatal(err)
	}
	recv, err := sacn.NewReceiverSocket("", ifi)
	if err != nil {
		log.Fatal(err)
	}
	recv.SetOnChangeCallback(func(old sacn.DataPacket, newD sacn.DataPacket) {
		fmt.Println("data changed on", newD.Universe())
		dmxChannel <- newD.Data()
	})
	recv.SetTimeoutCallback(func(univ uint16) {
		fmt.Println("timeout on", univ)
	})
	recv.Start()
	recv.JoinUniverse(1) // TODO Variablize
}

func worker() {
	fmt.Println("worker is now starting...")
	defer func() {
		fmt.Println("destroying the worker...")
		waitGroup.Done()
	}()
	for {
		value, ok := <-dmxChannel
		if !ok {
			fmt.Println("The channel is closed!")
			break
		}
		fmt.Println(value)
		for b := 0; b <= lifx.Count(); b = b + 4 {
			lifx.Set(b, int(value[b+0]), int(value[b+1]), int(value[b+2]), int(value[b+3]))
		}
	}
}
