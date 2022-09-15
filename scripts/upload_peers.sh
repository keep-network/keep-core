#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'           # new line + bold + color
LOG_END='\n\e[0m'                # new line + reset color
DONE_START='\n\e[1;32m'          # new line + bold + green
DONE_END='\n\n\e[0m'             # new line + reset
LOG_WARNING_START='\n\e\033[33m' # new line + bold + warning color
LOG_WARNING_END='\n\e\033[0m'    # new line + reset

SOURCE_PROJECT_ID=keep-test-f3e0
DESTINATION_BUCKET_NAME="diagnostics_test"
PROJECT_NAME="keep-test"
PEERS_LOCATION="diagnostics/peers.json"

help() {
  echo -e "\nUsage: $0" \
    "--key-file-path <GCP-key-file-path>"
  echo -e "\n\nCommand line arguments:\n"
  echo -e "\t--key-file: GCP key file path\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--key-file") set -- "$@" "-k" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "k:h" opt; do
  case "$opt" in
  k) key_file_path="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

KEY_FILE_PATH=${key_file_path:-""}

if [ "$KEY_FILE_PATH" == "" ]; then
  printf "${LOG_WARNING_START}Key file must be provided.${LOG_WARNING_END}"
fi

# Run script
printf "${LOG_START}Authenticating with gcloud...${LOG_END}"
gcloud auth activate-service-account --key-file $key_file_path

printf "${LOG_START}Setting cloud context to the ${PROJECT_NAME}...${LOG_END}"
gcloud config set project $SOURCE_PROJECT_ID

printf "${LOG_START}Uploading peers info to a ${DESTINATION_BUCKET_NAME} bucket...${LOG_END}"
gcloud storage cp ${PEERS_LOCATION} gs://$DESTINATION_BUCKET_NAME/

printf "${DONE_START}Upload completed!${DONE_END}"
