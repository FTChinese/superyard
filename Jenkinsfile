pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make'
            }
        }
        stage('Deploy') {
            steps {
                sh 'make publish'
            }
        }
    }
}