pragma solidity ^0.4.24;

import "./utils/BytesUtils.sol";


contract Verifier {

    using BytesUtils for bytes;

    // Obsolete when https://github.com/ethereum/solidity/issues/864 is implemented.
    function recoverAddress(bytes32 messageHash, bytes signature) public pure returns (address) {

        bytes32 r = signature.readBytes32(0);
        bytes32 s = signature.readBytes32(32);

        // Get ECDSA Recovery (parity) bit. Public key recovery equation
        // results in two `y`, hence the need of recovery bit, the specs:
        // recovery bit + 27 for uncompressed recovered pubkeys.
        // recovery bit + 31 for compressed recovered pubkeys.
        uint8 v = uint8(signature[64]) + 27;

        if (v != 27 && v != 28) {
            return 0;
        } else {
            bytes memory prefix = "\x19Ethereum Signed Message:\n32";
            return ecrecover(keccak256(prefix, messageHash), v, r, s);
        }
    }

    function isSigned(bytes32 messageHash, bytes signature, address signer) public pure returns (bool) {
        return recoverAddress(messageHash, signature) == signer;
    }
}
