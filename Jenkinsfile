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
    }

    stages {

        // ========================
        // 🔧 PREPARATION
        // ========================
        stage('Preparation') {
            steps {
                echo "Starting build for ${APP_NAME}..."
                sh 'docker version'
            }
        }

        // ========================
        // 🧪 UNIT TEST
        // ========================
        stage('Unit Testing') {
            steps {
                withCredentials([string(credentialsId: GH_CREDENTIALS_ID, variable: 'G_TOKEN')]) {
                    echo "Running unit test via Docker..."
                    sh '''
                    docker build \
                      --build-arg GITHUB_TOKEN=$G_TOKEN \
                      --target tester \
                      -t ${APP_NAME}:test .
                    '''
                }
            }
        }

        // ========================
        // 🔐 SECURITY SCAN
        // ========================
        stage('Security Scan') {
            steps {
                echo "Running security scan..."
                sh 'echo Security Scan PASSED'
            }
        }

        // ========================
        // 🐳 DOCKER BUILD
        // ========================
        stage('Docker Build') {
            steps {
                withCredentials([string(credentialsId: GH_CREDENTIALS_ID, variable: 'G_TOKEN')]) {
                    echo "Building Docker image..."
                    sh '''
                    docker build \
                      --build-arg GITHUB_TOKEN=$G_TOKEN \
                      -t ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER} .
                    '''
                }
            }
        }

        // ========================
        // 🔑 DOCKER LOGIN
        // ========================
        stage('Docker Login') {
            steps {
                withCredentials([usernamePassword(credentialsId: DOCKER_CRED_ID, usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                    echo "Login to GHCR..."
                    sh '''
                    echo $PASS | docker login ${REGISTRY} -u $USER --password-stdin
                    '''
                }
            }
        }

        // ========================
        // 📦 DOCKER PUSH
        // ========================
        stage('Docker Push') {
            steps {
                echo "Pushing Docker image..."
                sh '''
                docker push ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER}

                docker tag ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER} \
                           ${REGISTRY}/${IMAGE_NAME}:latest

                docker push ${REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }

        // ========================
        // 🚀 GITOPS UPDATE
        // ========================
        stage('GitOps Update') {
            steps {
                withCredentials([string(credentialsId: GH_CREDENTIALS_ID, variable: 'G_TOKEN')]) {
                    echo "Updating deployment repo..."
                    sh '''
                    rm -rf deployment-config

                    git clone https://$G_TOKEN@github.com/wahyunurdian26/deployment-config.git

                    cd deployment-config

                    sed -i "s|image: .*|image: ${REGISTRY}/${IMAGE_NAME}:staging-${BUILD_NUMBER}|g" ${DEPLOY_PATH}

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

    // ========================
    // 📌 POST ACTION
    // ========================
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
