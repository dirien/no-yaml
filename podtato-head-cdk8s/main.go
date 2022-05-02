package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/dirien/podtato-head-cdk8s/pkg/podtato"
)

var podtatoHead = &podtato.PodtatoHeadProps{
	PodtatoParts: []podtato.PodtatoParts{
		{
			PartName:     "entry",
			ImageVersion: "0.2.7",
			ServicePort:  9000,
			ServiceType:  "LoadBalancer",
		},
		{
			PartName:     "hat",
			ImageVersion: "0.2.7",
			ServicePort:  9001,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "left-leg",
			ImageVersion: "0.2.7",
			ServicePort:  9002,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "left-arm",
			ImageVersion: "0.2.7",
			ServicePort:  9003,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "right-leg",
			ImageVersion: "0.2.7",
			ServicePort:  9004,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "right-arm",
			ImageVersion: "0.2.7",
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
