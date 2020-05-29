pragma solidity 0.5.17;

import "../utils/AddressArrayUtils.sol";
import "../TokenStaking.sol";
import "../TokenGrant.sol";

interface ManagedGrant {
    function grantee() external view returns(address);
}

library RolesLookup {

    using AddressArrayUtils for address[];

    function isTokenOwnerForOperator(
        address tokenOwner,
        address operator,
        TokenStaking tokenStaking
    ) public view returns (bool) {
        return tokenStaking.ownerOf(operator) == tokenOwner;
    }

    function isGranteeForOperator(
        address grantee,
        address operator,
        TokenGrant tokenGrant
    ) public view returns (bool) {
        address[] memory operators = tokenGrant.getGranteeOperators(grantee);
        return operators.contains(operator);
    }

    function isManagedGranteeForOperator(
        address grantee,
        address operator,
        address granteeContract,
        TokenGrant tokenGrant
    ) public view returns (bool) {
        require(
            ManagedGrant(granteeContract).grantee() == grantee,
            "Not a grantee of the provided contract"
        );

        address[] memory operators = tokenGrant.getGranteeOperators(
            granteeContract
        );
        return operators.contains(operator);
    }
}