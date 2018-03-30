// Proxy contract for testing throws
// http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests

contract ThrowProxy {
  address public target;
  bytes data;

  function ThrowProxy(address _target) {
    target = _target;
  }

  //prime the data using the fallback function.
  function() {
    data = msg.data;
  }

  function execute() returns (bool) {
    return target.call(data);
  }
}
