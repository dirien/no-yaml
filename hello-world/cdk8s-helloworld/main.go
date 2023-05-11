package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	config := cdk8splus26.NewConfigMap(chart, jsii.String("hello-world-config"), nil)
	config.AddData(jsii.String("nginx.conf"), jsii.String(`events {}
http {
	server {
		listen 80;
		location / {
			return 200 'Hello, World!';
		}
	}
}`))

	volume := cdk8splus26.Volume_FromConfigMap(chart, jsii.String("hello-world-volume"), config, &cdk8splus26.ConfigMapVolumeOptions{
		Items: &map[string]*cdk8splus26.PathMapping{
			"nginx.conf": {
				Path: jsii.String("nginx.conf"),
			},
		},
	})

	deployment := cdk8splus26.NewDeployment(chart, jsii.String("hello-world-deployment"), &cdk8splus26.DeploymentProps{
		Replicas: jsii.Number(1),
	})
	deployment.AddVolume(volume)

	container := deployment.AddContainer(&cdk8splus26.ContainerProps{
		Name:       jsii.String("hello-world"),
		Image:      jsii.String("nginx:latest"),
		PortNumber: jsii.Number(80),
		SecurityContext: &cdk8splus26.ContainerSecurityContextProps{
			EnsureNonRoot:          jsii.Bool(false),
			ReadOnlyRootFilesystem: jsii.Bool(false),
		},
	})
	container.Mount(jsii.String("/etc/nginx/"), volume, &cdk8splus26.MountOptions{
		ReadOnly: jsii.Bool(true),
	})

	_ = deployment.ExposeViaService(&cdk8splus26.DeploymentExposeViaServiceOptions{
		ServiceType: cdk8splus26.ServiceType_LOAD_BALANCER,
	})
	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "cdk8s-helloworld", nil)
	app.Synth()
}
