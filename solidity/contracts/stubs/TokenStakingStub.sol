pragma solidity 0.5.17;

import "../TokenStaking.sol";
import "../TokenStakingEscrow.sol";
import "../TokenGrant.sol";
import "../KeepRegistry.sol";

contract TokenStakingStub is TokenStaking {
    constructor(
        ERC20Burnable _token,
        TokenGrant _tokenGrant,
        TokenStakingEscrow _escrow,
        KeepRegistry _registry,
        uint256 _initializationPeriod,
        uint256 _undelegationPeriod
    ) TokenStaking(_token, _tokenGrant, _escrow, _registry, _initializationPeriod, _undelegationPeriod) public {
    }

    function setInitializationPeriod(uint256 _initializationPeriod) public {
        initializationPeriod = _initializationPeriod;
    }
}
