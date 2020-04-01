pragma solidity ^0.5.4;

contract DelegatedAuthorityStub {
    address recognizedContract;

    constructor(address _recognizedContract) public {
        recognizedContract = _recognizedContract;
    }

    function __isRecognized(address _contract) public view returns (bool) {
        return _contract == recognizedContract;
    }
}
