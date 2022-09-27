#!/bin/bash

gcloud_auth() {
  if [[ ! -f /root/.ssh/google_compute_engine ]]; then
    echo "Authenticating with gcloud..."
    openssl aes-256-cbc -d -in .circleci/service-keys/${PROJECT_KEY} -k "${GCLOUD_KEY}" -md md5 >gcloud.json
    gcloud auth activate-service-account ${SERVICE_ACCOUNT} --key-file gcloud.json
    ssh-keygen -f ~/.ssh/google_compute_engine -N ""
    export GOOGLE_APPLICATION_CREDENTIALS=gcloud.json
  else
    echo "Already authenticated, skipping gcloud authentication."
  fi
}

deploy_pods() {
  service_name=$1
  project=$2
  zone=$3
  cluster=$4

  envsubst <.circleci/kubernetes/deployment.yml >services-deployment.yml

  gcloud config set project "${project}"
  gcloud config set compute/zone "${zone}"
  gcloud config set container/cluster "${cluster}"
  gcloud container clusters get-credentials "${cluster}"

  export CURRENT_DEPLOYMENT=$(kubectl get deployments | awk '$1 ~ /^keyplay-metadata-api/ {print $1}')

  if [[ -z ${CURRENT_DEPLOYMENT} ]]; then
    # Deploy
    kubectl create -f services-deployment.yml --record
    kubectl create -f .circleci/kubernetes/service.yml || true
  else
    # Apply the updated deployment config
    kubectl apply -f services-deployment.yml
  fi

  kubectl rollout status deploy/"${service_name}"

  # apply hpa configuration
  envsubst < .circleci/kubernetes/hpa.yml > hpa.yml
  kubectl apply -f hpa.yml
}
