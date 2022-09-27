#!/bin/bash

# Exit on any error
set -e

export CUR_DIR=`dirname $0`
source ${CUR_DIR}/common.sh
export GITHASH=`git rev-parse --short HEAD`
export DEPLOY_TIMESTAMP=`date -u`

# Configure gcloud for the cluster (Source Project Environment)
gcloud config set project "${SOURCE_PROJECT_ID}"
gcloud config set compute/zone us-east4-b
gcloud config set container/cluster api-cluster-east4

# Auth
gcloud_auth

# Set the credentials to be used by kubectl
gcloud container clusters get-credentials api-cluster-east4

# Get the current image running in the Source Project Environment
export CURRENT_SOURCE_DEPLOYMENT=`kubectl get deployments | awk '$1 ~ /^keyplay-metadata-api/ {print $1}'`
export NEW_IMAGE_NAME=`kubectl get deployment ${CURRENT_SOURCE_DEPLOYMENT} -o jsonpath="{.spec.template.spec.containers[0].image}"`

# Deploy
service_name=keyplay-metadata-api
deploy_pods ${service_name} "${GCLOUD_PROJECT_ID}" us-west2-b api-cluster-west2
deploy_pods ${service_name} "${GCLOUD_PROJECT_ID}" us-east4-b api-cluster-east4
if [[ ${DEPLOY_ENV} == "production" ]]; then
    deploy_pods ${service_name} "${GCLOUD_PROJECT_ID}" us-central1-f api-cluster-central1
fi