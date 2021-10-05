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
import "../../../contracts/KeepToken.sol";

// TODO: Documentation
// TODO: Unit tests
library Relay {
    using SafeERC20 for IERC20;

    struct Request {
        uint256 id;
        uint256 startBlock;
        uint256 groupIndex;
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

    function requestEntry(
        Data storage self,
        uint256 groupIndex,
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

        self.requestCount++;

        self.currentRequest = Request(
            self.requestCount,
            block.number,
            groupIndex,
            previousEntry
        );

        // TODO: Emit event.
    }

    function submitEntry(
        Data storage self,
        bytes calldata entry
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        require(!isRequestTimedOut(self), "Relay request timed out");

        // TODO: Verify entry and submitter eligibility.
        // TODO: Implement bleeding.

        delete self.currentRequest;

        // TODO: Emit event.
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