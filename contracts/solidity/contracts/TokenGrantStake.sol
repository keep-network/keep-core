pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./TokenStaking.sol";
import "./utils/BytesLib.sol";

/**
   @dev Interface of sender contract for approveAndCall pattern.
*/
interface tokenSender {
    function approveAndCall(address _spender, uint256 _value, bytes calldata _extraData) external;
}

contract TokenGrantStake {
    using SafeMath for uint256;
    using BytesLib for bytes;

    ERC20Burnable token;
    address tokenGrant; // Address of the master grant contract.
    uint256 grantId; // ID of the grant for this stake.
    TokenStaking stakingContract; // Staking contract.
    uint256 amount; // Amount of staked tokens.
    address operator; // Operator of the stake.

    constructor(
        address _tokenAddress,
        uint256 _grantId,
        address _stakingContract
    ) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = ERC20Burnable(_tokenAddress);
        tokenGrant = msg.sender;
        grantId = _grantId;
        stakingContract = TokenStaking(_stakingContract);
    }

    function stake(
        uint256 _amount,
        bytes memory _extraData
    ) onlyGrant public {
        amount = _amount;
        operator = _extraData.toAddress(20);
        tokenSender(address(token)).approveAndCall(
            address(stakingContract),
            _amount,
            _extraData
        );
    }

    function getGrantId() onlyGrant public view returns (uint256) {
        return grantId;
    }

    function getAmount() onlyGrant public view returns (uint256) {
        return amount;
    }

    function getStakingContract() onlyGrant public view returns (address) {
        return address(stakingContract);
    }

    function getDetails() onlyGrant public view returns (
        uint256 _grantId,
        uint256 _amount,
        address _stakingContract
    ) {
        return (
            grantId,
            amount,
            address(stakingContract)
        );
    }

    function cancelStake() onlyGrant public returns (uint256) {
        stakingContract.cancelStake(operator);
        return returnTokens();
    }

    function undelegate() onlyGrant public {
        stakingContract.undelegate(operator);
    }

    function recoverStake() onlyGrant public returns (uint256) {
        stakingContract.recoverStake(operator);
        return returnTokens();
    }

    function returnTokens() onlyGrant public returns (uint256) {
        uint256 returnedAmount = token.balanceOf(address(this));
        amount -= returnedAmount;
        token.transfer(tokenGrant, returnedAmount);
        return returnedAmount;
    }

    modifier onlyGrant {
        require(
            msg.sender == tokenGrant,
            "For token grant contract only"
        );
        _;
    }
}
