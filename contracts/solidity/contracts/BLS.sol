pragma solidity ^0.4.24;

import "./AltBn128.sol";


/**
 * @title BLS signatures verification
 * @dev Library for verification of aggregated or reconstructed threshold BLS signatures
 * generated using AltBn128 curve.
 */
library BLS {

    /**
     * @dev Verify performs the pairing operation to check if the signature
     * is correct for the provided message and the corresponding public key.
     */
    function verify(bytes publicKey, bytes message, bytes32 signature) public view returns (bool) {

        uint256[2] memory _signature;
        (_signature[0], _signature[1]) = AltBn128.g1Decompress(signature);

        uint256[2] memory _message;
        (_message[0], _message[1]) = AltBn128.g1HashToPoint(message);

        uint256[2] memory _publicKeyX;
        uint256[2] memory _publicKeyY;
        (_publicKeyX, _publicKeyY) = AltBn128.g2Decompress(publicKey);

        return AltBn128.pairing(
            [_signature[0], AltBn128.getP() - _signature[1]],
            AltBn128.g2(),
            _message,
            [_publicKeyX[0], _publicKeyX[1], _publicKeyY[0], _publicKeyY[1]]
        );
    }
}
