pragma solidity 0.5.17;

import "../../TokenGrant.sol";
import "../RolesLookup.sol";

library GrantStaking {

    using RolesLookup for address payable;

    /// @dev Grant ID is flagged with the most significant bit set, to
    /// distinguish the grant ID `0` from default (null) value. The flag is
    /// toggled with bitwise XOR (`^`) which keeps all other bits intact but
    /// flips the flag bit. The flag should be set before writing to
    /// `operatorToGrant`, and unset after reading from `operatorToGrant`
    /// before using the value.
    uint256 constant GRANT_ID_FLAG = 1 << 255;

    struct Storage {
        /// @dev Do not read or write this mapping directly; please use
        /// `hasGrantDelegated`, `setGrantForOperator`, and `getGrantForOperator`
        /// instead.
        mapping (address => uint256) _operatorToGrant;
    }

    function tryCapturingGrantId(
        Storage storage self,
        TokenGrant tokenGrant,
        address operator
    ) public {
        (bool success, bytes memory data) = address(tokenGrant).call(
            abi.encodeWithSignature("getGrantStakeDetails(address)", operator)
        );
        if (success) {
            uint256 grantId = abi.decode(data, (uint256));
            setGrantForOperator(self, operator, grantId);
        }
    }

    function hasGrantDelegated(
        Storage storage self,
        address operator
    ) public view returns (bool) {
        return self._operatorToGrant[operator] != 0;
    }

    function setGrantForOperator(
        Storage storage self,
        address operator,
        uint256 grantId
    ) public {
        self._operatorToGrant[operator] = grantId ^ GRANT_ID_FLAG;
    }

    function getGrantForOperator(
        Storage storage self,
        address operator
    ) public view returns (uint256) {
        uint256 grantId = self._operatorToGrant[operator];
        require (grantId != 0, "No grant for the operator");
        return grantId ^ GRANT_ID_FLAG;
    }

    function canUndelegate(
        Storage storage self,
        address owner,
        address operator,
        TokenGrant tokenGrant
    ) public returns (bool) {
        // First, check three simple cases:
        // - msg.sender is the operator,
        // - msg.sender is the owner (liquid tokens case or TokenGrantStaking
        //   contract is calling),
        // - msg.sender is a grantee of a standard grant.
        //
        // If one of them is the case, we return true.
        if (
            msg.sender == operator ||
            msg.sender == owner ||
            (msg.sender).isGranteeForOperator(operator, tokenGrant)
        ) {
            return true;
        }

        // If none of the above is true, we need to dig deeper and see if we
        // are dealing with a grantee from a managed grant.
        //
        // First of all, we need to see if the operator has grant delegated.
        // If it does not, there is no chance it's a managed grantee calling.
        if (!hasGrantDelegated(self, operator)) {
            return false;
        }

        // We know the operator has grant delegated, we are going to retrieve
        // the grant ID and check if msg.sender is grantee from a managed grant.
        uint256 grantId = getGrantForOperator(self, operator);
        if ((msg.sender).isManagedGranteeForOperatorAndGrant(
            operator, grantId, tokenGrant
        )) {
            return true;
        }

        // There is only one possibility left - grant has been revoked and
        // grant manager wants to take back delegated tokens.
        (,,,,uint256 revokedAt,) = tokenGrant.getGrant(grantId);
        if (revokedAt == 0) {
            return false;
        }
        (address grantManager,,,,) = tokenGrant.getGrantUnlockingSchedule(grantId);
        return msg.sender == grantManager;
    }
}