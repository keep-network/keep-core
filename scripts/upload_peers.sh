#!/bin/bash
set -eou pipefail

# This script should be run after scripts/diagnostics.sh which creates a single
# file with peers data.

LOG_START='\n\e[1;36m'           # new line + bold + color
LOG_END='\n\e[0m'                # new line + reset color
DONE_START='\n\e[1;32m'          # new line + bold + green
DONE_END='\n\n\e[0m'             # new line + reset
LOG_WARNING_START='\n\e\033[33m' # new line + bold + warning color
LOG_WARNING_END='\n\e\033[0m'    # new line + reset

BUCKET_NAME_DEFAULT="diagnostics_test"
PEERS_DIR_PATH="diagnostics"

help() {
  echo -e "\nUsage: $0" \
    "--oauth2-token-path <GCP-oauth2-token-path>" \
    "--bucket-name <GCP-bucket-name>"
  echo -e "\nCommand line arguments:\n"
  echo -e "\t--oauth2-token: GCP oauth2 token path"
  echo -e "\t--bucket-name: GCP destination bucket name\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--oauth2-token") set -- "$@" "-k" ;;
  "--bucket-name") set -- "$@" "-b" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "k:b:h" opt; do
  case "$opt" in
  k) oauth2_token_path="$OPTARG" ;;
  b) bucket_name="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

OAUTH2_TOKEN_PATH=${oauth2_token_path:-""}
BUCKET_NAME=${bucket_name:-${BUCKET_NAME_DEFAULT}}

if [ "$OAUTH2_TOKEN_PATH" == "" ]; then
  printf "${LOG_WARNING_START}OAuth2 token must be provided.${LOG_WARNING_END}"
  exit 1
fi

# Read the file name to be uploaded to GCP bucket
file_name=`find ${PEERS_DIR_PATH} -type f -exec basename {} \;`

curl -X POST --data-binary @${PEERS_DIR_PATH}"/"${file_name} \
     -H "Authorization: Bearer `cat ${OAUTH2_TOKEN_PATH}`" \
     -H "Content-Type: application/json" \
    "https://storage.googleapis.com/upload/storage/v1/b/${BUCKET_NAME}/o?name=${file_name}"

printf "${DONE_START}Upload completed!${DONE_END}"
