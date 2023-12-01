# KubeBuilder split files

This very minimal project helps with the following sceraio:

When a user is developing a kubebuilder based project, there are make targets to use kustomize to generate the manifests.
If the user wants to have an easy way to maintain a helm chart based on the output of those manifests, this would be the tasks:

1. create a helm chart and its basic directory structure
2. use kustomize to generate the manifests
3. split the output into single files
4. change the templated values (such as the image name) in the manifests to a helm template variable

Assuming 1. is done by the user and 2. is the standard function of a kubebuilder project, this project helps with 3. and 4.

## Usage

Assume you created a very basic helm chart with the following structure:

```
charts
└── my-operator
    ├── Chart.yaml
    ├── templates
    │   ├── NOTES.txt
    └── values.yaml
```

Then you can add the following new targets to your Makefile:

```
get-split-files:
	go install github.com/syseleven/kubebuilder-split-files@latest

generate-chart: kustomize get-split-files
	$(KUSTOMIZE) build config/default | ~/go/bin/kubebuilder-split-files --chart-dir ./charts --operator-name my-operator --image registry.example.com/orga/my-operator:1.2.3
```

NOTE: change the parameters according to what is generated by kubebuilder. See `config/manager/kustomization.yaml` in your project.


## How it works

The tool uses the kustomize API to read the manifests and then splits them into single files. It also replaces the image name with a helm template variable.
