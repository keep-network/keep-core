pragma solidity ^0.5.17;

import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

contract NewRewardsStub {

    using SafeERC20 for IERC20;

    function receiveApproval(
        address _from,
        uint256 _value,
        address _token,
        bytes memory
    ) public {
        IERC20(_token).safeTransferFrom(_from, address(this), _value);
    }
}