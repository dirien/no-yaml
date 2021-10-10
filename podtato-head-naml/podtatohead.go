package podtato

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Version = "0.0.1"

type PodtatoHeadApp struct {
	name         string
	description  string
	objects      []runtime.Object
	podtatoParts []podtatoParts
}

type podtatoParts struct {
	PartName     string
	ImageVersion string
	ServicePort  int
	ServiceType  apiv1.ServiceType
}

func NewPodtatoHeadApp() *PodtatoHeadApp {
	return &PodtatoHeadApp{
		name:        "podtato-kubectl",
		description: "📨🚚 CNCF App Delivery SIG Demo",
		podtatoParts: []podtatoParts{
			{
				PartName:     "podtato-main",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9000,
				ServiceType:  apiv1.ServiceTypeLoadBalancer,
			},
			{
				PartName:     "podtato-hats",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9001,
				ServiceType:  apiv1.ServiceTypeClusterIP,
			},
			{
				PartName:     "podtato-left-leg",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9002,
				ServiceType:  apiv1.ServiceTypeClusterIP,
			},
			{
				PartName:     "podtato-left-arm",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9003,
				ServiceType:  apiv1.ServiceTypeClusterIP,
			},
			{
				PartName:     "podtato-right-leg",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9004,
				ServiceType:  apiv1.ServiceTypeClusterIP,
			},
			{
				PartName:     "podtato-right-arm",
				ImageVersion: "v1-latest-dev",
				ServicePort:  9005,
				ServiceType:  apiv1.ServiceTypeClusterIP,
			},
		},
	}
}

func (p *PodtatoHeadApp) Install(client *kubernetes.Clientset) error {
	ctx := context.Background()
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.name,
		},
	}
	ns, err := client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to install namespace in Kubernetes: %v", err)
	}
	p.objects = append(p.objects, ns)
	appLabels := map[string]string{"app": "podtato-head"}

	err = p.buildPodtatoHeadComponent(ctx, client, appLabels)
	if err != nil {
		return err
	}

	return nil
}

func (p *PodtatoHeadApp) Uninstall(client *kubernetes.Clientset) error {
	ctx := context.Background()
	err := client.CoreV1().Namespaces().Delete(ctx, p.name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (p *PodtatoHeadApp) Meta() *metav1.ObjectMeta {
	return &metav1.ObjectMeta{
		Name: p.name,
	}
}

func (p *PodtatoHeadApp) Description() string {
	return p.description
}

func (p *PodtatoHeadApp) Objects() []runtime.Object {
	return p.objects
}

func (p *PodtatoHeadApp) buildPodtatoHeadComponent(ctx context.Context, client *kubernetes.Clientset, appName map[string]string) error {

	for _, podtatoPart := range p.podtatoParts {
		componentLables := map[string]string{"component": podtatoPart.PartName}
		terminationGracePeriodSeconds := int64(5)
		deployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podtatoPart.PartName,
				Namespace: p.name,
				Labels:    appName,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: componentLables,
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: componentLables,
					},
					Spec: apiv1.PodSpec{
						TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
						Containers: []apiv1.Container{
							{
								Name:            "server",
								Image:           fmt.Sprintf("ghcr.io/podtato-head/%s:%s", podtatoPart.PartName, podtatoPart.ImageVersion),
								ImagePullPolicy: apiv1.PullAlways,
								Ports: []apiv1.ContainerPort{
									{
										ContainerPort: 9000,
									},
								},
								Env: []apiv1.EnvVar{
									{
										Name:  "PORT",
										Value: "9000",
									},
								},
							},
						},
					},
				},
			},
		}
		_, err := client.AppsV1().Deployments(p.name).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to install deployment in Kubernetes: %v", err)
		}
		p.objects = append(p.objects, deployment)

		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podtatoPart.PartName,
				Namespace: p.name,
				Labels:    appName,
			},
			Spec: apiv1.ServiceSpec{
				Selector: componentLables,
				Ports: []apiv1.ServicePort{
					{
						Name:       "http",
						Port:       int32(podtatoPart.ServicePort),
						Protocol:   apiv1.ProtocolTCP,
						TargetPort: intstr.FromInt(9000),
					},
				},
				Type: podtatoPart.ServiceType,
			},
		}
		_, err = client.CoreV1().Services(p.name).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to install service in Kubernetes: %v", err)
		}
		p.objects = append(p.objects, service)
	}
	return nil
}
