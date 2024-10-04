package main

import (
	"flag"
	"sync"

	"github.com/lekuruu/hexagon/common"
	"github.com/lekuruu/hexagon/hnet"
	"github.com/lekuruu/hexagon/hscore"
)

type Config struct {
	HNet struct {
		Host string
		Port int
	}
	HScore struct {
		Host string
		Port int
	}
}

func loadConfig() Config {
	var config Config

	flag.StringVar(&config.HNet.Host, "hnet-host", "0.0.0.0", "Host for the hnet server")
	flag.IntVar(&config.HNet.Port, "hnet-port", 21556, "Port for the hnet server")

	flag.StringVar(&config.HScore.Host, "hscore-host", "0.0.0.0", "Host for the hscore server")
	flag.IntVar(&config.HScore.Port, "hscore-port", 80, "Port for the hscore server")

	flag.Parse()

	return config
}

func runService(wg *sync.WaitGroup, worker func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker()
	}()
}

func main() {
	config := loadConfig()

	hnetServer := hnet.NewServer(
		config.HNet.Host,
		config.HNet.Port,
		common.CreateLogger("hnet", common.DEBUG),
	)

	hscoreServer := hscore.NewServer(
		config.HScore.Host,
		config.HScore.Port,
		common.CreateLogger("hscore", common.DEBUG),
	)

	var wg sync.WaitGroup

	runService(&wg, hnetServer.Serve)
	runService(&wg, hscoreServer.Serve)

	wg.Wait()
}
