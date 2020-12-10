package artnet

import (
	"fmt"
	"net"
	"time"

	"github.com/jnovack/dmx2lifx/internal/logging"
)

type pollEvent struct {
	recvTime time.Time
	srcAddr  *net.UDPAddr
	packet   ArtDmx
}

type Config struct {
	Listen    string   `long:"artnet-listen" value-name:"ADDR" default:"0.0.0.0"`
	Discovery []string `long:"artnet-discovery" value-name:"ADDR" default:"255.255.255.255" env:"ARTNET_DISCOVERY" env-delim:","`

	DMXRefresh time.Duration `long:"artnet-dmx-refresh" value-name:"DURATION" default:"1s"`

	Log logging.Option `long:"log.artnet"`
}

func (config Config) Controller() (*Controller, error) {
	config.Log.Package = "artnet"

	var controller = Controller{
		config: config,

		universes: make(map[Address]*Universe),
	}

	controller.log = config.Log.Logger("controller", &controller)

	listenAddr := net.JoinHostPort(config.Listen, fmt.Sprintf("%d", Port))

	if udpAddr, err := net.ResolveUDPAddr("udp", listenAddr); err != nil {
		return nil, err
	} else if udpConn, err := net.ListenUDP("udp", udpAddr); err != nil {
		return nil, err
	} else {
		controller.transport = &Transport{
			udpConn: udpConn,
		}
	}

	return &controller, nil
}

type Controller struct {
	log logging.Logger

	config Config

	transport *Transport
	pollChan  chan pollEvent

	// state
	universes map[Address]*Universe
}

func (controller *Controller) String() string {
	return fmt.Sprintf("%v", controller.transport)
}

func (controller *Controller) Start() {
	controller.pollChan = make(chan pollEvent)

	go controller.recv()
}

func (controller *Controller) recv() {
	for {
		if packet, srcAddr, err := controller.transport.recv(); err != nil {
			// XXX: fatal if socket is dead?
			controller.log.Errorf("recv %v: %v", srcAddr, err)
		} else if err := controller.recvPacket(packet, srcAddr); err != nil {
			controller.log.Warnf("recv %v: %v", srcAddr, err)
		}
	}
}

func (controller *Controller) recvPacket(packet ArtPacket, srcAddr *net.UDPAddr) error {
	switch packetType := packet.(type) {
	case *ArtDmx:
		if !ProtVer14.IsCompatible(packetType.ProtVer) {
			return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer14)
		}

		// ignore
		return nil
	}

	return nil
}
