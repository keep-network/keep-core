// SPDX-License-Identifier: GPL-3.0-only

pragma solidity 0.8.17;

import "../WalletRegistry.sol";
import "../libraries/EcdsaDkg.sol";

contract DkgChallenger {
    WalletRegistry internal walletRegistry;

    constructor(WalletRegistry _walletRegistry) {
        walletRegistry = _walletRegistry;
    }

    function challengeDkgResult(EcdsaDkg.Result calldata dkgResult) external {
        walletRegistry.challengeDkgResult(dkgResult);
    }
}
