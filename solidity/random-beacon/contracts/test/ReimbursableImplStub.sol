// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {Reimbursable} from "../Reimbursable.sol";

contract ReimbursableImplStub is Reimbursable {
    address public admin;

    constructor(address _admin) {
        admin = _admin;
    }

    modifier onlyReimbursableAdmin() override {
        require(admin == msg.sender, "Caller is not the admin");
        _;
    }
}
