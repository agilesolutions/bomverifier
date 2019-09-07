# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.json. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent.
Let the jenkins build fail if any of the included libraries on that spring boot app are violating the compliancy of the BOM test.
## functionality

1. wget the BOM yaml file from github
2. go into the springboot jar file zip file and discover all libraries
3. check the compliancy againt the BOM yaml
4. report and conditionally break off

## setup

* [goto](https://www.katacoda.com/courses/docker/deploying-first-container)
* git clone https://github.com/agilesolutions/bomverifier.git
* curl -LO https://dl.google.com/go/go1.13.linux-amd64.tar.gz
* tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
* export PATH=$PATH:/usr/local/go/bin	
* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH
* go get gopkg.in/yaml.v2

## build

```
go build -o bomverify .

verify https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml false

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

1 [check this](https://www.callicoder.com/docker-golang-image-container-example/)
2 [parse yaml](https://stackoverflow.com/questions/28682439/go-parse-yaml-file/28683173)
3 [wget to file](https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go)
4 [get go package](https://gopkg.in/yaml.v2)
5 []()
6 []()
