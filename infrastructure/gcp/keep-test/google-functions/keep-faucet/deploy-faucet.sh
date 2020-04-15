gcloud functions deploy keep-faucet-ropsten \
  --trigger-http \
  --runtime nodejs10 \
  --entry-point issueGrant \
  --set-env-vars KEEP_CONTRACT_OWNER_ADDRESS=$KEEP_CONTRACT_OWNER_ADDRESS,KEEP_CONTRACT_OWNER_PRIVATE_KEY=$KEEP_CONTRACT_OWNER_PRIVATE_KEY,ETHEREUM_HOST=$ETHEREUM_HOST,ETHEREUM_NETWORK_ID=$ETHEREUM_NETWORK_ID
