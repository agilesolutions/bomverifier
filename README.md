# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.json. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent.
Let the jenkins build fail if any of the included libraries on that spring boot app are violating the compliancy of the BOM test.
## setup

* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH

## build

´´´
go build -o main .
docker build -t agilesolutions/bomverifier:latest
´´´

## run
bomverfier <URL bom.json>

## read
[check this](https://www.callicoder.com/docker-golang-image-container-example/)
