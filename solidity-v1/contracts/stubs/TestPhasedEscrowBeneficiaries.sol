pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

// Simple beneficiary that does nothing when notified that it has received
// tokens.
contract TestSimpleBeneficiary {
    function __escrowSentTokens(uint256 amount) external {}
}

// CurveRewards contract mock for curvefi.
contract TestCurveRewards {
    using SafeERC20 for IERC20;

    IERC20 public token;

    event RewardAdded(uint256 reward);

    constructor(IERC20 _token) public {
        token = _token;
    }

    function notifyRewardAmount(uint256 reward) external {
        token.safeTransferFrom(msg.sender, address(this), reward);
        emit RewardAdded(reward);
    }
}

// Simple reward contract mock for testing purposes.
contract TestSimpleStakerRewards {
    using SafeERC20 for IERC20;

    IERC20 public token;

    constructor(IERC20 _token) public {
        token = _token;
    }

    function receiveApproval(
        address _from,
        uint256 _value,
        address _token,
        bytes memory
    ) public {
        token.safeTransferFrom(_from, address(this), _value);
    }
}
