package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"podtato-head-pulumi/internal/podtato"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := podtato.NewPodtatoHead(ctx, "podtato-head", &podtato.PodtatoHeadArgs{
			NamespaceName: pulumi.String("podtato-head"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
