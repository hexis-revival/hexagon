package main

import (
	"sync"

	"github.com/lekuruu/hexagon/common"
	"github.com/lekuruu/hexagon/hnet"
	"github.com/lekuruu/hexagon/hscore"
)

func runService(wg *sync.WaitGroup, worker func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker()
	}()
}

func main() {
	// TODO: Add command line arguments for host and port
	hnetServer := hnet.NewServer(
		"0.0.0.0",
		21556,
		common.CreateLogger("hnet", common.DEBUG),
	)

	hscoreServer := hscore.NewServer(
		"0.0.0.0",
		80,
		common.CreateLogger("hscore", common.DEBUG),
	)

	var wg sync.WaitGroup

	runService(&wg, hnetServer.Serve)
	runService(&wg, hscoreServer.Serve)

	wg.Wait()
}
