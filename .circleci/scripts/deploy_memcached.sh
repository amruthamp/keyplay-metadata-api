#!/bin/bash

# Exit on any error
set -e

export GITHASH=`git rev-parse --short HEAD`
export DEPLOY_TIMESTAMP=`date -u`
export CUR_DIR=`dirname $0`

# Auth
gcloud_auth

gcloud config set project ${GCLOUD_PROJECT_ID}
envsubst < "${CUR_DIR}"/../kubernetes/memcached-server-deployment.yml > dist-memcached.yml

deploy_memcached() {
    zone=$1
    cluster=$2
    service_name=$3
    gcloud config set compute/zone "${zone}"
    gcloud config set container/cluster ${cluster}

    # Set the credentials to be used by kubectl
    gcloud container clusters get-credentials "${cluster}"

    #initialize the helm repo after setting all environment variables
    helm init --upgrade

    export MEMCACHED_SERVICE=$(kubectl get service | awk -v pat="${service_name}-memcached" '$1 ~ pat  {print $1}')

    if [[ -n ${MEMCACHED_SERVICE} ]]; then
        echo "stopping memcached cluster ..."
        helm delete "${service_name}" --purge
    fi

    echo "starting memcached cluster ..."
    helm install --name "${service_name}" -f dist-memcached.yml stable/memcached
    sleep 10

    local pod_count=0

    while [[ ${pod_count} != "$MEMCACHED_REPLICAS" ]]; do
        sleep 3
        v=$(kubectl get statefulset "${service_name}-memcached" | awk 'FNR ==2 {print $2} ')
        IFS='/ ' read -r -a array <<< "${v}"
        pod_count="${array[0]}"
        echo "${pod_count} / ${MEMCACHED_REPLICAS} servers started"
    done

    echo "memcached cluster started"
}

deploy_memcached us-west2-b api-cluster-west2 "${MEMCACHED_SERVICE_NAME}"
deploy_memcached us-east4-b api-cluster-east4 "${MEMCACHED_SERVICE_NAME}"
if [[ ${DEPLOY_ENV} == "production" ]]; then
    deploy_memcached us-central1-f api-cluster-central1 "${MEMCACHED_SERVICE_NAME}"
fi