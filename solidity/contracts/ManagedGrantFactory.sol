pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

import "./utils/BytesLib.sol";

import "./KeepToken.sol";
import "./ManagedGrant.sol";
import "./TokenGrant.sol";
import "./GrantStakingPolicy.sol";

contract ManagedGrantFactory {
    using SafeERC20 for ERC20Burnable;
    using BytesLib for bytes;
    using BytesLib for address;

    ERC20Burnable public token;
    mapping(address => uint256) public grantFundingPool;
    TokenGrant public tokenGrant;
    GrantStakingPolicy nonRevocableStakingPolicy;
    GrantStakingPolicy revocableStakingPolicy;

    constructor(
        address _tokenAddress,
        address _tokenGrant,
        address _nonRevocableStakingPolicy,
        address _revocableStakingPolicy
    ) public {
        token = ERC20Burnable(_tokenAddress);
        tokenGrant = TokenGrant(_tokenGrant);
        nonRevocableStakingPolicy = GrantStakingPolicy(_nonRevocableStakingPolicy);
        revocableStakingPolicy = GrantStakingPolicy(_revocableStakingPolicy);
    }

    function receiveApproval(
        address _from,
        uint256 _amount,
        address _token,
        bytes memory _extraData
    ) public {
        require(KeepToken(_token) == token, "Invalid token contract");
        require(_amount <= token.balanceOf(_from), "Insufficient sender balance");

        token.safeTransferFrom(_from, address(this), _amount);

        grantFundingPool[_from] += _amount;
    }
}
