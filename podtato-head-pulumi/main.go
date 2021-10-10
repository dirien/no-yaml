package main

import (
	"fmt"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func buildPodtatoHeadComponent(ctx *pulumi.Context, ns *corev1.Namespace, appName pulumi.StringMap, componentName,
	imageVersion string, servicePort int, serviceType string) error {

	componentLabel := pulumi.StringMap{
		"component": pulumi.String(componentName),
	}

	_, err := appsv1.NewDeployment(ctx, componentName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(componentName),
			Namespace: ns.Metadata.Name(),
			Labels:    appName,
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: componentLabel,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: componentLabel,
				},
				Spec: &corev1.PodSpecArgs{
					TerminationGracePeriodSeconds: pulumi.Int(5),
					Containers: &corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("server"),
							Image:           pulumi.String(fmt.Sprintf("ghcr.io/podtato-head/%s:%s", componentName, imageVersion)),
							ImagePullPolicy: pulumi.String("Always"),
							Ports: &corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(9000),
								},
							},
							Env: &corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name:  pulumi.String("PORT"),
									Value: pulumi.String("9000"),
								},
							},
						}},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, componentName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(componentName),
			Namespace: ns.Metadata.Name(),
			Labels:    appName,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: componentLabel,
			Ports: &corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					Port:       pulumi.Int(servicePort),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.Int(9000),
				},
			},
			Type: pulumi.String(serviceType),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

type PodtatoHead struct {
	PodtatoParts []PodtatoParts
}

type PodtatoParts struct {
	PartName     string
	ImageVersion string
	ServicePort  int
	ServiceType  string
}

var podtatoHead = PodtatoHead{
	PodtatoParts: []PodtatoParts{
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
	pulumi.Run(func(ctx *pulumi.Context) error {
		namespace, err := corev1.NewNamespace(ctx, "podtato-kubectl", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("podtato-kubectl"),
			},
		})
		if err != nil {
			return err
		}
		appLabels := pulumi.StringMap{
			"app": pulumi.String("podtato-head"),
		}

		for _, podtatoPart := range podtatoHead.PodtatoParts {
			err = buildPodtatoHeadComponent(ctx, namespace, appLabels, podtatoPart.PartName,
				podtatoPart.ImageVersion, podtatoPart.ServicePort, podtatoPart.ServiceType)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
