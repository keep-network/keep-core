BASEDIR=$(dirname "$0")

EXTERNAL_PARTICIPANT_PATH="${BASEDIR}/../keep-test/ropsten/participants"
EXTERNAL_PARTICIPANTS=$(ls $EXTERNAL_PARTICIPANT_PATH)

function clean_deployment_bundles() {
  echo "=====REMOVING OLD BUNDLES====="
  for participant in $EXTERNAL_PARTICIPANTS
  do
    echo $EXTERNAL_PARTICIPANT_PATH
    PARTICIPANT_PATH="$EXTERNAL_PARTICIPANT_PATH/$participant"
    CURRENT_BUNDLE=$(ls $PARTICIPANT_PATH | grep ".tar.gz")
    echo $PARTICIPANT_PATH

    if [ -z $CURRENT_BUNDLE ]
    then
      echo "No bundle in $PARTICIPANT_PATH to remove..."
    else
      echo "Removing $CURRENT_BUNDLE from $PARTICIPANT_PATH..."
      rm $PARTICIPANT_PATH/$CURRENT_BUNDLE
      echo "Bundle removed!"
    fi
  done
  echo "==============================\n"
}

function create_deployment_bundles() {
  echo "=====CREATING DEPLOYMENT BUNDLES====="
  DATE=$(date +%F)

  for participant in $EXTERNAL_PARTICIPANTS
  do
    PARTICIPANT_PATH=$EXTERNAL_PARTICIPANT_PATH/$participant

    echo "$participant"
    tar -zcvf $PARTICIPANT_PATH/$DATE-keep-client-deployment-bundle.tar.gz \
      -C $BASEDIR/../../../docs/ keep-client-quickstart.adoc \
      -C ../infrastructure/eth-networks/keep-test/ropsten changelog.adoc \
         ./keep-client-snapshot.tar \
      -C ./participants/$participant config \
         ./persistence \
         ./.secrets/eth-account-password.txt
     echo "==============================\n"
   done
}

function fetch_keep_client_image() {
  echo "=====FETCHING LATEST KEEP-CLIENT IMAGE====="
  docker pull gcr.io/keep-test-f3e0/keep-client
  echo "==============================\n"
}

function save_keep_client_image() {
  echo "=====SAVING LATEST KEEP-CLIENT IMAGE====="
  echo "This will take several seconds..."
  docker save -o ../keep-test/ropsten/keep-client-snapshot.tar gcr.io/keep-test-f3e0/keep-client
  echo "==============================\n"
}

clean_deployment_bundles
fetch_keep_client_image
save_keep_client_image
create_deployment_bundles
echo "All done!"

