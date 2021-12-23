pipeline {
  agent any

  stages {
    stage("Checkout") {
      steps {
        checkout(
          [$class: 'GitSCM',
          branches: [[name: GIT_BUILD_REF]],
          userRemoteConfigs: [[
            url: GIT_REPO_URL,
              credentialsId: CREDENTIALS_ID
            ]]]
        )
      }
    }

    stage("Go build"){
      steps{
        sh 'GOOS=windows GOARCH=amd64 go build .'
        sh 'ls'
      }
    }

  }
}