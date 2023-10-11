// ALL Scripts
pipeline {
    // 指定在哪个机器执行
    agent any

    // 声明环境变量后续使用
    environment {
        // 设置 Go 语言环境变量，可以根据你的需求进行修改
        GOBIN = "${tool 'go1.17'}/bin"
        // GOPATH = "/var/lib/jenkins/go"
    }

    stages {
        stage('git clone') {
            steps {
                checkout scmGit(branches: [[name: '${branch}']], extensions: [], userRemoteConfigs: [[credentialsId: '9c3e157d-921f-4b05-a304-78568ddd5b69', url: 'http://192.168.71.220:13888/gitlab-instance-6564fb61/dingtalk.git']])
            }
        }

        stage ('build') {
            steps {
                sh "/data/go/bin/go version"
                sh "/data/go/bin/go build"
            }
        }

        stage ('docker') {
            steps {
                sh "/usr/bin/docker build -t 192.168.71.199:80/project/dingtalk:${branch} ."
                sh "/usr/bin/docker login -u admin -p 123456 192.168.71.199:80"
                sh "/usr/bin/docker push 192.168.71.199:80/project/dingtalk:${branch}"
            }
        }
        stage ('deploy') {
            steps {
                sh "kubectl apply -f kubernetes/deployment.yaml"
            }
        }
    }

    // // jenkins 通知dingding
    // post {
    //     success {
    //         dingtalk(
    //             robot: "Jenkins-DingDing".
    //             type: 'MARKDOWN'
    //             title: "success:  ${JOB_NAME}",
    //             text: ["-成功构建:  ${JOB_NAME}!" \n- 版本: "${branch}" n- 持续时间: ${currentBuild.durationString}]
    //         )
    //     }

    //     failure {
    //         dingtalk(
    //             robot: "Jenkins-DingDing".
    //             type:  'MARKDOWN'
    //             title: "success:  ${JOB_NAME}",
    //             text: ["-成功构建: ${JOB_NAME}!" \n- 版本: "${branch}" n- 持续时间: ${currentBuild.durationString}]
    //         )
    //     }
    // }
}