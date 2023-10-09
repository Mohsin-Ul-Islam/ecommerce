#!/usr/bin/env bash

if [[ ! $(pgrep -x docker) ]]; then
    echo 'you need docker daemon running to run the project locally'
    exit 1
fi

if [[ ! $(which minikube) ]]; then
    echo 'you need minikube to run the project locally'
    exit 1
fi

if [[ ! $(which bazel) ]]; then
    echo 'you need bazel to run the project locally'
    exit 1
fi

if [[ ! $(which grpcurl) ]]; then
    echo 'you need grpcurl to run the project locally'
    exit 1
fi

if [[ ! $(which migrate) ]]; then
    echo 'you need golang migrate to run the project locally'
    exit 1
fi

if [[ ! $(pgrep -x minikube) ]]; then
    echo 'starting minikube'
    minikube start
fi

echo 'creating the temporal deployment'
minikube kubectl -- create deployment temporal --image=temporalio/server:latest

echo 'creating the cockroach deployment'
minikube kubectl -- apply -f database/deploy/deployment.yaml

echo 'creating the cockroach service'
minikube kubectl -- apply -f database/deploy/service.yaml

echo 'port forwarding the cockroach service'
minikube kubectl -- port-forward deployments/database-deployment 26257:26257 > /dev/null 2>&1 &

echo 'running the migrations'
migrate -path database/migrations -database 'cockroachdb://root:@localhost:26257/defaultdb?sslmode=disable' up

echo 'seeding the data'
cockroach sql --insecure --host localhost --port 26257 < data/seed.sql

echo 'building the customers docker image'
bazel run //customers/cmd/serve:image

echo 'loading the customers image in minikube'
minikube image load bazel/customers/cmd/serve:image

echo 'creating the customers deployment'
minikube kubectl -- apply -f customers/deploy/deployment.yaml

echo 'creating the customers service'
minikube kubectl -- apply -f customers/deploy/service.yaml

echo 'port forwarding the customers service'
minikube kubectl -- port-forward deployments/customers-deployment 8080:8080 > /dev/null 2>&1 &

echo 'testing the connection to customers service'
if [[ ! $(grpcurl -plaintext localhost:8080 list) ]]; then
    echo 'customers service not deployed successfully'
fi

