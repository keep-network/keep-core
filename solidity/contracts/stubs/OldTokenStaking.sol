pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
import "../utils/OperatorParams.sol";

/// Staking contract stub for testing purposes that mimics the behavior of:
/// - v1.0.1 TokenStaking (Mainnet)
/// - 1.3.0-rc.0 TokenStaking (Ropsten)
contract OldTokenStaking {
    using OperatorParams for uint256;

    event Undelegated(address indexed operator, uint256 undelegatedAt);

    mapping(address => address[]) public ownerOperators;
    mapping(address => Operator) public operators;

    struct Operator {
        uint256 packedParams;
        address owner;
        address payable beneficiary;
        address authorizer;
    }

    function operatorsOf(address _address) public view returns (address[] memory) {
        return ownerOperators[_address];
    }

    function balanceOf(address _address) public view returns (uint256 balance) {
        return operators[_address].packedParams.getAmount();
    }

    function ownerOf(address _operator) public view returns (address) {
        return operators[_operator].owner;
    }

    function beneficiaryOf(address _operator) public view returns (address payable) {
        return operators[_operator].beneficiary;
    }

    function authorizerOf(address _operator) public view returns (address) {
        return operators[_operator].authorizer;
    }

    function getDelegationInfo(address _operator)
    public view returns (uint256 amount, uint256 createdAt, uint256 undelegatedAt) {
        return operators[_operator].packedParams.unpack();
    }

    function undelegationPeriod() public view returns(uint256) {
        return 5184000; // two months
    }

    function undelegate(address _operator) public {
        uint256 oldParams = operators[_operator].packedParams;
        operators[_operator].packedParams = oldParams.setUndelegationTimestamp(
            block.timestamp
        );
        emit Undelegated(_operator, block.timestamp);
    }

    function setOperator(
        address _owner,
        address _operator,
        address payable _beneficiary,
        address _authorizer,
        uint256 _value
    ) public {
        operators[_operator] = Operator(
            OperatorParams.pack(_value, block.timestamp, 0),
            _owner,
            _beneficiary,
            _authorizer
        );
        ownerOperators[_owner].push(_operator);
    }
}
