pragma solidity 0.5.17;

contract MigrationEscrowStub {

    struct Received {
        uint256 grantId;
        uint256 amount;
    }

    mapping (address => Received) public received;

    function receiveApproval(
        address from,
        uint256 value,
        address token,
        bytes memory extraData
    ) public {
        (address operator, uint256 grantId) = abi.decode(
            extraData, (address, uint256)
        );
        received[operator] = Received(grantId, value);
    }

    function depositedAmount(address operator) public view returns (uint256) {
        return received[operator].amount;
    }

    function depositGrantId(address operator) public view returns (uint256) {
        return received[operator].grantId;
    }
}