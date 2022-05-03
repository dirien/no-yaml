# no-yaml

Read the detailed article on my blog:

- [Kubernetes: No YAML, please!?](https://blog.ediri.io/kubernetes-no-yaml-please)
- [Kubernetes: No YAML, please!? - Part 2](https://blog.ediri.io/kubernetes-no-yaml-please-part-2)

No YAML deployments to K8s with following approaches:

- Pulumi
- NAML
- cdk8s
- isopod

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
go build .
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

Update to the `no-yaml` project with `isopod`:

# Isopod

isopod is currently only available for MacOS and Linux systems.

```bash
wget https://github.com/cruise-automation/isopod/releases/download/v1.8.6/isopod-darwin
chmod +x isopod-darwin
mv isopod-darwin /usr/local/bin/isopod
```

To execute the deployment, just run:

```bash
cd podtato-head-isopod
isopod -kubeconfig $HOME/.kube/config install main.ipd
```

and remove the deployment with:

```bash
isopod -kubeconfig $HOME/.kube/config remove main.ipd
```

If you want to know more about [skylark](https://docs.bazel.build/versions/main/skylark/language.html).

