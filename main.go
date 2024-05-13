package main

import (
	"flag"
	"log"

	"github.com/rajchandvaniya/orangedb/config"
	"github.com/rajchandvaniya/orangedb/server"
)

func main() {
	setupFlags()
	log.Println("peeling the 🍊")
	server.StartAsynchronousTCPServer()
}

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the orange server")
	flag.IntVar(&config.Port, "port", 3690, "port for the orange server")
	flag.IntVar(&config.MaxConnections, "max-connections", 20000, "maximum concurrent connections")
	flag.Parse()
}
