#!/bin/bash
if [[ -z $GOOGLE_PROJECT_NAME || -z $GOOGLE_PROJECT_ID || -z $BUILD_TAG || -z $GOOGLE_REGION || -z $GOOGLE_COMPUTE_ZONE_A || -z $TRUFFLE_NETWORK ]]; then
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

# Copy deployment artifacts over
ssh utilitybox rm -rf /tmp/$BUILD_TAG
ssh utilitybox mkdir /tmp/$BUILD_TAG
scp -r contracts/solidity utilitybox:/tmp/$BUILD_TAG/

# Run deployment
ssh utilitybox << EOF
  gcloud container clusters get-credentials $GOOGLE_PROJECT_NAME --region $GOOGLE_REGION --internal-ip --project=$GOOGLE_PROJECT_ID

  nohup timeout 600 kubectl port-forward svc/eth-tx-node 8545:8545 2>&1 > /dev/null &

  sleep 10s

  geth --exec "personal.unlockAccount(\"${CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS}\", \"${CONTRACT_OWNER_ETH_ACCOUNT_PASSWORD}\", 600)" attach http://localhost:8545

  cd /tmp/$BUILD_TAG/solidity

  npm init -y
  npm install truffle@5.0.7
  npm install openzeppelin-solidity
  npm install solidity-bytes-utils
  npm install babel-register
  npm install babel-polyfill

  cp ./truffle_sample.js ./truffle.js
  ./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
EOF

scp utilitybox:/tmp/$BUILD_TAG/contracts/solidity/build/contracts/* /tmp/contracts
ssh utilitybox rm -rf /tmp/$BUILD_TAG