// SPDX-License-Identifier: MIT
/*
▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓▌        ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
  ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓

                           Trust math, not hardware.
*/
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "../MaintenancePool.sol";

// TODO: Documentation
// TODO: Unit tests
library Relay {
    using SafeERC20 for IERC20;

    struct Request {
        uint256 id;
        uint256 startBlock;
        bytes groupPublicKey;
        bytes previousEntry;
    }

    struct Data {
        uint256 requestCount;
        Request currentRequest;

        IERC20 tToken;
        MaintenancePool maintenancePool;
        uint256 relayRequestFee;

        uint256 groupSize;
        uint256 entrySubmissionBlockStep;
        uint256 bleedingPeriod;
    }

    event RelayEntryRequested(
        uint256 indexed requestId,
        bytes groupPublicKey,
        bytes previousEntry
    );
    event RelayEntrySubmitted(uint256 indexed requestId, bytes entry);

    function requestEntry(
        Data storage self,
        bytes memory groupPublicKey,
        bytes calldata previousEntry
    ) internal {
        require(!isRequestInProgress(self), "Another relay request in progress");

        // TODO: Transfer can be done directly to the maintenance pool. In that
        //       case the requester needs to know and approve the maintenance
        //       pool address and we lose the possibility to make additional
        //       actions upon deposit within the maintenance pool. We need to
        //       discuss the approach here.
        self.tToken.safeTransferFrom(msg.sender, address(this), self.relayRequestFee);
        self.tToken.safeIncreaseAllowance(address(self.maintenancePool), self.relayRequestFee);
        self.maintenancePool.deposit(self.relayRequestFee);

        uint256 currentRequestId = ++self.requestCount;

        self.currentRequest = Request(
            currentRequestId,
            block.number,
            groupPublicKey,
            previousEntry
        );

        emit RelayEntryRequested(currentRequestId, groupPublicKey, previousEntry);
    }

    function submitEntry(
        Data storage self,
        bytes calldata entry
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        require(!isRequestTimedOut(self), "Relay request timed out");

        // TODO: Verify submitter eligibility.

        // TODO: Attach BLS lib.
        // require(
        //     BLS.verify(
        //         self.currentRequest.groupPublicKey,
        //         self.currentRequest.previousEntry,
        //         entry
        //     ),
        //     "Invalid signature"
        // );

        // TODO: Implement bleeding.

        delete self.currentRequest;

        emit RelayEntrySubmitted(self.requestCount, entry);
    }

    function isRequestInProgress(
        Data storage self
    ) internal view returns (bool) {
        return self.currentRequest.id != 0;
    }

    function isRequestTimedOut(
        Data storage self
    ) internal view returns (bool) {
        // TODO: Bleeding period should be counted in.
        return isRequestInProgress(self) &&
            block.number > self.groupSize * self.entrySubmissionBlockStep;
    }
}