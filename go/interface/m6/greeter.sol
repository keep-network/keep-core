pragma solidity ^0.4.18;

contract Greeter {
  event _Greet(string name);

  function greet(string name) public {
    _Greet(name);
  }
}
