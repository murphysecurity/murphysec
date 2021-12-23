pipeline {
  agent {
    docker {
      image 'golang:latest'
      reuseNode 'true'
    }

  }
  stages {
    stage('Checkout') {
      steps {
        checkout([$class: 'GitSCM',
        branches: [[name: GIT_BUILD_REF]],
        userRemoteConfigs: [[
          url: GIT_REPO_URL,
          credentialsId: CREDENTIALS_ID
        ]]])
      }
    }
    stage('Go build') {
      steps {
        sh 'go env -w GOPROXY=https://goproxy.cn,direct'
        sh 'GOOS=windows GOARCH=amd64 go build .'
        sh 'ls'
      }
    }
  }
}