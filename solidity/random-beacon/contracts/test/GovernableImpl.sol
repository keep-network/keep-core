pragma solidity ^0.8.9;

import "../Governable.sol";

contract GovernableImpl is Governable {
    function _transferGovernanceExposed(address newGovernance) external {
        _transferGovernance(newGovernance);
    }
}
