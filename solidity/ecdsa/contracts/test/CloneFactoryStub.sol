pragma solidity ^0.8.6;

import "@thesis/solidity-contracts/contracts/clone/CloneFactory.sol";

contract CloneFactoryStub is CloneFactory {
    address public masterContract;

    constructor(address _masterContract) {
        masterContract = _masterContract;
    }

    function createClone() external returns (address result) {
        return createClone(masterContract);
    }
}
