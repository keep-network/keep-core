pragma solidity 0.5.17;

contract TokenDistributor {
    // In the claim function, you need to provide Ethereum address and a signed
    // address of token recipient (the signature by Ethereum address from a
    // merkle tree). We'll validate the signature and see how many tokens should
    // be claimable by that address based on the information in Merkle tree.
    // TODO: Add support of Merkle Tree distribution
    function claim(
        address recipient,
        address destination,
        uint8 v,
        bytes32 r,
        bytes32 s // TODO: Add more parameters for claims from merkle tree.
    ) public {
        require(recipient != address(0), "recipient address cannot be zero");
        require(
            destination != address(0),
            "destination address cannot be zero"
        );

        require(
            recipient ==
                recoverSignerAddress(abi.encodePacked(destination), v, r, s),
            "invalid signature of destination address"
        );
    }

    function recoverSignerAddress(
        bytes memory message,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) internal pure returns (address) {
        // Validate `s` and `v` values for a malleability concern described in EIP-2.
        // Only signatures with `s` value in the lower half of the secp256k1
        // curve's order and `v` value of 27 or 28 are considered valid.
        require(
            uint256(s) <=
                0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A0,
            "Invalid signature 's' value"
        );
        require(v == 27 || v == 28, "Invalid signature 'v' value");

        bytes32 digest = keccak256(message);
        bytes32 prefixedDigest =
            keccak256(
                abi.encodePacked("\x19Ethereum Signed Message:\n32", digest)
            );

        return ecrecover(prefixedDigest, v, r, s);
    }

    // TODO: Tokens not claimed within a given timeout should go to a treasury
    // wallet address set on that contract.
}
