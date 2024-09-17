# My stuff 

## Description
This is a simple go application that runs has 4 endpoints setup
```
/health: Returns a 200 OK status if service is healthy
/ready: Returns a 200 OK status and the service is ready to run.
/payload: Calculate a Fibonacci sequence to the random number and
return it as a JSON response with code 200.
/metrics: Returns basic metrics about the service's operation.
```
### Load Balancing K8s Load Balancer
The service is running behind a K8s loadbalancer to ensure traffic are distributed evenly as shown in this Grafana dashboard.

![Screenshot 2024-09-14 at 1 54 44 PM](https://github.com/user-attachments/assets/04a1bacb-5881-479a-b69a-09670ab7cbd9)

```PromQL sum(rate(service_requests_total[5m])) by (pod) ```

### Load Balancing Ingress-Nginx 
The service configuration can also be configured with an Ingress Nginx controller to control traffic routes. With this setup it can be more cost effective as you don't need as many elastic IPs from either GCP or AWS when you create a K8s Load Balancer.
Ingress Nginx also gives an option to either TLS termination or proxy passthrough depending on how your backend services wants to handle traffic.

* Add a line ``` 127.0.0.1 service.local ``` at the bottom of /etc/hosts *

This service is deployed on kubernetes Deployment which has rolling upgrade incrementally replacing instances of the old version with instances of the new version which will provide us with zero-downtime deployment.
Also, we are using Helm to package it all into a chart, which provides version control and rollbacks and configuration management with values.yaml. With Helm you can also switch between deploying in Production or Develpment environments by adjust values or using ```--set``` switch. The best is this can all be integrated with a CI/CD pipelines for ease of deployments.

However, in this case with minikube, as we are developing everything locally, you have to set image pull policy to Never since we don't have this image pushing to a public registry.

## Building the service
Build steps are defined in Makefile.

To build go binary
```
make build 
```
To build the container 
```
make container
```
To build Helm chart
```
make helm
```

## Local Development
We will use minikube to run this service locally with images built by docker and deploy them using helm packages.

### Setup
If you don't have minikube setup, please follow the [installation guide](https://minikube.sigs.k8s.io/docs/) for your local machine

If you don't have docker setup, please follow the [installation guide](https://docs.docker.com/engine/install/) for your local machine

If you don't have helm setup, please follow the [installation guide](https://helm.sh/docs/intro/install/) for your local machine

To start minikube run
```
minikube start
```
Once you have ran ```make helm``` and built all the artifacts you can run ```make refresh-minikube-env``` to load your local image into minikube environment and helm install your service from locally built chart.
As there is a LoadBalancer set in the k8s services.yaml you will also need to run ```minikube tunnel``` on a separate terminal so minikube can setup an external IP with localhost IP.

To execute send requests test.go run 
```
make send-requests
```

### Monitoring Setup
In local setup we'll use kube-prometheus-stack [link here](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack/) which will default provide us a set of Prometheus stack and Grafana. If there are values that need to be adjusted please use helm install / upgrade with a values.yaml.

To setup a basic kube-prometheus-stack run ```make install-monitoring``` which will install Prometheus and Grafana into minikube. 

#### Grafana
On a separate terminal you can run port-forwarding to access the UI on your browser
```
export POD_NAME=$(kubectl get po -l "app.kubernetes.io/name=grafana" -o jsonpath='{.items[0].metadata.name}')
kubectl --namespace default port-forward $POD_NAME 3000
```
On the UI to login you can get admin credentials by running.
```
kubectl get secret --namespace default prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

#### Prometheus
On a separate terminal you can run port-forwarding to access the UI on your browser
```
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=prometheus,app.kubernetes.io/instance=prometheus-kube-prometheus-prometheus" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace default port-forward $POD_NAME 9090
```
