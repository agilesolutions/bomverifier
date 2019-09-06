# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.json. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent
## setup

* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH

## run
bomverfier <URL bom.json>

## read
[check this](https://www.callicoder.com/docker-golang-image-container-example/)
