#!/bin/bash
set -ex

if [[ -z $GOOGLE_PROJECT_NAME || -z $GOOGLE_PROJECT_ID || -z $BUILD_TAG || -z $GOOGLE_REGION || -z $GOOGLE_COMPUTE_ZONE_A ]]; then
  echo "one or more required variables are undefined"
  exit 1
fi

UTILITYBOX_IP=$(gcloud compute instances --project $GOOGLE_PROJECT_ID describe $GOOGLE_PROJECT_NAME-utility-box --zone $GOOGLE_COMPUTE_ZONE_A --format json | jq .networkInterfaces[0].networkIP -r)

# Setup ssh environment
gcloud compute config-ssh --project $GOOGLE_PROJECT_ID -q
cat >> ~/.ssh/config << EOF
Host *
  StrictHostKeyChecking no

Host utilitybox
  HostName $UTILITYBOX_IP
  IdentityFile ~/.ssh/google_compute_engine
  ProxyCommand ssh -W %h:%p $GOOGLE_PROJECT_NAME-jumphost.$GOOGLE_COMPUTE_ZONE_A.$GOOGLE_PROJECT_ID
EOF

ssh utilitybox rm -rf /tmp/$BUILD_TAG
