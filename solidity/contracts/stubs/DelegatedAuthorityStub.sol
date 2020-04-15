pragma solidity 0.5.17;

import "../TokenStaking.sol";

contract DelegatedAuthorityStub {
    address recognizedContract;

    constructor(address _recognizedContract) public {
        recognizedContract = _recognizedContract;
    }

    function __isRecognized(address _contract) public view returns (bool) {
        return _contract == recognizedContract;
    }

    function claimAuthorityRecursively(
        address stakingContract,
        address source
    ) public {
        TokenStaking(stakingContract).claimDelegatedAuthority(source);
    }
}
