package main

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/jnovack/dmx2lifx/internal/lifx"
	"github.com/jnovack/flag"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var waitGroup sync.WaitGroup
var dmxChannel chan []byte

var (
	sInterface = flag.String("interface", "eth0", "interface to listen on")
	iUniverse  = flag.Int("universe", 1, "universe to listen on")
)

func main() {
	flag.Parse()

	dmxChannel = make(chan []byte)

	waitGroup.Add(1)
	go worker()

	startListener()
	waitGroup.Wait()
}

func startListener() {
	ifi, err := net.InterfaceByName(*sInterface)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open interface")
	}
	recv, err := sacn.NewReceiverSocket("", ifi)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open socket on interface")
	}
	recv.SetOnChangeCallback(func(old sacn.DataPacket, newD sacn.DataPacket) {
		log.Debug().Msg("received update for universe " + strconv.FormatInt(int64(newD.Universe()), 10))
		dmxChannel <- newD.Data()
	})
	recv.SetTimeoutCallback(func(univ uint16) {
		log.Debug().Msg("timeout receiving data on universe " + strconv.FormatInt(int64(univ), 10))
	})
	recv.Start()
	recv.JoinUniverse(uint16(*iUniverse))
	log.Info().Msg("listening for data on universe " + strconv.FormatInt(int64(*iUniverse), 10))
}

func worker() {
	log.Info().Msg("processing data on universe " + strconv.FormatInt(int64(*iUniverse), 10))
	defer func() {
		log.Warn().Msg("stopped processing data on universe " + strconv.FormatInt(int64(*iUniverse), 10))
		waitGroup.Done()
	}()
	for {
		value, ok := <-dmxChannel
		if !ok {
			log.Error().Msgf("channel for data on universe %d has closed", strconv.FormatInt(int64(*iUniverse), 10))
			break
		}
		log.Trace().Str("universe", strconv.FormatInt(int64(*iUniverse), 10)).Msgf("raw data %v", []byte(value))

		for b := 0; b <= lifx.Count(); b = b + 4 {
			// This always assumes there is at least ONE bulb, otherwise, honestly, why are you running this??
			lifx.Set(b, int(value[b+0]), int(value[b+1]), int(value[b+2]), int(value[b+3]))
		}
	}
}

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		// Format using ConsoleWriter if running straight
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().In(time.Local)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// Format using JSON if running as a service (or container)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}
