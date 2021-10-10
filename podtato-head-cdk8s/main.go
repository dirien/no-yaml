package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/dirien/podtato-head-cdk8s/pkg/podtato"
)

var podtatoHead = &podtato.PodtatoHeadProps{
	PodtatoParts: []podtato.PodtatoParts{
		{
			PartName:     "podtato-main",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9000,
			ServiceType:  "LoadBalancer",
		},
		{
			PartName:     "podtato-hats",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9001,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "podtato-left-leg",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9002,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "podtato-left-arm",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9003,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "podtato-right-leg",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9004,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "podtato-right-arm",
			ImageVersion: "v1-latest-dev",
			ServicePort:  9005,
			ServiceType:  "ClusterIP",
		},
	},
}

func main() {
	app := cdk8s.NewApp(nil)
	podtato.PodtatoHeadChart(app, "podtato-head-cdk8s", podtatoHead)
	app.Synth()
}
