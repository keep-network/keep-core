FROM ethereumex/node

RUN git clone https://github.com/lispmeister/eth-netstats.git /var/lib/eth-netstats
WORKDIR /var/lib/eth-netstats
RUN npm install
RUN grunt all

ENTRYPOINT ["/bin/bash"]
