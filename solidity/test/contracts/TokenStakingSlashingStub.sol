pragma solidity 0.5.17;

import "../../contracts/TokenStaking.sol";


contract TokenStakingSlashingStub is TokenStaking {
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

    function slash(uint256 amountToSlash, address[] memory misbehavedOperators)
        public
    {
        for (uint256 i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            emit TokensSlashed(operator, 1 ether);
        }
    }

    function seize(
        uint256 amountToSeize,
        uint256 rewardMultiplier,
        address tattletale,
        address[] memory misbehavedOperators
    ) public {
        for (uint256 i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            emit TokensSeized(operator, 1 ether);
        }
    }
}
