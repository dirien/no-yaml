package podtato

import (
	"fmt"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PodtatoHeadPartArgs struct {
	Namespace            pulumi.StringInput `pulumi:"namespace"`
	PartName             string             `pulumi:"partName"`
	ImageVersion         string             `pulumi:"imageVersion"`
	ServicePort          int                `pulumi:"servicePort"`
	ServiceType          string             `pulumi:"serviceType"`
	skipAwait            bool               `pulumi:"skipAwait"`
	isEntry              bool               `pulumi:"isEntry"`
	ServiceDiscoveryData string             `pulumi:"serviceDiscoveryData"`
}

type PodtatoHeadPart struct {
	pulumi.ResourceState
}

func NewPodtatoHeadPart(ctx *pulumi.Context, name string, args *PodtatoHeadPartArgs, opts ...pulumi.ResourceOption) (*PodtatoHeadPart, error) {
	podtatoHeadPart := &PodtatoHeadPart{}
	err := ctx.RegisterComponentResource("pkg:index:PodtatoHeadPart", name, podtatoHeadPart, opts...)
	if err != nil {
		return nil, err
	}
	labels := pulumi.StringMap{
		"app.kubernetes.io/name":      pulumi.String(name),
		"app.kubernetes.io/component": pulumi.String(args.PartName),
		"app.kubernetes.io/version":   pulumi.String(args.ImageVersion),
	}
	matchedLabels := pulumi.StringMap{
		"app.kubernetes.io/name":      pulumi.String(name),
		"app.kubernetes.io/component": pulumi.String(args.PartName),
	}

	k8sName := fmt.Sprintf("podtato-head-%s", args.PartName)

	annotations := pulumi.StringMap{
		"pulumi.com/skipAwait": pulumi.Sprintf("%t", args.skipAwait),
	}

	var volumeArray *corev1.VolumeArray
	var volumeMountArray *corev1.VolumeMountArray

	envVarArray := &corev1.EnvVarArray{
		&corev1.EnvVarArgs{
			Name: pulumi.String("POD_NAME"),
			ValueFrom: &corev1.EnvVarSourceArgs{
				FieldRef: &corev1.ObjectFieldSelectorArgs{
					FieldPath: pulumi.String("metadata.name"),
				},
			},
		},
	}

	if args.isEntry {
		configMapName := "podtato-head-service-discovery"
		config, err := corev1.NewConfigMap(ctx, "podtato-head-config", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(configMapName),
				Labels:    labels,
				Namespace: args.Namespace,
			},
			Data: pulumi.StringMap{
				"servicesConfig.yaml": pulumi.String(args.ServiceDiscoveryData),
			},
		}, pulumi.Parent(podtatoHeadPart))
		if err != nil {
			return nil, err
		}
		volumeMountArray = &corev1.VolumeMountArray{
			&corev1.VolumeMountArgs{
				Name:      pulumi.String(configMapName),
				MountPath: pulumi.String("/config"),
			},
		}
		volumeArray = &corev1.VolumeArray{
			&corev1.VolumeArgs{
				Name: config.Metadata.Name().Elem(),
				ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
					Name: config.Metadata.Name().Elem(),
				},
			},
		}
		envVarArray = &corev1.EnvVarArray{
			envVarArray.ToEnvVarArrayOutput().Index(pulumi.Int(0)),
			&corev1.EnvVarArgs{
				Name:  pulumi.String("SERVICES_CONFIG_FILE_PATH"),
				Value: pulumi.String("/config/servicesConfig.yaml"),
			},
		}

	}
	podSpec := &corev1.PodSpecArgs{
		TerminationGracePeriodSeconds: pulumi.Int(5),
		Containers: &corev1.ContainerArray{
			&corev1.ContainerArgs{
				Name:            pulumi.String("server"),
				Image:           pulumi.String(fmt.Sprintf("ghcr.io/podtato-head/%s:%s", args.PartName, args.ImageVersion)),
				ImagePullPolicy: pulumi.String("Always"),
				Ports: &corev1.ContainerPortArray{
					&corev1.ContainerPortArgs{
						ContainerPort: pulumi.Int(9000),
						Protocol:      pulumi.String("TCP"),
						Name:          pulumi.String("http"),
					},
				},
				LivenessProbe: &corev1.ProbeArgs{
					TcpSocket: &corev1.TCPSocketActionArgs{
						Port: pulumi.String("http"),
					},
				},
				ReadinessProbe: &corev1.ProbeArgs{
					TcpSocket: &corev1.TCPSocketActionArgs{
						Port: pulumi.String("http"),
					},
				},
				VolumeMounts: volumeMountArray,
				Env:          envVarArray,
			}},
		Volumes: volumeArray,
	}

	_, err = appsv1.NewDeployment(ctx, args.PartName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(k8sName),
			Namespace:   args.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: matchedLabels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: matchedLabels,
				},
				Spec: podSpec,
			},
		},
	}, pulumi.Parent(podtatoHeadPart))
	if err != nil {
		return nil, err
	}

	_, err = corev1.NewService(ctx, args.PartName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(k8sName),
			Namespace:   args.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: matchedLabels,
			Ports: &corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					Port:       pulumi.Int(args.ServicePort),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.Int(9000),
				},
			},
			Type: pulumi.String(args.ServiceType),
		},
	}, pulumi.Parent(podtatoHeadPart))
	if err != nil {
		return nil, err
	}

	return podtatoHeadPart, nil
}
