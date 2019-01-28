#!/bin/bash
# This script extracts the keyfiles for pregenerated Ethereum
# accound from a Docker image and uploads them into a Kubernetes
# ConfigMap.

display_usage() {
	echo "This script requires arguments."
  echo -e "\nUsage:\nmanage-keyfiles-as-ConfigMaps [arguments]\n"
  echo -e "-e\n"
  echo -e "\tExtract keyfiles from Docker image into keystore directory.\n"
  echo -e "-c\n"
  echo -e "\tCreate ConfigMap from files in keystore directory using account IDs as key.\n"
  echo -e "-d\n"
  echo -e "\tDelete ConfigMap from files found in keystore directory using account IDs as key.\n"
  echo -e "-l\n"
  echo -e "\tList ConfigMap from files found in keystore directory using account IDs as key.\n"
  echo -e "-h\n"
  echo -e "\tDisplay this help message."
}

# if less than two arguments supplied, display usage
if [  $# -le 0 ]
then
  display_usage
  exit 1
fi

# check whether user had supplied -h or --help . If yes display usage
if [[ ( $# == "--help") ||  $# == "-h" ]]
then
display_usage
exit 0
fi

# parse commandline arguments
POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -e)
    EXTRACT=true
    shift # past argument
    ;;
    -c)
    CREATE=true
    shift # past argument
    ;;
    -d)
    DELETE=true
    shift # past argument
    ;;
    -l)
    LIST=true
    shift # past argument
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

# set keystore directory
KEYSTORE=`pwd`/keystore

# extract keyfiles from docker image
if [ "$EXTRACT" = true ]; then
  echo -e "\nExtracting files from container $DOCKER_ID_USER/geth-node into $KEYSTORE\n"
  docker run --entrypoint="" --rm -v `pwd`:/out $DOCKER_ID_USER/geth-node \
    cp -rv /root/.geth/keystore/ /out
fi

# upload to ConfigMap
if [ "$CREATE" = true ]; then
  echo -e "\nCreating ConfigMap from files found in $KEYSTORE\n"
  ls -1 $KEYSTORE | while read FILENAME; do
    ACCOUNT=`echo "$FILENAME" | cut -d "-" -f9`
    echo -e "$ACCOUNT\t$FILENAME"
    kubectl create configmap $ACCOUNT --from-file=$ACCOUNT=$KEYSTORE/$FILENAME
  done;
fi

# list ConfigMap
if [ "$LIST" = true ]; then
  echo -e "\nDescribing ConfigMap from keys found in $KEYSTORE\n"
  ls -1 $KEYSTORE | while read FILENAME; do
    ACCOUNT=`echo "$FILENAME" | cut -d "-" -f9`
    kubectl describe configmaps $ACCOUNT
  done;
fi

# delete ConfigMapls -1 keystore | while read FILENAME; do
if [ "$DELETE" = true ]; then
    echo -e "\nDeleting keys found in $KEYSTORE\n"
    ls -1 $KEYSTORE | while read FILENAME; do
      ACCOUNT=`echo "$FILENAME" | cut -d "-" -f9`
      kubectl delete configmaps $ACCOUNT
    done;
    rm -rf $KEYSTORE
fi
