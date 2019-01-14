FROM lispmeister/geth-node
MAINTAINER "Markus Fix <Markus.Fix@keep.network>

# Change to /root before provisioning our services
WORKDIR /root

# prebake the DAG to speedup geht-node boot time
COPY geth-dag.sh geth-dag.sh
RUN /root/geth-dag.sh

# when starting geth we assume the chain has already been initialized
# and we are running with a prebaked DAG
COPY run-geth.sh run-geth.sh

ENTRYPOINT ["pm2", "start", "--no-daemon", "app.json"]
