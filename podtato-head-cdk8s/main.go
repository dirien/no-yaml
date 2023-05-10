package main

import (
	"fmt"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/dirien/podtato-head-cdk8s/internal/podtato"
)

const (
	// Default values
	entryPort = 9000
	hatPort   = 9001
	leftLeg   = 9002
	leftArm   = 9003
	rightLeg  = 9004
	rightArm  = 9005
)

var podtatoHead = &podtato.PodtatoHeadProps{
	NamespaceName: "podtato-head",
	PodtatoParts: []podtato.PodtatoPartArgs{
		{
			PartName:     "entry",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  entryPort,
			ServiceType:  "LoadBalancer",
			IsEntry:      true,
			ServiceDiscoveryData: fmt.Sprintf(`hat:       "http://podtato-head-hat:%d"
left-leg:  "http://podtato-head-left-leg:%d"
left-arm:  "http://podtato-head-left-arm:%d"
right-leg: "http://podtato-head-right-leg:%d"
right-arm: "http://podtato-head-right-arm:%d"
`, hatPort, leftLeg, leftArm, rightLeg, rightArm),
		},
		{
			PartName:     "hat",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  hatPort,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "left-leg",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  leftLeg,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "left-arm",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  leftArm,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "right-leg",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  rightLeg,
			ServiceType:  "ClusterIP",
		},
		{
			PartName:     "right-arm",
			ImageVersion: "0.2.8-chainguard",
			ServicePort:  rightArm,
			ServiceType:  "ClusterIP",
		},
	},
}

func main() {
	app := cdk8s.NewApp(nil)
	podtato.PodtatoHeadChart(app, "podtato-head", podtatoHead)
	app.Synth()
}
