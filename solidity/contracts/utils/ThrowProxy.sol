pragma solidity 0.5.17;
// Proxy contract for testing throws
// http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests

contract ThrowProxy {
    address public target;
    bytes data;

    constructor(address _target) public {
        target = _target;
    }

    //prime the data using the fallback function.
    function() external payable {
        data = msg.data;
    }

    function execute() public returns (bool) {
        // The contract is used only for tests, disabling ethlint warning.
        // solium-disable-next-line
        (bool result, ) = target.call(data);
        return result;
    }
}
