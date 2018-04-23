pragma solidity ^0.4.18;
// Proxy contract for testing throws
// http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests

contract ThrowProxy {
  address public target;
  bytes data;

  function ThrowProxy(address _target) public {
    target = _target;
  }

  //prime the data using the fallback function.
  function() public {
    data = msg.data;
  }

  function execute() public returns (bool) {
    return target.call(data);
  }
}
