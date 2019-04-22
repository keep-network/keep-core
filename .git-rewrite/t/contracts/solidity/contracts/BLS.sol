pragma solidity ^0.5.4;

import "./utils/AltBn128.sol";


/**
 * @title BLS signatures verification
 * @dev Library for verification of 2-pairing-check BLS signatures, including
 * basic, aggregated, or reconstructed threshold BLS signatures, generated
 * using the AltBn128 curve.
 */
library BLS {

    /**
     * @dev Verify performs the pairing operation to check if the signature
     * is correct for the provided message and the corresponding public key.
     */
    function verify(bytes memory publicKey, bytes memory message, bytes32 signature) public view returns (bool) {

        AltBn128.G1Point memory _signature;
        _signature = AltBn128.g1Decompress(signature);

        return AltBn128.pairing(
            AltBn128.G1Point(_signature.x, AltBn128.getP() - _signature.y),
            AltBn128.g2(),
            AltBn128.g1HashToPoint(message),
            AltBn128.g2Decompress(publicKey)
        );
    }
}
