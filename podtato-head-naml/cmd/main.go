package main

import (
	"os"

	"github.com/dirien/podtato-head-naml"
	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml"
)

func main() {

	podtatoHead := podtato.NewPodtatoHeadApp()
	naml.Register(podtatoHead)
	err := naml.RunCommandLine()
	if err != nil {
		logger.Critical("%v", err)
		os.Exit(1)
	}
}
