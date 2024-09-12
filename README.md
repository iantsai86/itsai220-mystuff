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
As there is a LoadBalancer set in the k8s services.yaml you will also need to run ```minikube tunnel``` on a separate terminal