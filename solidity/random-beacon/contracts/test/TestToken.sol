// SPDX-License-Identifier: MIT

pragma solidity 0.8.5;

import "@thesis/solidity-contracts/contracts/token/ERC20WithPermit.sol";

contract TestToken is ERC20WithPermit {
    /* solhint-disable-next-line no-empty-blocks */
    constructor() ERC20WithPermit("Test Token", "TT") {}
}
