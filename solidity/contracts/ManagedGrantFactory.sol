pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

import "./utils/BytesLib.sol";

import "./KeepToken.sol";
import "./ManagedGrant.sol";
import "./TokenGrant.sol";
import "./GrantStakingPolicy.sol";

contract ManagedGrantFactory {
    using SafeERC20 for KeepToken;
    using BytesLib for bytes;
    using BytesLib for address;

    KeepToken public token;
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
        token = KeepToken(_tokenAddress);
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

        token.safeTransferFrom(_from, address(this), _amount);

        grantFundingPool[_from] += _amount;
    }

    function createGrant(
        address grantee,
        uint256 amount,
        uint256 duration,
        uint256 start,
        uint256 cliff,
        bool revocable
    ) public returns (address _managedGrant) {
        require(
            grantFundingPool[msg.sender] >= amount,
            "Insufficient funding"
        );
        grantFundingPool[msg.sender] -= amount;

        GrantStakingPolicy stakingPolicy = revocable
            ? revocableStakingPolicy
            : nonRevocableStakingPolicy;

        // Grant ID is predictable in advance
        uint256 grantId = tokenGrant.numGrants();

        ManagedGrant managedGrant = new ManagedGrant(
            address(token),
            address(tokenGrant),
            msg.sender,
            grantId,
            grantee
        );
        _managedGrant = address(managedGrant);

        bytes memory grantData = abi.encodePacked(
            _managedGrant,
            duration,
            start,
            cliff,
            revocable ? 0x01 : 0x00,
            address(stakingPolicy)
        );

        token.approveAndCall(
            address(tokenGrant),
            amount,
            grantData
        );

        return _managedGrant;
    }
}
