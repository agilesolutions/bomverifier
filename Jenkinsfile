pipeline {
  agent none
  environment {
    URL = "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml"
    BREAK-BUILD = "false"
  }
  stages {
    stage('Build') {
      agent {
          docker {
              image 'agilesolutions/bomverifier'
          }
      }
      steps {
        sh 'ls -latr'
      }
    }
  }
}