/**
▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓▌        ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
  ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓

                           Trust math, not hardware.
*/

pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

import "./libraries/grant/UnlockingSchedule.sol";
import "./TokenGrant.sol";
import "./ManagedGrant.sol";

/// @title TokenStakingEscrow
/// @notice Escrow lets the staking contract to deposit undelegated, granted
/// tokens and either withdraw them based on the grant unlocking schedule or
/// re-delegate them to another operator.
/// @dev The owner of TokenStakingEscrow is TokenStaking contract and only owner
/// can deposit. This contract works with an assumption that operator is unique.
contract TokenStakingEscrow is Ownable {

    using SafeERC20 for IERC20;
    using SafeMath for uint256;
    using UnlockingSchedule for uint256;

    event Deposited(
        address indexed operator,
        uint256 amount
    );
    event DepositWithdrawn(
        address indexed operator,
        address indexed grantee,
        uint256 amount
    );

    IERC20 public keepToken;
    TokenGrant public tokenGrant;

    struct Deposit {
        uint256 grantId;
        uint256 amount;
        uint256 withdrawn;
    }

    // operator address -> KEEP deposit
    mapping(address => Deposit) internal deposits;

    constructor(
        IERC20 _keepToken,
        TokenGrant _tokenGrant
    ) public {
        keepToken = _keepToken;
        tokenGrant = _tokenGrant;
    }

    /// @notice receiveApproval accepts deposits from staking contract and
    /// stores them in the escrow by the operator address from which they were
    /// undelegated. Function expects operator address and grant identifier to
    /// be passed as ABI-encoded information in extraData. Grant with the given
    /// identifier has to exist.
    function receiveApproval(
        address from,
        uint256 value,
        address token,
        bytes memory extraData
    ) public {
        require(IERC20(token) == keepToken, "Not a KEEP token");
        require(extraData.length == 64, "Unexpected data length");

        (address operator, uint256 grantId) = abi.decode(
            extraData, (address, uint256)
        );
        receiveDeposit(from, value, operator, grantId);
    }

    // TODO: redelegateTo(fromOperator, toOperator)

    /// @notice Returns the total amount deposited in the escrow after
    /// undelegating it from the provided operator.
    function depositedAmount(address operator) public view returns (uint256) {
        return deposits[operator].amount;
    }

    /// @notice Returns grant ID for the amount deposited in the escrow after
    /// undelegating it from the provided operator.
    function depositGrantId(address operator) public view returns (uint256) {
        return deposits[operator].grantId;
    }

    /// @notice Returns the amount withdrawn so far from the value deposited
    /// in the escrow contract after undelegating it from the provided operator.
    function depositWithdrawnAmount(address operator) public view returns (uint256) {
        return deposits[operator].withdrawn;
    }

    /// @notice Returns the currently withdrawable amount that was previously
    /// deposited in the escrow after undelegating it from the provided operator.
    /// Tokens are unlocked base on their grant unlocking schedule.
    /// Function returns 0 for non-existing deposits and revoked grants.
    function withdrawable(address operator) public view returns (uint256) {
        Deposit memory deposit = deposits[operator];

        // Staked tokens can be only withdrawn by grantee for non-revoked grant.
        // It is not possible for the escrow to determine the number of tokens
        // it should return to the grantee of a revoked grant given different
        // possible staking contracts and staking policies.
        if (getAmountRevoked(deposit.grantId) == 0) {
            (
                uint256 duration,
                uint256 start,
                uint256 cliff
            ) = getUnlockingSchedule(deposit.grantId);

            uint256 unlocked = now.getUnlockedAmount(
                deposit.amount,
                duration,
                start,
                cliff
            );

            if (deposit.withdrawn < unlocked) {
              return unlocked - deposit.withdrawn;
            }
        }

        return 0;
    }

    /// @notice Withdraws currently unlocked tokens deposited in the escrow
    /// after undelegating them from the provided operator. Only grantee or
    /// operator can call this function. Important: this function can not be
    /// called for a `ManagedGrant` grantee. This may lead to locking tokens.
    /// For `ManagedGrant`, please use `withdrawToManagedGrantee` instead.
    function withdraw(address operator) public {
        Deposit memory deposit = deposits[operator];
        address grantee = getGrantee(deposit.grantId);

        // Make sure this function is not called for a managed grant.
        // If called for a managed grant, tokens could be locked there.
        // Better be safe than sorry.
        (bool success, ) = address(this).call(
            abi.encodeWithSignature("getManagedGrantee(address)", grantee)
        );
        require(!success, "Can not be called for managed grant");

        require(
            msg.sender == grantee || msg.sender == operator,
            "Only grantee or operator can withdraw"
        );

        withdraw(deposit, operator, grantee);
    }

    /// @notice Withdraws currently unlocked tokens deposited in the escrow
    /// after undelegating them from the provided operator. Only grantee or
    /// operator can call this function. This function works only for
    /// `ManagedGrant` grantees. For a standard grant, please use `withdraw`
    /// instead.
    function withdrawToManagedGrantee(address operator) public {
        Deposit memory deposit = deposits[operator];
        address managedGrant = getGrantee(deposit.grantId);
        address grantee = getManagedGrantee(managedGrant);

        require(
            msg.sender == grantee || msg.sender == operator,
            "Only grantee or operator can withdraw"
        );

        withdraw(deposit, operator, grantee);
    }

    function getManagedGrantee(
        address managedGrantee
    ) internal view returns(address) {
        ManagedGrant grant = ManagedGrant(managedGrantee);
        return grant.grantee();
    }

    function receiveDeposit(
        address from,
        uint256 value,
        address operator,
        uint256 grantId
    ) internal {
        // This contract works with an assumption that operator is unique.
        // This is fine as long as the staking contract works with the same
        // assumption so we are limiting deposits to the staking contract only.
        require(from == owner(), "Only owner can deposit");
        require(
            getAmountGranted(grantId) > 0,
            "Grant with this ID does not exist"
        );

        keepToken.safeTransferFrom(from, address(this), value);
        deposits[operator] = Deposit(grantId, value, 0);

        emit Deposited(operator, value);
    }

    function withdraw(
        Deposit memory deposit,
        address operator,
        address grantee
    ) internal {
        uint256 amount = withdrawable(operator);

        deposits[operator].withdrawn = deposit.withdrawn.add(amount);
        keepToken.safeTransfer(grantee, amount);

        emit DepositWithdrawn(operator, grantee, amount);
    }

    function getAmountGranted(uint256 grantId) internal view returns (
        uint256 amountGranted
    ) {
       (amountGranted,,,,,) = tokenGrant.getGrant(grantId);
    }

    function getAmountRevoked(uint256 grantId) internal view returns (
        uint256 amountRevoked
    ) {
        (,,,amountRevoked,,) = tokenGrant.getGrant(grantId);
    }

    function getUnlockingSchedule(uint256 grantId) internal view returns (
        uint256 duration,
        uint256 start,
        uint256 cliff
    ) {
        (,duration,start,cliff,) = tokenGrant.getGrantUnlockingSchedule(grantId);
    }

    function getGrantee(uint256 grantId) internal view returns (
        address grantee
    ) {
        (,,,,,grantee) = tokenGrant.getGrant(grantId);
    }
}