package podtato

import (
	"fmt"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/dirien/podtato-head-cdk8s/imports/k8s"
)

type PodtatoHeadProps struct {
	cdk8s.ChartProps
	PodtatoParts []PodtatoParts
}

type PodtatoParts struct {
	PartName     string
	ImageVersion string
	ServicePort  int
	ServiceType  string
}

func buildPodtatoHeadComponent(scope constructs.Construct, ns k8s.KubeNamespace, appName map[string]*string, componentName,
	imageVersion string, servicePort int, serviceType string) {
	componentLabel := map[string]*string{"component": jsii.String(componentName)}

	name := fmt.Sprintf("podtato-head-%s", componentName)

	k8s.NewKubeDeployment(scope, jsii.String(fmt.Sprintf("%s-depl", componentName)), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String(name),
			Namespace: ns.Metadata().Name(),
			Labels:    &appName,
		},
		Spec: &k8s.DeploymentSpec{
			Selector: &k8s.LabelSelector{
				MatchLabels: &componentLabel,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &componentLabel,
				},
				Spec: &k8s.PodSpec{
					TerminationGracePeriodSeconds: jsii.Number(5),
					Containers: &[]*k8s.Container{
						{
							Name:            jsii.String("server"),
							Image:           jsii.String(fmt.Sprintf("ghcr.io/podtato-head/%s:%s", componentName, imageVersion)),
							ImagePullPolicy: jsii.String("Always"),
							Ports: &[]*k8s.ContainerPort{
								{
									ContainerPort: jsii.Number(9000),
								},
							},
							Env: &[]*k8s.EnvVar{
								{
									Name:  jsii.String("PORT"),
									Value: jsii.String("9000"),
								},
							},
						},
					},
				},
			},
		},
	})

	k8s.NewKubeService(scope, jsii.String(fmt.Sprintf("%s-svc", componentName)), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String(name),
			Namespace: ns.Metadata().Name(),
			Labels:    &appName,
		},
		Spec: &k8s.ServiceSpec{
			Selector: &componentLabel,
			Ports: &[]*k8s.ServicePort{
				{
					Name:       jsii.String("http"),
					Port:       jsii.Number(float64(servicePort)),
					Protocol:   jsii.String("TCP"),
					TargetPort: k8s.IntOrString_FromNumber(jsii.Number(float64(9000))),
				},
			},
			Type: jsii.String(serviceType),
		},
	})
}

func PodtatoHeadChart(scope constructs.Construct, id string, props *PodtatoHeadProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	appLabels := map[string]*string{"app": jsii.String("podtato-head")}

	namespace := k8s.NewKubeNamespace(chart, jsii.String("podtato"), &k8s.KubeNamespaceProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String("podtato"),
		},
	})
	for _, podtatoPart := range props.PodtatoParts {
		buildPodtatoHeadComponent(chart, namespace, appLabels, podtatoPart.PartName,
			podtatoPart.ImageVersion, podtatoPart.ServicePort, podtatoPart.ServiceType)
	}

	return chart
}
