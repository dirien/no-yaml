package main

import (
	"github.com/dirien/podtato-head-naml/internal/podtato"
	"os"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml"
)

type version struct {
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
}

func main() {

	podtatoHead := podtato.NewPodtatoHeadApp()
	naml.Register(podtatoHead)
	err := naml.RunCommandLine()
	if err != nil {
		logger.Critical("%v", err)
		os.Exit(1)
	}
}
