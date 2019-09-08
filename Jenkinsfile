pipeline {
    agent {
        docker {
            image 'agilesolutions:bomverifier'
        }
    }
    stages {
        stage('Build') {
            steps {
                sh 'bomverifier'
            }
        }
    }
}