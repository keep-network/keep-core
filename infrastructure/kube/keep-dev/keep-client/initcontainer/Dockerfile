FROM ethereum/client-go AS runtime

COPY ./unlock-and-stake-eth-account.sh /tmp/unlock-and-stake-eth-account.sh

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["./tmp/unlock-and-stake-eth-account.sh"]
