#!/bin/bash

HELP="Usage: ./$(basename $0) -e <ENVIRONMENT> -r <REGION>"

while getopts ":e:r:" opt; do
  case $opt in
    e ) ENVIRONMENT=$OPTARG;;
    r ) REGION=$OPTARG;;

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

function download_gke_creds() {

  PROJECT_ID=`gcloud projects list | grep -i $ENVIRONMENT | awk '{print $1}'`
  CLUSTER_NAME=`gcloud container clusters list --project $PROJECT_ID | grep -i $ENVIRONMENT | awk '{print $1}'`

  gcloud container clusters get-credentials $CLUSTER_NAME --region $REGION --project $PROJECT_ID --internal-ip
}

download_gke_creds