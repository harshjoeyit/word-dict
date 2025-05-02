pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'go build -o word-dict main.go'
                stash name: 'binary', includes: 'word-dict'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sh 'go test ./... -v'
                sh 'go vet ./...'
                sh 'golint ./...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                unstash 'binary'
                sh './word-dict'  // Run the binary
                // sh 'chmod +x deploy.sh && ./deploy.sh'  // Make script executable and run it
            }
        }
    }
}