pragma solidity ^0.4.17;

// Proof of Existence contract, version 1
contract ProofOfExistence1 {
  // state
  bytes32 public proof;
  // calculate and store the proof for a document
  // *transactional function*
  function notarize(string document) public {
    proof = proofFor(document);
  }
  // helper function to get a document's sha256
  // *read-only function*
  function proofFor(string document) pure public returns (bytes32) {
    return sha256(document);
  }

  function double(int a) pure public returns (int) {
    return 10*a;
  } 
}

