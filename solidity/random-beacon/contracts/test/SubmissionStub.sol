// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Submission.sol";

contract SubmissionStub {
    function getEligibilityRange(
        uint256 seed,
        uint256 protocolSubmissionBlock,
        uint256 protocolStartBlock,
        uint256 eligibilityDelay,
        uint256 groupSize
    )
        external
        view
        returns (uint256 firstEligibleIndex, uint256 lastEligibleIndex)
    {
        return
            Submission.getEligibilityRange(
                seed,
                protocolSubmissionBlock,
                protocolStartBlock,
                eligibilityDelay,
                groupSize
            );
    }

    function isEligible(
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        uint256 lastEligibleIndex
    ) external view returns (bool) {
        return
            Submission.isEligible(
                submitterIndex,
                firstEligibleIndex,
                lastEligibleIndex
            );
    }
}
