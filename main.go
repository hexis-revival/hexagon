package main

import (
	"github.com/lekuruu/hexagon/common"
	"github.com/lekuruu/hexagon/hnet"
)

func main() {
	// TODO: Add command line arguments for host and port
	server := hnet.NewServer(
		"0.0.0.0",
		21556,
		common.CreateLogger("hnet", common.DEBUG),
	)

	server.Serve()
}
