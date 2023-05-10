package podtato

import (
	"fmt"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PodtatoHeadArgs struct {
	NamespaceName pulumi.StringInput `pulumi:"namespaceName"`
}

type PodtatoHead struct {
	pulumi.ResourceState
}

const (
	// Default values
	entryPort = 9000
	hatPort   = 9001
	leftLeg   = 9002
	leftArm   = 9003
	rightLeg  = 9004
	rightArm  = 9005
)

func NewPodtatoHead(ctx *pulumi.Context, name string, args *PodtatoHeadArgs, opts ...pulumi.ResourceOption) (*PodtatoHead, error) {
	podtatoHead := &PodtatoHead{}
	err := ctx.RegisterComponentResource("podtato:podtatohead:PodtatoHead", name, podtatoHead, opts...)
	if err != nil {
		return nil, err
	}
	namespace, err := corev1.NewNamespace(ctx, "podtato-namespace", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: args.NamespaceName,
		},
	}, pulumi.Parent(podtatoHead))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-entry", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "entry",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  entryPort,
		ServiceType:  "LoadBalancer",
		isEntry:      true,
		ServiceDiscoveryData: fmt.Sprintf(`
hat:       "http://podtato-head-hat:%d"
left-leg:  "http://podtato-head-left-leg:%d"
left-arm:  "http://podtato-head-left-arm:%d"
right-leg: "http://podtato-head-right-leg:%d"
right-arm: "http://podtato-head-right-arm:%d"
`, hatPort, leftLeg, leftArm, rightLeg, rightArm),
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-hat", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "hat",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  hatPort,
		ServiceType:  "ClusterIP",
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-left-leg", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "left-leg",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  leftLeg,
		ServiceType:  "ClusterIP",
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-left-arm", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "left-arm",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  leftArm,
		ServiceType:  "ClusterIP",
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-right-leg", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "right-leg",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  rightLeg,
		ServiceType:  "ClusterIP",
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	_, err = NewPodtatoHeadPart(ctx, "podtato-head-right-arm", &PodtatoHeadPartArgs{
		Namespace:    namespace.Metadata.Name().Elem(),
		PartName:     "right-arm",
		ImageVersion: "0.2.8-chainguard",
		ServicePort:  rightArm,
		ServiceType:  "ClusterIP",
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return nil, err
	}
	return podtatoHead, nil
}
