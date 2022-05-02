package podtato

import (
	"context"
	"fmt"
	"github.com/kris-nova/naml"
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
	naml.AppMeta
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
		AppMeta: naml.AppMeta{
			Description: "📨🚚 CNCF App Delivery SIG Demo",
			ObjectMeta: metav1.ObjectMeta{
				Name:            "podtato",
				ResourceVersion: Version,
			},
		},
		podtatoParts: []podtatoParts{
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
}

func (p *PodtatoHeadApp) Install(client kubernetes.Interface) error {
	ctx := context.Background()
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.AppMeta.Name,
		},
	}
	if client != nil {
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
	}
	return nil
}

func (p *PodtatoHeadApp) Uninstall(client kubernetes.Interface) error {
	ctx := context.Background()
	err := client.CoreV1().Namespaces().Delete(ctx, p.AppMeta.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (p *PodtatoHeadApp) Meta() *naml.AppMeta {
	return &p.AppMeta
}

func (p *PodtatoHeadApp) Objects() []runtime.Object {
	return p.objects
}

func (p *PodtatoHeadApp) buildPodtatoHeadComponent(ctx context.Context, client kubernetes.Interface, appName map[string]string) error {
	for _, podtatoPart := range p.podtatoParts {
		name := fmt.Sprintf("podtato-head-%s", podtatoPart.PartName)
		componentLables := map[string]string{"component": podtatoPart.PartName}
		terminationGracePeriodSeconds := int64(5)
		deployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.AppMeta.Name,
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
		_, err := client.AppsV1().Deployments(p.AppMeta.Name).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to install deployment in Kubernetes: %v", err)
		}
		p.objects = append(p.objects, deployment)

		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.AppMeta.Name,
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
		_, err = client.CoreV1().Services(p.AppMeta.Name).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to install service in Kubernetes: %v", err)
		}
		p.objects = append(p.objects, service)
	}
	return nil
}
