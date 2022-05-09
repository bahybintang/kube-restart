<h1>Kube Restart</h1>

<h4>Helps you periodically restart your k8s deployments!</h4>

<p>
  <a href="https://github.com/kubernetes/kubernetes/releases">
    <img src="https://img.shields.io/badge/Kubernetes-%3E%3D%201.18-brightgreen" alt="kubernetes">
  </a>
  <a href="https://golang.org/doc/go1.18">
    <img src="https://img.shields.io/github/go-mod/go-version/bahybintang/kube-restart?color=blueviolet" alt="go-version">
  </a>
</p>

<div>
<hr>
</div>


## Project Summary

This project helps to manage periodically restart deployment job for Kubernetes. It all centralized using one config or via annotations so you don't have to configure cronjob for every deployments manually. This project can work with either deployment or statefulsets.

## Kube Restart Config

You can either apply the config using config file or via annotations on the deployments or statefulsets.

### Via Config File

The config file is pretty simple consisting the list of deployments and statefulsets.

```yaml
deployments:
- name: demo-deployment
  namespace: default
  schedule: "* * * * *"
- name: ingress-nginx-controller
  namespace: ingress-nginx
  schedule: "* * * * *"

statefulsets:
- name: demo-statefulsets-redis-replicas
  namespace: default
  schedule: "*/5 * * * *"
```

### Via Kubernetes Annotations

You can add `kube.restart/schedule` annotation to the deployments or statefulsets to configure the restart schedule.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kube.restart/schedule: '*/1 * * * *'
  name: demo-deployment
  namespace: default
...
```

## Build, Run, and Deploy

### How to Build

```bash
make build
```

### How to Run Outside Cluster

```bash
make build
CONFIG_PATH=<CONFIG_PATH> ./kube-restart  --kubeconfig=<KUBECONFIG_PATH>
```

### Install in Cluster

```bash
helm repo add kube-restart https://bahybintang.github.io/kube-restart
helm install kube-restart kube-restart/kube-restart
```
