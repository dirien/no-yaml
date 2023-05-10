package podtato

import (
	"fmt"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2/k8s"
)

type PodtatoHeadProps struct {
	cdk8s.ChartProps
	PodtatoParts  []PodtatoPartArgs
	NamespaceName string
}

type PodtatoPartArgs struct {
	PartName             string
	ImageVersion         string
	ServicePort          int
	ServiceType          string
	IsEntry              bool
	ServiceDiscoveryData string
}

func buildPodtatoHeadComponent(scope constructs.Construct, ns k8s.KubeNamespace, name string, args PodtatoPartArgs) {
	labels := map[string]*string{
		"app.kubernetes.io/name":      jsii.String(name),
		"app.kubernetes.io/component": jsii.String(args.PartName),
		"app.kubernetes.io/version":   jsii.String(args.ImageVersion),
	}
	matchedLabels := map[string]*string{
		"app.kubernetes.io/name":      jsii.String(name),
		"app.kubernetes.io/component": jsii.String(args.PartName),
	}

	k8sName := fmt.Sprintf("podtato-head-%s", args.PartName)

	envVarArray := []*k8s.EnvVar{
		&k8s.EnvVar{
			Name: jsii.String("POD_NAME"),
			ValueFrom: &k8s.EnvVarSource{
				FieldRef: &k8s.ObjectFieldSelector{
					FieldPath: jsii.String("metadata.name"),
				},
			},
		},
	}

	var volumeArray []*k8s.Volume
	var volumeMountArray []*k8s.VolumeMount

	if args.IsEntry {
		configMapName := "podtato-head-service-discovery"
		k8s.NewKubeConfigMap(scope, jsii.String("podtato-head-config"), &k8s.KubeConfigMapProps{
			Metadata: &k8s.ObjectMeta{
				Name:      jsii.String(configMapName),
				Namespace: ns.Metadata().Name(),
				Labels:    &labels,
			},
			Data: &map[string]*string{
				"servicesConfig.yaml": jsii.String(args.ServiceDiscoveryData),
			},
		})

		volumeArray = append(volumeArray, &k8s.Volume{
			Name: jsii.String(configMapName),
			ConfigMap: &k8s.ConfigMapVolumeSource{
				Name: jsii.String(configMapName),
			},
		})

		volumeMountArray = append(volumeMountArray, &k8s.VolumeMount{
			Name:      jsii.String(configMapName),
			MountPath: jsii.String("/config"),
		})

		envVarArray = append(envVarArray, &k8s.EnvVar{
			Name:  jsii.String("SERVICES_CONFIG_FILE_PATH"),
			Value: jsii.String("/config/servicesConfig.yaml"),
		})

	}

	k8s.NewKubeDeployment(scope, jsii.String(fmt.Sprintf("%s-dep", args.PartName)), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String(k8sName),
			Namespace: ns.Metadata().Name(),
			Labels:    &labels,
		},
		Spec: &k8s.DeploymentSpec{
			Selector: &k8s.LabelSelector{
				MatchLabels: &matchedLabels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &matchedLabels,
				},
				Spec: &k8s.PodSpec{
					TerminationGracePeriodSeconds: jsii.Number(5),
					Containers: &[]*k8s.Container{
						{
							Name:            jsii.String("server"),
							Image:           jsii.String(fmt.Sprintf("ghcr.io/podtato-head/%s:%s", args.PartName, args.ImageVersion)),
							ImagePullPolicy: jsii.String("Always"),
							Ports: &[]*k8s.ContainerPort{
								{
									ContainerPort: jsii.Number(9000),
									Protocol:      jsii.String("TCP"),
									Name:          jsii.String("http"),
								},
							},
							LivenessProbe: &k8s.Probe{
								TcpSocket: &k8s.TcpSocketAction{
									Port: k8s.IntOrString_FromString(jsii.String("http")),
								},
							},
							ReadinessProbe: &k8s.Probe{
								TcpSocket: &k8s.TcpSocketAction{
									Port: k8s.IntOrString_FromString(jsii.String("http")),
								},
							},
							Env:          &envVarArray,
							VolumeMounts: &volumeMountArray,
						},
					},
					Volumes: &volumeArray,
				},
			},
		},
	})

	k8s.NewKubeService(scope, jsii.String(fmt.Sprintf("%s-svc", args.PartName)), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String(k8sName),
			Namespace: ns.Metadata().Name(),
			Labels:    &labels,
		},
		Spec: &k8s.ServiceSpec{
			Selector: &matchedLabels,
			Ports: &[]*k8s.ServicePort{
				{
					Name:       jsii.String("http"),
					Port:       jsii.Number(float64(args.ServicePort)),
					Protocol:   jsii.String("TCP"),
					TargetPort: k8s.IntOrString_FromNumber(jsii.Number(float64(9000))),
				},
			},
			Type: jsii.String(args.ServiceType),
		},
	})
}

func PodtatoHeadChart(scope constructs.Construct, id string, props *PodtatoHeadProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	namespace := k8s.NewKubeNamespace(chart, jsii.String("podtato"), &k8s.KubeNamespaceProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String(props.NamespaceName),
		},
	})
	for _, podtatoPart := range props.PodtatoParts {
		fmt.Println(podtatoPart)
		buildPodtatoHeadComponent(chart, namespace, id, podtatoPart)
	}

	return chart
}
