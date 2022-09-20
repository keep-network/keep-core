#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m'  # new line + bold + color
LOG_END='\n\e[0m'       # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

IP_ADDRESS_DEFAULT="127.0.0.1" # bootstrap node
PORT_DEFAULT="9701"            # bootstrap port
PEER_PORTS_DEFAULT=("9701")    # a range of ports can be provided
ENDPOINT="diagnostics"         # diagnostics api endpoint
DIAGNOSTICS_DIR="diagnostics"  # dir for saving peers info
# At least one of the ports in the range of PEER_PORT_START_DEFAULT..PEER_PORT_END_DEFAULT
# should be opened for diagnostics.
PEER_PORT_START_DEFAULT=9701 # peer port start lookup
PEER_PORT_END_DEFAULT=9701   # peer port end lookup

help() {
  echo -e "\nUsage: $0" \
    "--address <bootstrap address>" \
    "--port <bootstrap port>" \
    "--peer-port-start <beginning-ports-lookup>" \
    "--peer-port-end <end-ports-lookup>"
  echo -e "\n\nOptional command line arguments:\n"
  echo -e "\t--address: Address of the bootstrap node"
  echo -e "\t--port: Port of the bootstrap node"
  echo -e "\t--peer-ports-start: Beginning of the ports lookup"
  echo -e "\t--end-ports-start: End of the ports lookup\n"
  exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--address") set -- "$@" "-a" ;;
  "--port") set -- "$@" "-p" ;;
  "--peer-port-start") set -- "$@" "-s" ;;
  "--peer-port-end") set -- "$@" "-e" ;;
  "--help") set -- "$@" "-h" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

# Parse short options
OPTIND=1
while getopts "a:p:s:e:h" opt; do
  case "$opt" in
  a) address="$OPTARG" ;;
  p) port="$OPTARG" ;;
  s) peer_port_start="$OPTARG" ;;
  e) peer_port_end="$OPTARG" ;;
  h) help ;;
  ?) help ;; # Print help in case parameter is non-existent
  esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
ADDRESS=${address:-$IP_ADDRESS_DEFAULT}
PORT=${port:-$PORT_DEFAULT}
PEER_PORTS=${peer_ports:-$PEER_PORTS_DEFAULT}
PEER_PORT_START=${peer_port_start:-$PEER_PORT_START_DEFAULT}
PEER_PORT_END=${peer_port_end:-$PEER_PORT_END_DEFAULT}

# Run script
printf "${LOG_START}Starting bootstrap diagnostics...${LOG_END}"

# Clean up for fresh data
rm -rf ${DIAGNOSTICS_DIR}
mkdir ${DIAGNOSTICS_DIR}

peerAddresses=$(curl ${ADDRESS}:${PORT}/${ENDPOINT} | jq -r '.connected_peers[].address')

peers=()
# Iterate over the connected peers using their IP/dns addresses
for peerAddress in ${peerAddresses[@]}; do
  # Iterate over the eligible peer ports. First port that returns no error breaks
  # the loop and assign data to peers[] array.
  for port in $(seq $PEER_PORT_START $PEER_PORT_END); do
    if peer=$(curl ${peerAddress}:${port}/${ENDPOINT}); then
      peerExtractedData=$(jq -r '{chain_address: .client_info.chain_address, version: .client_info.version, preParamPoolSize: .tbtc.preParamsPoolSize}' <<<"${peer}")
      peers+=(${peerExtractedData})
      break
    else
      echo "${peerAddress}:${port}/${ENDPOINT} is not a valid endpoint"
    fi
  done
done

peersJsonArray=$(jq -s <<<"${peers[@]}")
peersJson=$(jq --null-input -r --argjson peersJsonArray "${peersJsonArray}" '{peers: $peersJsonArray}')

timestamp=$(date +%s) # unix timestamp, seconds since Jan 01 1970
printf "${LOG_START}Saving diagnostics to a file ${DIAGNOSTICS_DIR}/peers_$timestamp.json ..${LOG_END}"
echo $peersJson >"${DIAGNOSTICS_DIR}/peers_$timestamp.json"

printf "${DONE_START}Bootstrap diagnostics completed!${DONE_END}"
