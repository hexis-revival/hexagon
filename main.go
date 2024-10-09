package main

import (
	"flag"
	"sync"

	"github.com/hexis-revival/hexagon/common"
	"github.com/hexis-revival/hexagon/hnet"
	"github.com/hexis-revival/hexagon/hscore"
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
	State *common.StateConfiguration
}

func loadConfig() Config {
	var config Config
	config.State = common.NewStateConfiguration()

	flag.StringVar(&config.HNet.Host, "hnet-host", "0.0.0.0", "Host for the hnet server")
	flag.IntVar(&config.HNet.Port, "hnet-port", 21556, "Port for the hnet server")

	flag.StringVar(&config.HScore.Host, "hscore-host", "0.0.0.0", "Host for the hscore server")
	flag.IntVar(&config.HScore.Port, "hscore-port", 80, "Port for the hscore server")

	flag.StringVar(&config.State.Database.Host, "db-host", "localhost", "Database host")
	flag.IntVar(&config.State.Database.Port, "db-port", 5432, "Database port")
	flag.StringVar(&config.State.Database.Username, "db-username", "postgres", "Database username")
	flag.StringVar(&config.State.Database.Password, "db-password", "examplePassword", "Database password")
	flag.StringVar(&config.State.Database.Database, "db-database", "hexagon", "Database name")

	flag.IntVar(&config.State.Database.MaxIdle, "db-max-idle", 10, "Database max idle connections")
	flag.IntVar(&config.State.Database.MaxOpen, "db-max-open", 100, "Database max open connections")
	flag.DurationVar(&config.State.Database.MaxLifetime, "db-max-lifetime", 0, "Database max connection lifetime")

	flag.StringVar(&config.State.Redis.Host, "redis-host", "localhost", "Redis host")
	flag.IntVar(&config.State.Redis.Port, "redis-port", 6379, "Redis port")
	flag.StringVar(&config.State.Redis.Password, "redis-password", "", "Redis password")
	flag.IntVar(&config.State.Redis.Database, "redis-database", 0, "Redis database")

	flag.StringVar(&config.State.DataPath, "data-path", ".data", "Path to store data")
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
	logger := common.CreateLogger("main", common.DEBUG)
	config := loadConfig()

	state, err := common.NewState(config.State)
	if err != nil {
		logger.Error(err)
		return
	}

	hnetServer := hnet.NewServer(
		config.HNet.Host,
		config.HNet.Port,
		common.CreateLogger("hnet", common.DEBUG),
		state,
	)

	hscoreServer := hscore.NewServer(
		config.HScore.Host,
		config.HScore.Port,
		common.CreateLogger("hscore", common.DEBUG),
		state,
	)

	var wg sync.WaitGroup

	runService(&wg, hnetServer.Serve)
	runService(&wg, hscoreServer.Serve)

	wg.Wait()
}
