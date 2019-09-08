pipeline {
  agent none
  environment {
    BOM_URL = "https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml"
    BREAK_BUILD = "false"
  }
  stages {
    stage('Build') {
      agent {
          docker {
              image 'kut'
          }
      }
      steps {
        sh 'ls -latr'
      }
    }
  }
}