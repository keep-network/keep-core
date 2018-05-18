#!/usr/bin/env bash

read -d '' INFO << TEXT
Purpose:     To return the requested version number increment for this project
             based on the contents of the .version file and to store the new
             value back to the .version file.

Assumptions: 1) You run this command in the root directory of the versioned project
             2) The version binary is in the path

Usage:       version [ + | ++ | +++ | --reset=<NEW_VERSION_NUMBER> | --help | --usage]
Notes:       1) This script uses https://github.com/stretchr/version
             2) The actual version number is stored in a file named .version
             3) The version binary (***) has been placed in the /usr/bin directory or equivalent
             (***) Download and extract version binary from: https://github.com/stretchr/version

Examples:
  version.sh                # To return the current version number
  version.sh +              # To increment the build number
  version.sh ++             # To increment the minor number
  version.sh +++            # To increment the major number
  version.sh --reset=0.0.0  # To reset the version number to "0.0.0"
TEXT

USAGE=`echo "$INFO"|grep Usage`

if [ "$1" == "" ]; then    
   echo "`version -n -v=false ./`"
elif [ "$1" == "+" ] || [ "$1" == "++" ] || [ "$1" == "+++" ]; then    
   version="`version -n -v=false ./ $1`"
   echo $version

else    
   RESET_ARG=`echo $1| cut -d'=' -f 1`
   if [ "$RESET_ARG" == "--reset" ]; then
      RESET_VERSION_NUMBER=`echo $1| cut -d'=' -f 2`
      if [ "$RESET_VERSION_NUMBER" == "0.0.0" ]; then
         rm .version 2>/dev/null
      else
         printf "v$RESET_VERSION_NUMBER" > .version
      fi
      version="`version -n -v=false ./`"
      echo $version
   elif [ "$1" == "--usage" ]; then
      echo "$USAGE" 
      exit 1
   else
      echo "$INFO"
      exit 1
   fi
fi
