pipeline {
  agent {
    docker {
      image 'iseki0/go-github-env:latest'
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
      environment{
        GOPROXY='https://goproxy.cn,direct'
        GOOS = 'windows'
        GOARCH = 'amd64'
      }

      steps {
        sh (script: 'set', label:'Show env')
        sh (script:'go build .', label:'Go build')
      }
    }

    stage("Github release"){
      when {
//         expression { env.TAG_NAME ==~ /v.+/ }
        tag "v*"
      }
      environment{
        GIT_COMMITTER_NAME='CI.working'
        GIT_AUTHOR_NAME = 'CI.working'
        GIT_AUTHOR_EMAIL= 'working@iseki.space'
        GIT_COMMITTER_EMAIL='working@iseki.space'
      }
      steps{
        sh (script: 'set', label:'Show env')
        sh script:'''
          mkdir ~/.ssh || true
          echo "Host *\n    StrictHostKeyChecking no" > ~/.ssh/config
          chmod 0600 ~/.ssh/config
        ''', label: 'Setup SSH config'
        sh label:'Setup GitHub CLI', script:'''
           echo $GITHUB_PERSONAL_TOKEN | gh auth login --with-token
        '''
        sh script:'gh repo clone murphysec/murphysec-cli-release', label: 'Clone release repo'
        sh script:'cp murphysec-cli-simple.exe murphysec-cli-release/murphysec-cli.exe', label:'Copy files'
        dir('murphysec-cli-release'){
          sh (script:'gh release create $TAG_NAME murphysec-cli.exe --generate-notes', label:'Create release')
        }
        sh script:'''
        echo $SCOOP_REPO_DEPLOY_SSH_KEY | base64 -d > ~/.ssh/id_rsa
        chmod 600 ~/.ssh/id_rsa
        ''', label: 'Setup scoop repo ssh key'
        sh (script:'git clone git@github.com:murphysec/scoop-bucket.git', label: 'Clone scoop repo')
        dir('scoop-bucket'){
          sh script:'''
              echo "{
              \\"version\\": \\"$TAG_NAME\\",
              \\"url\\":\\"https://github.com/murphysec/murphysec-cli-release/releases/download/$TAG_NAME/murphysec-cli.exe\\",
              \\"bin\\": \\"murphysec-cli.exe\\"
              }" > bucket/murphysec-cli.json
          ''', label:"Write json"
          sh script:'''
             git add .
             git commit -am "Time: `date -Iseconds`
             Commit:$GIT_COMMIT"
             git push
          ''', label:"Commit & Push"
        }
      }
    }
  }
}