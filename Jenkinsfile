#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.44.3', 'github.com/cloudogu/dogu-build-lib@v1.2.0', 'github.com/cloudogu/zalenium-build-lib@v2.0.0'])
import com.cloudogu.ces.cesbuildlib.*
import com.cloudogu.ces.dogubuildlib.*
import com.cloudogu.ces.zaleniumbuildlib.*

node('vagrant') {
    projectName = 'go-redmine'
    project = "github.com/cloudogu/${projectName}"
    branch = "${env.BRANCH_NAME}"

    timestamps {
        stage('Checkout') {
            checkout scm
        }

        new Docker(this)
            .image('golang:1.14.13')
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
        {

            stage('Build') {
                make 'clean package checksum'
                archiveArtifacts 'target/**/*.deb'
                archiveArtifacts 'target/**/*.sha256sum'
            }

            stage('Unit Test') {
                make 'unit-test'
                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
            }

            stage('Static Analysis') {
                def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
                    withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=cloudogu", "CI_REPO_NAME=${projectName}"]) {
                        make 'static-analysis-ci'
                    }
                }
            }
        }

        stage('SonarQube') {
            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
            withSonarQubeEnv {
                Git git = new Git(this, "cesmarvin")

                sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
                gitWithCredentials("fetch --all")

                if (branch == "master") {
                    echo "This branch has been detected as the master branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName}"
                } else if (branch == "develop") {
                    echo "This branch has been detected as the develop branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=master  "
                } else if (env.CHANGE_TARGET) {
                    echo "This branch has been detected as a pull request."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} -Dsonar.branch.name=${env.CHANGE_BRANCH}-PR${env.CHANGE_ID} -Dsonar.branch.target=${env.CHANGE_TARGET}    "
                } else if (branch.startsWith("feature/")) {
                    echo "This branch has been detected as a feature branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
                } else if (branch.startsWith("bugfix/")) {
                    echo "This branch has been detected as a bugfix branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
                }
            }
            timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
                def qGate = waitForQualityGate()
                if (qGate.status != 'OK') {
                    unstable("Pipeline unstable due to SonarQube quality gate failure")
                }
            }
        }

    }
}

String project
String projectName
String branch

void make(goal) {
    sh "cd /go/src/${project} && make ${goal}"
}

void gitWithCredentials(String command) {
    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
        sh(
                script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
                returnStdout: true
        )
    }
}
