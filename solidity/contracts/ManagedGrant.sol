pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "./TokenGrant.sol";

/// @title ManagedGrant
/// @dev A managed grant acts as the grantee towards the token grant contract,
/// proxying instructions from the actual grantee.
/// The address used by the actual grantee
/// to issue instructions and withdraw tokens
/// can be reassigned with the consent of the grant manager.
contract ManagedGrant {
    using SafeERC20 for ERC20Burnable;

    ERC20Burnable public token;
    TokenGrant public tokenGrant;
    address public grantManager;
    uint256 public grantId;
    address public grantee;
    address public requestedNewGrantee;

    constructor(
        address _tokenAddress,
        address _tokenGrant,
        address _grantManager,
        uint256 _grantId,
        address _grantee
    ) public {
        token = ERC20Burnable(_tokenAddress);
        tokenGrant = TokenGrant(_tokenGrant);
        grantManager = _grantManager;
        grantId = _grantId;
        grantee = _grantee;
    }

    function requestGranteeReassignment(address _newGrantee)
        public
        onlyGrantee
        noRequestedReassignment
    {
        _setRequestedNewGrantee(_newGrantee);
    }

    function cancelReassignmentRequest()
        public
        onlyGrantee
        withRequestedReassignment
    {
        requestedNewGrantee = address(0);
    }

    function changeReassignmentRequest(address _newGrantee)
        public
        onlyGrantee
        withRequestedReassignment
    {
        require(
            requestedNewGrantee != _newGrantee,
            "Unchanged reassignment request"
        );
        _setRequestedNewGrantee(_newGrantee);
    }

    function confirmGranteeReassignment(address _newGrantee)
        public
        onlyManager
        withRequestedReassignment
    {
        require(
            requestedNewGrantee == _newGrantee,
            "Reassignment address mismatch"
        );
        grantee = requestedNewGrantee;
        requestedNewGrantee = address(0);
    }

    function withdraw() public onlyGrantee {
        require(
            requestedNewGrantee == address(0),
            "Can not withdraw with pending reassignment"
        );
        tokenGrant.withdraw(grantId);
        uint256 amount = token.balanceOf(address(this));
        token.safeTransfer(grantee, amount);
    }

    function stake(
        address _stakingContract,
        uint256 _amount,
        bytes memory _extraData
    ) public onlyGrantee {
        tokenGrant.stake(grantId, _stakingContract, _amount, _extraData);
    }

    function cancelStake(address _operator) public onlyGranteeOr(_operator) {
        tokenGrant.cancelStake(_operator);
    }

    function undelegate(address _operator) public onlyGranteeOr(_operator) {
        tokenGrant.undelegate(_operator);
    }

    function recoverStake(address _operator) public onlyGrantee {
        tokenGrant.recoverStake(_operator);
    }

    function _setRequestedNewGrantee(address _newGrantee) internal {
        require(_newGrantee != address(0), "Invalid new grantee address");
        require(_newGrantee != grantee, "New grantee same as current grantee");

        requestedNewGrantee = _newGrantee;
    }

    modifier withRequestedReassignment {
        require(
            requestedNewGrantee != address(0),
            "No reassignment requested"
        );
        _;
    }

    modifier noRequestedReassignment {
        require(
            requestedNewGrantee == address(0),
            "Reassignment already requested"
        );
        _;
    }

    modifier onlyGrantee {
        require(
            msg.sender == grantee,
            "Only grantee may perform this action"
        );
        _;
    }

    modifier onlyGranteeOr(address _operator) {
        require(
            msg.sender == grantee || msg.sender == _operator,
            "Only grantee or operator may perform this action"
        );
        _;
    }

    modifier onlyManager {
        require(
            msg.sender == grantManager,
            "Only grantManager may perform this action"
        );
        _;
    }
}
