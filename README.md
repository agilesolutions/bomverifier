# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.txt. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent.
Let the jenkins build fail if any of the included libraries on that spring boot app are violating the compliancy of the BOM test.
## functionality

1. wget the BOM txt file from github
2. go into the springboot jar file zip file and discover all libraries
3. check the compliancy againt the BOM txt
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

## build

```
go build -o bomverifier .

bomverifier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate

docker build -t agilesolutions/bomverifier:latest .
```

## GO build and Docker build through Multistage Docker builds
GO compile this code and produce a Docker image by simply running

```
docker build -t agilesolutions/bomverifier:latest .
```

## run
bomverfier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate

## where to find the Springboot BOM details and release trains

* [Tutorial on Maven BOM and Spring](https://www.baeldung.com/spring-maven-bom)
* [Ablout Spring Boot and Clou BOM and Release Trains](https://spring.io/blog/2019/07/24/simplifying-the-spring-cloud-release-train)
* [spring-boot-dependencies BOM](https://github.com/mahendra-shinde/maven-repo-springboot/blob/master/repository/org/springframework/spring-framework-bom/5.1.8.RELEASE/spring-framework-bom-5.1.8.RELEASE.pom)

## now run this docker agent on a jenkins pipeline, lets spin up jenkins

* [go to katacoda](https://www.katacoda.com/courses/kubernetes/helm-package-manager)
* create directory /jenkins
* docker run -d --name jenkins --user root --privileged=true -p 8080:8080 -v /jenkins:/var/jenkins_home -v /var/run/docker.sock:/var/run/docker.sock jenkinsci/blueocean
* docker logs -f jenkins
* docker exec -ti jenkins bash
* docker ps -a
* browse to http://localhost:8080 and wait until the Unlock Jenkins page appears.
* get password from /jenkins/secrets/initialAdminPassword
* create new pipeline job from https://github.com/agilesolutions/bomverifier.git

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
        sh 'bomverifier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate'
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

* [check this](https://www.callicoder.com/docker-golang-image-container-example/)
* [parse yaml](https://stackoverflow.com/questions/28682439/go-parse-yaml-file/28683173)
* [wget to file](https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go)
* [get go package](https://gopkg.in/yaml.v2)
* [jenkins pipelines and docker agents](https://jenkins.io/doc/book/pipeline/docker/)
