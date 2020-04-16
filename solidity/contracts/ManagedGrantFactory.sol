pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

import "./utils/BytesLib.sol";

import "./KeepToken.sol";
import "./ManagedGrant.sol";
import "./TokenGrant.sol";
import "./GrantStakingPolicy.sol";

/// @title ManagedGrantFactory
/// @notice Creates managed grants that permit grantee reassignment
/// and use pre-defined staking policies.
contract ManagedGrantFactory {
    using SafeERC20 for KeepToken;
    using BytesLib for bytes;
    using BytesLib for address;

    KeepToken public token;
    TokenGrant public tokenGrant;
    GrantStakingPolicy nonRevocableStakingPolicy;
    GrantStakingPolicy revocableStakingPolicy;

    event ManagedGrantCreated(
        address grantAddress,
        address grantee
    );

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

    /// @notice Create a managed grant
    /// with the parameters specified in `_extraData`.
    /// @dev Requires no setup beforehand,
    /// but only provides the managed grant address through an event.
    /// @param _from The owner of the tokens who approved them to transfer.
    /// @param _amount Approved amount for the transfer to create the grant.
    /// @param _token Address of the token contract;
    /// must match the token specified when the factory was created.
    /// @param _extraData The following values encoded with `abi.encode`:
    /// grantee (address) Address of the grantee.
    /// duration (uint256) Duration in seconds of the unlocking period.
    /// cliff (uint256) Duration in seconds of the cliff before which no tokens will unlock.
    /// start (uint256) Timestamp at which unlocking will start.
    /// revocable (bool) Whether the token grant is revocable or not.
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

    /// @notice Create a managed grant with the given parameters.
    /// @dev At least `amount` tokens to be approved for the factory beforehand.
    /// The grant will use the staking policy specified for its type
    /// (revocable or non-revocable)
    /// when the factory was created.
    /// @param grantee The initial grantee.
    /// @param amount The number of tokens to grant.
    /// @param duration Duration in seconds of the unlocking period.
    /// @param start Timestamp at which unlocking will start.
    /// @param cliff Duration in seconds of the cliff before which no tokens will unlock.
    /// @param revocable Whether the token grant is revocable or not.
    /// @return The address of the managed grant.
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
            _from,
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

        emit ManagedGrantCreated(
            _managedGrant,
            grantee
        );
        return _managedGrant;
    }
}
