pragma solidity ^0.8.15;

import "./AltBn128.sol";

contract Secret {
    AltBn128.G1Point public sP;
    AltBn128.G1Point public xP;
    AltBn128.G2Point public xQ;

    bytes32 public currentRequestHash;
    // Request public currentRequest;
    // Response public currentResponse;

    event Request(
        AltBn128.G2Point yQ,
        AltBn128.G1Point xyP,
        bytes32 requestHash
    );
      

    event Response(
        AltBn128.G1Point s_xyP,
        AltBn128.G2Point s_xyQ
    );

    function registerGroup(
        bytes memory _sP,
        bytes memory _xP,
        bytes memory _xQ
    ) public {
        sP = AltBn128.g1Unmarshal(_sP);
        xP = AltBn128.g1Unmarshal(_xP);
        xQ = AltBn128.g2Unmarshal(_xQ);

        require(
            AltBn128.pairing(
                AltBn128.G1Point(xP.x, AltBn128.getP() - xP.y),
                AltBn128.g2(),
                AltBn128.g1(),
                xQ
            ),
            "Invalid group registration"
        );
    }

    function makeRequest(
        bytes memory yQ,
        bytes memory xyP
    ) public returns (bool) {
        AltBn128.G1Point memory _xyP = AltBn128.g1Unmarshal(xyP);
        AltBn128.G2Point memory _yQ = AltBn128.g2Unmarshal(yQ);
        currentRequestHash = keccak256(abi.encode(_yQ, _xyP));
        emit Request(_yQ, _xyP, currentRequestHash);
        return true;
    }

    function challengeRequest(
        bytes memory yQ,
        bytes memory xyP
    ) public returns (bool) {
        AltBn128.G1Point memory _xyP = AltBn128.g1Unmarshal(xyP);
        AltBn128.G2Point memory _yQ = AltBn128.g2Unmarshal(yQ);
        bytes32 requestHash = keccak256(abi.encode(_yQ, _xyP));
        require(requestHash == currentRequestHash, "Invalid challenge");
        return AltBn128.pairing(
            AltBn128.G1Point(_xyP.x, AltBn128.getP() - _xyP.y),
            AltBn128.g2(),
            xP,
            _yQ
        );
    }

    function respond(
        bytes memory yQ,
        bytes memory xyP,
        bytes memory s_xyQ
    ) public returns (bool) {
        AltBn128.G1Point memory _xyP = AltBn128.g1Unmarshal(xyP);
        AltBn128.G2Point memory _yQ = AltBn128.g2Unmarshal(yQ);
        bytes32 requestHash = keccak256(abi.encode(_yQ, _xyP));
        require(requestHash == currentRequestHash, "Invalid challenge");
        AltBn128.G2Point memory _s_xyQ = AltBn128.g2Unmarshal(s_xyQ);
        AltBn128.G1Point memory _s_xyP = AltBn128.g1Add(_xyP, sP);
        bool valid = AltBn128.pairing(
            AltBn128.G1Point(_s_xyP.x, AltBn128.getP() - _s_xyP.y),
            AltBn128.g2(),
            AltBn128.g1(),
            _s_xyQ
        );
        if (valid) {
            emit Response(_s_xyP, _s_xyQ);
        }
        return valid;
    }
        
}