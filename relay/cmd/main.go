package main

import "github.com/notedit/relay-server/relay"

var configFile = "../../config.yaml"

func main() {
	server, err := relay.NewRelayServer(configFile)
	if err != nil {
		panic(err)
	}

	server.Run()
}
