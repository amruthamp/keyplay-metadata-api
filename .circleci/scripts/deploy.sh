#!/bin/bash

# Exit on any error
set -e

export CUR_DIR=`dirname $0`
source ${CUR_DIR}/common.sh
export GITHASH=`git rev-parse --short HEAD`
export DEPLOY_TIMESTAMP=`date -u`

# Auth
gcloud_auth

service_name=keyplay-metadata-api
gcloud docker -- push gcr.io/"${GCLOUD_PROJECT_ID}"/${service_name}:"${CIRCLE_BUILD_NUM}"
export NEW_IMAGE_NAME="gcr.io/${GCLOUD_PROJECT_ID}/${service_name}:${CIRCLE_BUILD_NUM}"

# Deploy
deploy_pods ${service_name} "${GCLOUD_PROJECT_ID}" us-west2-b api-cluster-west2
deploy_pods ${service_name} "${GCLOUD_PROJECT_ID}" us-east4-b api-cluster-east4
