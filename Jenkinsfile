pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make build'
                sh 'make syncconfig'
                archiveArtifacts artifacts: 'build/linux/*', fingerprint: true
            }
        }
        stage('Deploy') {
            when {
                expression {
                    currentBuild.result == null || currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                sh 'make publish'
            }
        }
    }
}
