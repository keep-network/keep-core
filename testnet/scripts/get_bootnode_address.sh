apk add --no-cache curl;
until nslookup ${BOOTNODE_SERVICE};
do
  echo "Waiting for bootnode to be ready..."; 
  sleep 2; 
done;
curl -m 5 -s ${BOOTNODE_SERVICE} | xargs echo -n >> /bootnode/bootnodes;