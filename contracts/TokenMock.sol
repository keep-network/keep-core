// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import "@thesis/solidity-contracts/contracts/token/ERC20WithPermit.sol";

contract TokenMock is ERC20WithPermit {
    constructor() ERC20WithPermit("Token Mock", "T") {}
}
