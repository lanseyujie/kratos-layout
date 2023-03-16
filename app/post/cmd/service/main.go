package main

import (
	"flag"

	_ "go.uber.org/automaxprocs"
)

func main() {
	configDir := flag.String("config", "./configs/", "config path, eg: -config ./configs/")
	flag.Parse()

	app, cleanup, err := wireApp(*configDir)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		panic(err)
	}
}
