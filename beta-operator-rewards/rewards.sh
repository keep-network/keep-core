#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'  # new line + bold + color
LOG_END='\n\e[0m'       # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

PROMETHEUS_API_DEFAULT="http://prometheus.monitoring.svc.cluster.local:9090/api/v1"
PROMETHEOUS_JOB_DEFAULT="keep-discovered-nodes"
REWARDS_END_DATE_DEFAULT=$(date +"%Y-%m-%d")
OUTPUT_JSON_FILE="peersData.json"
BUCKET_NAME_DEFAULT="diagnostics_test"
CLIENT_UPGRADE_DELAY_ACCEPTANCE=604800 # 7days in sec

help() {
  echo -e "\nUsage: $0" \
    "--prometheus-api <prometheus-api-address>" \
    "--prometheus-job <prometheus-job>" \
    "--rewards-start-date <rewards-start-date Y-m-d>" \
    "--rewards-end-date <rewards-end-date Y-m-d>" \
    "--bucket-name <GCP-bucket-name>"
  echo -e "\nCommand line arguments:\n"
  echo -e "\t--prometheus-api: Prometheus API"
  echo -e "\t--prometheus-job: Prometheus service discovery job name"
  echo -e "\t--rewards-start-date: Rewards start date Y-m-d"
  echo -e "\t--rewards-end-date: Rewards end date Y-m-d"
  echo -e "\t--bucket-name: GCP destination bucket name where peer data are stored\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--prometheus-api") set -- "$@" "-i" ;;
  "--prometheus-job") set -- "$@" "-p" ;;
  "--rewards-start-date") set -- "$@" "-k" ;;
  "--rewards-end-date") set -- "$@" "-e" ;;
  "--bucket-name") set -- "$@" "-b" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "a:d:p:k:e:b:h" opt; do
  case "$opt" in
  a) prometheus_api="$OPTARG" ;;
  p) prometheous_job="$OPTARG" ;;
  k) rewards_start_date="$OPTARG" ;;
  e) rewards_end_date="$OPTARG" ;;
  b) bucket_name="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

REWARDS_START_DATE=${rewards_start_date:-""}
REWARDS_END_DATE=${rewards_end_date:-${REWARDS_END_DATE_DEFAULT}}
PROMETHEUS_API=${prometheus_api:-${PROMETHEUS_API_DEFAULT}}
PROMETHEOUS_JOB=${prometheous_job:-${PROMETHEOUS_JOB_DEFAULT}}
BUCKET_NAME=${bucket_name:-${BUCKET_NAME_DEFAULT}}

if [ "$REWARDS_START_DATE" == "" ]; then
  printf "${LOG_WARNING_START}Rewards start date must be provided.${LOG_WARNING_END}"
  exit 1
fi

# TODO: convert to Y-m-d 23:59:59
rewardsStartDate=$(date -j -f "%Y-%m-%d" ${REWARDS_START_DATE} "+%s")
rewardsEndDate=$(date -j -f "%Y-%m-%d" ${REWARDS_END_DATE} "+%s")

printf "${LOG_START}Installing yarn dependencies...${LOG_END}"
yarn install

# latestTags=$(git tag --sort=-version:refname --list 'v[0-9]*.*-m[0-9]*') # TODO: REPLACE, this is correct regex to match latest versions
allTags=($(git tag --sort=-version:refname --list 'v[0-9]*')) # sorted latest -> oldest # TODO: For testing
latestTags=allTags[@]:0:1                                     # pick 2 latest tags

tagsInRewardInterval=()

latestTimestamp=$(git show -s --format=%ct ${allTags[0]}^{commit})
if [ $latestTimestamp -gt $rewardsStartDate ] && [ $latestTimestamp -lt $rewardsEndDate ]; then
  tagOneBeforeLatestTimestamp="${allTags[0]}_$(git show -s --format=%ct ${allTags[1]}^{commit})"
  tagLatestTimestamp="${allTags[1]}_$latestTimestamp"
  tagsInRewardInterval+=($tagOneBeforeLatestTimestamp)
  tagsInRewardInterval+=($tagLatestTimestamp)
elif [ $latestTimestamp -gt $rewardsEndDate ]; then
  # The latest tag was created after the given rewards interval. We do not process
  # such tags and we take the one before latest
  tagLatestTimestamp="${allTags[1]}_$latestTimestamp"
  tagsInRewardInterval+=($tagLatestTimestamp)
else
  # Latest tag was created before the start rewards interval and continue being
  # the latest tag. (no new releases)
  tagLatestTimestamp="${allTags[0]}_$latestTimestamp"
  tagsInRewardInterval+=($tagLatestTimestamp)
fi

# Converting array to string so we can pass to the rewards-requirements.ts
printf -v tags '%s|' "${tagsInRewardInterval[@]}"
tagsTrimmed="${tags%?}" # remove "|" at the end

# Run script
printf "${LOG_START}Fetching peers data...${LOG_END}"

yarn rewards-requirements \
  --api ${PROMETHEUS_API} \
  --job ${PROMETHEOUS_JOB} \
  --start $rewardsStartDate \
  --end $rewardsEndDate \
  --interval 5 \
  --versions $tagsTrimmed \
  --output ${OUTPUT_JSON_FILE}

# TODO: do we want to upload the output file to a GCP bucket?

printf "${DONE_START}Complete!${DONE_END}"
