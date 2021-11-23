// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

library Authorization {

    struct Data {
       uint96 minimumAuthorization;
    }

    function setMinimumAuthorization(
        Data storage self,
        uint96 _minimumAuthorization
    ) internal {
        self.minimumAuthorization = _minimumAuthorization;
    }
}