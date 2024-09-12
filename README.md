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
The Dockerfile is defined in docker dir.
Build steps are defined in Makefile.

To build go binary
```
make build 
```

To build the container 
```
make container
```

## Local Development
We will use minikube to run this service locally with docker

### Setup
If you don't have minikube setup, please follow the [installation guide](https://minikube.sigs.k8s.io/docs/) for your local machine

If you don't have docker setup, please follow the [installation guide](https://docs.docker.com/engine/install/) for your local machine

