pragma solidity ^0.4.24;

import "../mixins/StakingManager.sol";


contract StakingManagerExample is StakingManager {

    function slash(address staker, uint256 amount) public {
        _transfer(staker, amount);
    }

}
