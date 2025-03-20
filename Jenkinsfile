pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
        NODE_VERSION = '18'
    }
    
    stages {
        stage('Test') {
            steps {
                sh 'go test ./... -v'
                sh 'npm run test:e2e'
            }
        }
        
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
        
        stage('Deploy') {
            steps {
                sh 'make deploy'
            }
        }
        
        stage('Monitor') {
            steps {
                sh 'curl -X POST ${MONITORING_WEBHOOK}'
            }
        }
    }
    
    post {
        failure {
            script {
                rollback()
            }
        }
    }
}