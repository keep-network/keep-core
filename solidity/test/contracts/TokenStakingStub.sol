pragma solidity 0.5.17;

import "../../contracts/TokenStaking.sol";


contract TokenStakingStub is TokenStaking {
    constructor(
        address _tokenAddress,
        address _registry,
        uint256 _initializationPeriod,
        uint256 _undelegationPeriod
    )
        public
        TokenStaking(
            _tokenAddress,
            _registry,
            _initializationPeriod,
            _undelegationPeriod
        )
    {}

    function setInitializationPeriod(uint256 _initializationPeriod) public {
        initializationPeriod = _initializationPeriod;
    }
}
