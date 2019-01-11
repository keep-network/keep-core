FROM ubuntu
MAINTAINER "Markus Fix <Markus.Fix@keep.network>

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y build-essential
RUN apt-get install -y nodejs npm git curl

RUN npm install -g grunt
RUN npm install -g pm2

RUN git clone https://github.com/lispmeister/eth-netstats.git /var/lib/eth-netstats
WORKDIR /var/lib/eth-netstats
RUN npm install
RUN grunt all

RUN git clone https://github.com/lispmeister/bootnode-registrar.git /var/lib/bootnode
WORKDIR /var/lib/bootnode
RUN npm install

RUN useradd -ms /bin/bash dashboard
USER dashboard

WORKDIR /home/dashboard
COPY app.json /home/dashboard/app.json
COPY run.sh /home/dashboard/run.sh

COPY updateNode.sh /home/dashboard/updateNode.sh
RUN /bin/bash /home/dashboard/updateNode.sh

ENTRYPOINT ["/bin/bash", "run.sh"]
