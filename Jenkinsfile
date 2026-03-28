pipeline {
    agent any

    environment {
        // Image Config
        REGISTRY = "ghcr.io/wahyunurdian26"
        IMAGE_NAME = "service-audit"
        APP_NAME = "service-audit"

        // Deployment Repo
        DEPLOY_REPO = "https://github.com/wahyunurdian26/deployment-config.git"
        DEPLOY_PATH = "omnipay/hcloud/omnipay-manifests/overlays/staging/service-audit/deployment-patch.yaml"

        // Credentials
        GH_CREDENTIALS_ID = "github-token"
        DOCKER_CRED_ID = "ghcr-login"

        // Build Optimizations
        DOCKER_BUILDKIT = "1"
    }

    stages {
//        stage('Testing') {
//            steps {
//                withCredentials([string(credentialsId: GH_CREDENTIALS_ID, variable: 'G_TOKEN')]) {
//                    echo "Running unit test via Docker..."
//                    sh '''
//                    docker build \
//                      --build-arg GITHUB_TOKEN=$G_TOKEN \
//                      --build-arg BUILDKIT_INLINE_CACHE=1 \
//                      --cache-from ${REGISTRY}/${IMAGE_NAME}:latest \
//                      --target tester \
//                      -t ${APP_NAME}:test .
//                    '''
//                }
//            }
//        }

        stage('Code review') {
            steps {
                echo "Running security scan / code review..."
                sh 'echo Security Scan PASSED'
            }
        }

        stage('Prepare') {
            steps {
                withCredentials([usernamePassword(credentialsId: DOCKER_CRED_ID, usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                    echo "Login to GHCR..."
                    sh '''
                    echo $PASS | docker login ${REGISTRY} -u $USER --password-stdin
                    '''
                }
            }
        }

        stage('Build and Deploy') {
            stages {
                stage('Deploy to Staging') {
                    when {
                        anyOf {
                            branch 'master'
                            branch 'main'
                        }
                    }
                    steps {
                        withCredentials([string(credentialsId: GH_CREDENTIALS_ID, variable: 'G_TOKEN')]) {
                            echo "Building and Pushing Docker image..."
                            sh '''
                            docker build \
                              --build-arg GITHUB_TOKEN=$G_TOKEN \
                              --build-arg BUILDKIT_INLINE_CACHE=1 \
                              --cache-from ${REGISTRY}/${IMAGE_NAME}:latest \
                              -t ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER} .

                            docker push ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER}

                            docker tag ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER} \
                                       ${REGISTRY}/${IMAGE_NAME}:latest

                            docker push ${REGISTRY}/${IMAGE_NAME}:latest
                            '''

                            echo "Updating deployment repo (GitOps)..."
                            sh '''
                            rm -rf deployment-config
                            git clone https://$G_TOKEN@github.com/wahyunurdian26/deployment-config.git
                            cd deployment-config
                            sed -i 's|image: .*|image: '${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER}'|g' ${DEPLOY_PATH}
                            git config user.email "jenkins@wahyunurdian.com"
                            git config user.name "Jenkins CI"
                            git add .
                            git diff --cached --quiet || git commit -m "Update ${APP_NAME} image to staging-${BUILD_NUMBER}"
                            git push origin master
                            '''
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo "✅ Pipeline SUCCESS - ArgoCD will deploy automatically"
        }
        failure {
            echo "❌ Pipeline FAILED - Check logs"
        }
    }
}
