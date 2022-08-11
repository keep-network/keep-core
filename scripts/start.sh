#!/bin/bash
set -eou pipefail

# Dafault inputs.
LOG_LEVEL_DEFAULT="info"
KEEP_CORE_PATH=$PWD
CONFIG_DIR_PATH_DEFAULT="$KEEP_CORE_PATH/configs"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--config-dir")  set -- "$@" "-c" ;;
    *)               set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "c:" opt
do
   case "$opt" in
      c ) config_dir_path="$OPTARG" ;;
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}

config_files=($CONFIG_DIR_PATH/*.toml)
config_files_count=${#config_files[@]}
while :
do
    printf "\nSelect client config file: \n"
    i=1
    for o in "${config_files[@]}"; do
        echo "$i) ${o##*/}"
        let i++
    done

    read reply
    if [ "$reply" -ge 1 ] && [ "$reply" -le $config_files_count ]; then
        CONFIG_FILE_PATH=${config_files["$reply"-1]}
        break
    else
        printf "\nInvalid choice. Please choose an existing option number.\n"
    fi
done
printf "\nClient config file: \"$CONFIG_FILE_PATH\" \n\n"

log_level_options=("info" "debug" "custom...")
while :
do
    echo "Select log level [$LOG_LEVEL_DEFAULT]: "
    i=1
    for o in  "${log_level_options[@]}"; do
        echo "$i) $o"
        let i++
    done

    read reply
    case $reply in
        "1"|"${log_level_options[0]}") LOG_LEVEL=${log_level_options[0]}; break;;
        "2"|"${log_level_options[1]}") LOG_LEVEL=${log_level_options[1]}; break;;
        "3"|"${log_level_options[2]}") 
            read -p "Enter custom log level: [$LOG_LEVEL_DEFAULT]" log_level
            LOG_LEVEL=${log_level:-$LOG_LEVEL_DEFAULT}
            break
            ;;
        "") LOG_LEVEL=$LOG_LEVEL_DEFAULT; break;;
        *) echo "Invalid choice. Please choose an existing option number.";;
    esac
done
echo "Log level: \"$LOG_LEVEL\""

# Run script.
LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

KEEP_CORE_CONFIG_FILE_PATH=$CONFIG_FILE_PATH

printf "${LOG_START}Starting keep-core client...${LOG_END}"
cd $KEEP_CORE_PATH
KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
    LOG_LEVEL=${LOG_LEVEL} \
    ./keep-client --config $KEEP_CORE_CONFIG_FILE_PATH start
