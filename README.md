# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.json. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent.
Let the jenkins build fail if any of the included libraries on that spring boot app are violating the compliancy of the BOM test.
## functionality

1. wget the BOM yaml file from github
2. go into the springboot jar file zip file and discover all libraries
3. check the compliancy againt the BOM yaml
4. report and or break off build

## setup

* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH

## build

```
go build -o main .

docker build -t agilesolutions/bomverifier:latest
```

## run
bomverfier <URL bom.json>

## include on pipeline

```
pipeline {
  agent none
  environment {
    DOCKER_IMAGE = null
  }
  stages {
    stage('Verify') {
      agent {
          docker {
              image 'agilesolutions/bomverifier:latest'
          }
      }
      steps {
        sh 'main http;//whatever.com/bom.yaml'
      }
    }
    stage('Build') {
      agent {
          docker {
              image 'maven:3-alpine'
            // do some caching on maven here
              args '-v $HOME/.m2:/root/.m2'
          }
      }
      steps {
        sh 'mvn clean install'
      }
    }
    stage('dockerbuild') {
      steps {
        script {
          DOCKER_IMAGE = docker.build("katacodarob/demo:latest")
        }
      }
    }
```


## read
[check this](https://www.callicoder.com/docker-golang-image-container-example/)
