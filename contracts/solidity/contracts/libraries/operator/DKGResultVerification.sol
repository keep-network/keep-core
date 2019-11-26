pragma solidity ^0.5.4;

import "solidity-bytes-utils/contracts/BytesLib.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "./GroupSelection.sol";

library DKGResultVerification {
    using BytesLib for bytes;
    using ECDSA for bytes32;
    using GroupSelection for GroupSelection.Storage;

    struct Storage {
        // Time in blocks after DKG result is complete and ready to be published
        // by clients.
        uint256 timeDKG;
    }

    /**
    * @dev Verifies that submitter is eligible to submit the result and that the provided
    * members signatures of the DKG result were produced by the members stored previously
    * on-chain in the order of their ticket values and returns indices of members with a
    * boolean value of their signature validity.
    * @param submitterMemberIndex Claimed index of the staker. We pass this for gas efficiency purposes.
    * @param signatures Concatenation of user-generated signatures.
    * @param resultHash The result hash signed by the users.
    * @param signingMemberIndices Indices of members corresponding to each signature.
    * @param members Array of selected participants.
    * @param groupThreshold Minimum number of members needed to produce a relay entry.
    * @param resultPublicationBlockStep Time in blocks after which the next group member is eligible to submit the result.
    * @param groupSelectionEndBlock Block height at which the group selection ended.
    * @return true if submitter is eligible to submit and the signatures are valid.
    */
    function verify(
        Storage storage self,
        uint256 submitterMemberIndex,
        bytes memory signatures,
        bytes32 resultHash,
        uint256[] memory signingMemberIndices,
        address[] memory members,
        uint256 groupThreshold,
        uint256 resultPublicationBlockStep,
        uint256 groupSelectionEndBlock
    ) public view returns (bool) {
        require(submitterMemberIndex > 0, "Invalid submitter index");
        require(
            members[submitterMemberIndex - 1] == msg.sender,
            "Unexpected submitter index"
        );

        uint T_init = groupSelectionEndBlock + self.timeDKG;
        require(
            block.number >= (T_init + (submitterMemberIndex-1) * resultPublicationBlockStep),
            "Submitter not eligible"
        );

        uint256 signaturesCount = signatures.length / 65;
        require(signatures.length >= 65, "Too short signatures array");
        require(signatures.length % 65 == 0, "Malformed signatures array");
        require(signaturesCount == signingMemberIndices.length, "Unexpected signatures count");
        require(signaturesCount >= groupThreshold, "Too few signatures");

        bytes memory current; // Current signature to be checked.

        for(uint i = 0; i < signaturesCount; i++){
            require(signingMemberIndices[i] > 0, "Invalid index");
            require(signingMemberIndices[i] <= members.length, "Index out of range");
            current = signatures.slice(65*i, 65);
            address recoveredAddress = resultHash.toEthSignedMessageHash().recover(current);

            require(members[signingMemberIndices[i] - 1] == recoveredAddress, "Invalid signature");
        }

        return true;
    }
}
