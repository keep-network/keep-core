#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'           # new line + bold + color
LOG_END='\n\e[0m'                # new line + reset color
LOG_WARNING_START='\n\e\033[33m' # new line + bold + warning color
LOG_WARNING_END='\n\e\033[0m'    # new line + reset
DONE_START='\n\e[1;32m'          # new line + bold + green
DONE_END='\n\n\e[0m'             # new line + reset

IP_ADDRESS_DEFAULT="127.0.0.1"
PORT_DEFAULT="8081"
PEER_PORTS_DEFAULT=("8081" "8082" "8083")
ENDPOINT="diagnostics"

help() {
  echo -e "\nUsage: $0" \
    "--ip-address <bootstrap ip address>" \
    "--port <bootstrap port>" \
    "--peer-ports <range of the eligible peer ports>" \
  echo -e "\n\nOptional command line arguments:\n"
  echo -e "\t--ip-address: IP address of the bootstrap node\n"
  echo -e "\t--port: Port of the bootstrap node\n"
  echo -e "\t--peer-ports: Range of the eligible peer ports\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--ip-address") set -- "$@" "-i" ;;
  "--port") set -- "$@" "-p" ;;
  "--peer-ports") set -- "$@" "-s" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "i:p:s:h" opt; do
  case "$opt" in
  i) ip_address="$OPTARG" ;;
  p) port="$OPTARG" ;;
  s) peer_ports="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
IP_ADDRESS=${ip_address:-$IP_ADDRESS_DEFAULT}
PORT=${port:-$PORT_DEFAULT}
PEER_PORTS=${peer_ports:-$PEER_PORTS_DEFAULT}

# Run script
printf "${LOG_START}Starting bootstrap diagnostics...${LOG_END}"

ips=$(curl ${IP_ADDRESS}:${PORT}/${ENDPOINT} | jq -r '.connected_peers[].ip')

peers=()
# Iterate over the connected peers using their IP addresses
for ip in ${ips[@]}; do
  # Iterate over the eligible peer ports. First port that returns no error breaks
  # the loop and assign data to peer[] array.
  for port in ${PEER_PORTS[@]}; do
    if peer=$(curl ${ip}:${port}/${ENDPOINT}); then
      peerExtractedData=$(jq -r '{version: .client_info.version, preParamPoolSize: .tbtc.preParamsPoolSize}' <<<"${peer}")
      peers+=(${peerExtractedData})
      break
    else
      echo "${ip}:${port}/${ENDPOINT} is not a valid endpoint"
    fi
  done
done

peersJsonArray=$(jq -s <<<"${peers[@]}")
peersJson=$(jq --null-input -r --argjson peersJsonArray "${peersJsonArray}" '{peers: $peersJsonArray}')

echo $peersJson

printf "${DONE_START}Bootstrap diagnostics completed!${DONE_END}"
