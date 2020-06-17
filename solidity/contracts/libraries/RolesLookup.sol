pragma solidity 0.5.17;

import "../utils/AddressArrayUtils.sol";
import "../StakeDelegatable.sol";
import "../TokenGrant.sol";
import "../ManagedGrant.sol";

/// @title Roles Lookup
/// @notice Library facilitating lookup of roles in stake delegation setup.
library RolesLookup {
    using AddressArrayUtils for address[];

    /// @notice Returns true if the tokenOwner delegated tokens to operator
    /// using the provided tokenStaking contract. Othwerwise, returns false.
    /// This function works only for the case when tokenOwner own those tokens
    /// and those are not tokens from a grant.
    function isTokenOwnerForOperator(
        address tokenOwner,
        address operator,
        StakeDelegatable tokenStaking
    ) internal view returns (bool) {
        return tokenStaking.ownerOf(operator) == tokenOwner;
    }

    /// @notice Returns true if the grantee delegated tokens to operator
    /// with the provided tokenGrant contract. Otherwise, returns false.
    /// This function works only for the case when tokens were generated from
    /// a non-managed grant, that is, the grantee is a non-contract address to
    /// which the delegated tokens were granted.
    function isGranteeForOperator(
        address grantee,
        address operator,
        TokenGrant tokenGrant
    ) internal view returns (bool) {
        address[] memory operators = tokenGrant.getGranteeOperators(grantee);
        return operators.contains(operator);
    }

    /// @notice Returns true if the grantee from the given managed grant contract
    /// delegated tokens to operator with the provided tokenGrant contract.
    /// Otherwise, returns false. In case the grantee declared by the managed
    /// grant contract does not match the provided grantee, function reverts.
    /// This function works only for cases when grantee, from TokenGrant's
    /// perspective, is a smart contract exposing grantee() function returning
    /// the final grantee. One possibility is the ManagedGrant contract.
    function isManagedGranteeForOperator(
        address grantee,
        address operator,
        address managedGrantContract,
        TokenGrant tokenGrant
    ) internal view returns (bool) {
        require(
            ManagedGrant(managedGrantContract).grantee() == grantee,
            "Not a grantee of the provided contract"
        );

        address[] memory operators = tokenGrant.getGranteeOperators(
            managedGrantContract
        );
        return operators.contains(operator);
    }
}
