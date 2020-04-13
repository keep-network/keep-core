pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

import "./utils/BytesLib.sol";

import "./KeepToken.sol";
import "./ManagedGrant.sol";
import "./TokenGrant.sol";
import "./GrantStakingPolicy.sol";

/// @title ManagedGrantFactory
/// @dev Creates managed grants that permit grantee reassignment
/// and use pre-defined staking policies.
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
        (address _grantee,
         uint256 _duration,
         uint256 _start,
         uint256 _cliff,
         bool _revocable) = abi.decode(
             _extraData,
             (address, uint256, uint256, uint256, bool)
        );
        _createGrant(
            _grantee,
            _amount,
            _duration,
            _start,
            _cliff,
            _revocable,
            _from
        );
    }

    function createManagedGrant(
        address grantee,
        uint256 amount,
        uint256 duration,
        uint256 start,
        uint256 cliff,
        bool revocable
    ) public returns (address _managedGrant) {
        return _createGrant(
            grantee,
            amount,
            duration,
            start,
            cliff,
            revocable,
            msg.sender
        );
    }

    function _createGrant(
        address grantee,
        uint256 amount,
        uint256 duration,
        uint256 start,
        uint256 cliff,
        bool revocable,
        address _from
    ) internal returns (address _managedGrant) {
        require(grantee != address(0), "Grantee address can't be zero.");
        require(cliff <= duration, "Unlocking cliff duration must be less or equal total unlocking duration.");

        token.safeTransferFrom(_from, address(this), amount);

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

        bytes memory grantData = abi.encode(
            _managedGrant,
            duration,
            start,
            cliff,
            revocable,
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
