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

var Version = "0.2.8"

const (
	// Default values
	entryPort = 9000
	hatPort   = 9001
	leftLeg   = 9002
	leftArm   = 9003
	rightLeg  = 9004
	rightArm  = 9005
)

type PodtatoHeadApp struct {
	naml.AppMeta
	objects      []runtime.Object
	podtatoParts []podtatoParts
}

type podtatoParts struct {
	PartName             string
	ImageVersion         string
	ServicePort          int
	ServiceType          apiv1.ServiceType
	IsEntry              bool
	ServiceDiscoveryData string
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
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  entryPort,
				ServiceType:  "LoadBalancer",
				IsEntry:      true,
				ServiceDiscoveryData: fmt.Sprintf(`
hat:       "http://podtato-head-hat:%d"
left-leg:  "http://podtato-head-left-leg:%d"
left-arm:  "http://podtato-head-left-arm:%d"
right-leg: "http://podtato-head-right-leg:%d"
right-arm: "http://podtato-head-right-arm:%d"
`, hatPort, leftLeg, leftArm, rightLeg, rightArm),
			},
			{
				PartName:     "hat",
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  hatPort,
				ServiceType:  "ClusterIP",
			},
			{
				PartName:     "left-leg",
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  leftLeg,
				ServiceType:  "ClusterIP",
			},
			{
				PartName:     "left-arm",
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  leftArm,
				ServiceType:  "ClusterIP",
			},
			{
				PartName:     "right-leg",
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  rightLeg,
				ServiceType:  "ClusterIP",
			},
			{
				PartName:     "right-arm",
				ImageVersion: fmt.Sprintf("%s-chainguard", Version),
				ServicePort:  rightArm,
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

		err = p.buildPodtatoHeadComponent(ctx, client)
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

func ptrInt64(i int64) *int64 {
	return &i
}

func (p *PodtatoHeadApp) buildPodtatoHeadComponent(ctx context.Context, client kubernetes.Interface) error {
	for _, podtatoPart := range p.podtatoParts {

		labels := map[string]string{
			"app.kubernetes.io/name":      p.AppMeta.Name,
			"app.kubernetes.io/component": podtatoPart.PartName,
			"app.kubernetes.io/version":   podtatoPart.ImageVersion,
		}
		matchingLabels := map[string]string{
			"app.kubernetes.io/name":      p.AppMeta.Name,
			"app.kubernetes.io/component": podtatoPart.PartName,
		}

		name := fmt.Sprintf("podtato-head-%s", podtatoPart.PartName)

		envVarArray := []apiv1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &apiv1.EnvVarSource{
					FieldRef: &apiv1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
		}
		var volumeArray []apiv1.Volume
		var volumeMountArray []apiv1.VolumeMount

		if podtatoPart.IsEntry {
			configMapName := "podtato-head-service-discovery"
			configMap := &apiv1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: p.AppMeta.Name,
					Labels:    labels,
				},
				Data: map[string]string{
					"servicesConfig.yaml": podtatoPart.ServiceDiscoveryData,
				},
			}
			_, err := client.CoreV1().ConfigMaps(p.AppMeta.Name).Create(ctx, configMap, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("unable to create configmap in Kubernetes: %v", err)
			}
			p.objects = append(p.objects, configMap)
			volumeMountArray = append(volumeMountArray, apiv1.VolumeMount{
				Name:      configMapName,
				MountPath: "/config",
			})
			volumeArray = append(volumeArray, apiv1.Volume{
				Name: configMapName,
				VolumeSource: apiv1.VolumeSource{
					ConfigMap: &apiv1.ConfigMapVolumeSource{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: configMapName,
						},
					},
				},
			})
			envVarArray = append(envVarArray, apiv1.EnvVar{
				Name:  "SERVICES_CONFIG_FILE_PATH",
				Value: "/config/servicesConfig.yaml",
			})

		}

		deployment := &v1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.AppMeta.Name,
				Labels:    labels,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: matchingLabels,
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: matchingLabels,
					},
					Spec: apiv1.PodSpec{
						TerminationGracePeriodSeconds: ptrInt64(5),
						Containers: []apiv1.Container{
							{
								Name:            "server",
								Image:           fmt.Sprintf("ghcr.io/podtato-head/%s:%s", podtatoPart.PartName, podtatoPart.ImageVersion),
								ImagePullPolicy: apiv1.PullAlways,
								Ports: []apiv1.ContainerPort{
									{
										ContainerPort: 9000,
										Protocol:      apiv1.ProtocolTCP,
										Name:          "http",
									},
								},
								LivenessProbe: &apiv1.Probe{
									ProbeHandler: apiv1.ProbeHandler{
										TCPSocket: &apiv1.TCPSocketAction{
											Port: intstr.FromString("http"),
										},
									},
								},
								ReadinessProbe: &apiv1.Probe{
									ProbeHandler: apiv1.ProbeHandler{
										TCPSocket: &apiv1.TCPSocketAction{
											Port: intstr.FromString("http"),
										},
									},
								},
								VolumeMounts: volumeMountArray,
								Env:          envVarArray,
							},
						},
						Volumes: volumeArray,
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
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.AppMeta.Name,
				Labels:    labels,
			},
			Spec: apiv1.ServiceSpec{
				Selector: matchingLabels,
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
