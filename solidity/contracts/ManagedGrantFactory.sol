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
        require(_amount <= token.balanceOf(_from), "Insufficient sender balance");

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

        bytes memory extraData;

        {
        bytes memory fstBytes;
        bytes memory sndBytes;
        bytes memory trdBytes;

        bytes memory managedGrantBytes = abi.encodePacked(_managedGrant);
        bytes memory durationBytes = abi.encodePacked(duration);
        fstBytes = BytesLib.concat(managedGrantBytes, durationBytes);

        bytes memory startBytes = abi.encodePacked(start);
        bytes memory cliffBytes = abi.encodePacked(cliff);
        sndBytes = BytesLib.concat(startBytes, cliffBytes);

        bytes memory revocableBytes = abi.encodePacked(
            uint8(revocable ? 0x01 : 0x00)
        );
        bytes memory stakingPolicyBytes = abi.encodePacked(address(stakingPolicy));
        trdBytes = BytesLib.concat(revocableBytes, stakingPolicyBytes);

        extraData = fstBytes.concat(sndBytes).concat(trdBytes);
        }

        token.approveAndCall(
            address(tokenGrant),
            amount,
            extraData
        );

        return _managedGrant;
    }
}
