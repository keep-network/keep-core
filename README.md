# Merkle Distribution
[![Build Status](https://github.com/1inch/merkle-distribution/actions/workflows/test.yml/badge.svg)](https://github.com/1inch/merkle-distribution/actions)
[![Coverage Status](https://coveralls.io/repos/github/1inch/merkle-distribution/badge.svg?branch=master)](https://coveralls.io/github/1inch/merkle-distribution?branch=master)

Set of smart contracts for gas efficient merkle tree drops. 

## Sequential cumulative Merkle Tree drops

Each next Merkle Tree root replaces previous one and should contain cumulative balances of all the participants. Cumulative claimed amount is used as invalidation for every participant.

## Signature-based drop

Each entry of the drop contains private key which is used to sign the address of the receiver. This is done to safely distribute the drop and prevent MEV stealing.
