Please note that the below steps can be automated and optimized using a shell script.

Local Development Environment (using script)
==========================================

``chmod +x ./scripts/run.sh && ./scripts/run.sh``

Local Development Environment (no minikube)
==========================================

Setting Up the Database
-----------------------

First of all run a cockroach cluster locally

``cockroach start-single-node --advertise-addr='localhost' --insecure``

Run migrations

``migrate -database "cockroachdb://root@localhost:26257/defaultdb?sslmode=disable" -path migrations up``

Seed the initial data into database

``cockroach sql --insecure < data/seed.sql``

Setting up a task queue service
-------------------------------

Run a temporal server locally

``temporal server start-dev``

Setting up the ecommerce services
---------------------------------

Run the services in order.

``bazel run //customers/cmd/serve:serve``

``bazel run //catalogue/cmd/serve:serve``

``bazel run //transactions/cmd/serve:serve``

Example requests are in a file named `requests.http`.

Local Development Environment (with minikube)
===========================================

Create docker images of the services

``bazel run //customers/cmd/serve:image``

``bazel run //catalogue/cmd/serve:image``

``bazel run //transactions/cmd/serve:image``

Start minikube

``minikube start``

Create cockroachdb

``minikube kubectl -- create deployment cockroachdb --image=cockroachdb/cockroach:v23.1.8``

Create temporal server

``minikube kubectl -- create deployment temporal --image=temporalio/server:latest``


Push ecommerce images to minikube (built in previous section)

``minikube image load bazel/customers/cmd/serve:image``

``minikube image load bazel/catalogue/cmd/serve:image``

``minikube image load bazel/transactions/cmd/serve:image``

Create deployments

``minikube kubectl -- apply -f customers/deploy/deployment.yaml``

``minikube kubectl -- apply -f catalogue/deploy/deployment.yaml``

``minikube kubectl -- apply -f transactions/deploy/deployment.yaml``

Create services

``minikube kubectl -- apply -f customers/deploy/service.yaml``

``minikube kubectl -- apply -f catalogue/deploy/service.yaml``

``minikube kubectl -- apply -f transactions/deploy/service.yaml``

