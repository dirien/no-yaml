package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		appLabels := pulumi.StringMap{
			"app": pulumi.String("hello-world"),
		}

		config, err := corev1.NewConfigMap(ctx, "hello-world-config", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("hello-world-config"),
			},
			Data: pulumi.StringMap{
				"nginx.conf": pulumi.String(`events {}
http {
	server {
		listen 80;
		location / {
			return 200 'Hello, World!';
		}
	}
}`)},
		})
		if err != nil {
			return err
		}

		deploy, err := appsv1.NewDeployment(ctx, "hello-world-deploy", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("hello-world-deploy"),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Strategy: &appsv1.DeploymentStrategyArgs{
					Type: pulumi.String("Recreate"),
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("hello-world"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("hello-world"),
								Image: pulumi.String("nginx:latest"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String("hello-world-config"),
										MountPath: pulumi.String("/etc/nginx/"),
									},
								},
							},
						},
						Volumes: corev1.VolumeArray{
							&corev1.VolumeArgs{
								Name: pulumi.String("hello-world-config"),
								ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
									Name: config.Metadata.Elem().Name(),
									Items: corev1.KeyToPathArray{
										&corev1.KeyToPathArgs{
											Key:  pulumi.String("nginx.conf"),
											Path: pulumi.String("nginx.conf"),
										},
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = corev1.NewService(ctx, "hello-world", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("hello-world"),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type: pulumi.String("LoadBalancer"),
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(80),
					},
				},
				Selector: appLabels,
			},
		}, pulumi.DependsOn([]pulumi.Resource{deploy}))

		if err != nil {
			return err
		}

		return nil
	})
}
