#!/bin/bash

HELP="Usage: ./$(basename $0) -c <LOCAL_KUBE_CONTEXT> \nExample: ./$(basename $0) -c docker-for-desktop"

while getopts ":c:" opt; do
  case $opt in
    c ) LOCAL_KUBE_CONTEXT=$OPTARG;;

    \?)
      echo -e $HELP
      exit 1
  esac
done

if [ $# -eq 0 ]
then
  echo -e $HELP
  exit 1
fi

function use_local_context() {

  kubectl config use-context $LOCAL_KUBE_CONTEXT
}

function create_google_container_registry_secret() {

  DOCKER_PASSWORD="$(gcloud auth print-access-token)"
  DOCKER_EMAIL="$(gcloud info | grep Account | awk '{print $2}' | tr -d  "[]")"

  kubectl create secret docker-registry google-container-registry-auth \
    --docker-server=https://gcr.io \
    --docker-username=oauth2accesstoken \
    --docker-password=$DOCKER_PASSWORD \
    --docker-email=$DOCKER_EMAIL
}

echo "Setting kube context to local:"
use_local_context
echo "----------------"

echo "Creating secret for accessing Google private container registry:"
create_google_container_registry_secret