# no-yaml

No YAML deployments to K8s with following approaches:

- Pulumi
- NAML
- cdk8s

We will deploy the 📨🚚 CNCF App Delivery SIG Demo [podtato-head](https://github.com/podtato-head/podtato-head)

and use this manifest as a base -> https://github.com/podtato-head/podtato-head/blob/main/delivery/kubectl/manifest.yaml

# Pulumi

```
pulumi new
```

And chose the `kubernetes-go` template

# NAML

Create a new Go project and add the NAML libs to your project.

Run go build . to crate your app (or use goreleaser)

```
go buid .
```

Run your app. ./app

# cdk8s

```
npm install -g cdk8s-cli #for the 1.0.0-beta

brew install cdk8s # for the 0.33.0
```

Then create your project with

```
mkdir podtato-head-cdk8s
cd podtato-head-cdk8s
cdk8s init go-app
cdk8s import
cdk8s synth
kubectl apply -f dist/podtato-head-cdk8s.k8s.yaml 
```