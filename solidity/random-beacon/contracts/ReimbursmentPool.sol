// SPDX-License-Identifier: MIT
//
// ▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
//   ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
// ▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
//
//                           Trust math, not hardware.

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/access/Ownable.sol";

contract ReimbursementPool is Ownable {
    /// @notice Authorized contracts that can interact with the reimbursment pool.
    ///         Authorization can be granted and removed by the governance.
    mapping(address => bool) public isAuthorized;

    /// @notice Static gas cost of submitting a transaction.
    uint256 public staticGasCost;

    /// @notice Max gas price used to reimburse a transaction submitter. Protects
    ///         against malicious operator-miners.
    uint256 public maxGasPrice;

    modifier onlyAuthorized() {
        require(
            isAuthorized[msg.sender],
            "Contract is not authorized for a refund"
        );
        _;
    }

    /// @notice Refunds ETH to a spender for executing specific transactions.
    /// @dev Only authorized contracts are allowed calling this function.
    /// @param gasSpent Gas spent on a transaction that needs to be reimbursed.
    /// @param spender Address where the reimbursment is sent.
    function refund(uint256 gasSpent, address spender) external onlyAuthorized {
        uint256 gasPrice = tx.gasprice < maxGasPrice
            ? tx.gasprice
            : maxGasPrice;
        uint256 refundAmount = (gasSpent + staticGasCost) * gasPrice;

        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = spender.call{value: refundAmount}("");
        require(sent, "Failed to refund Ether");
    }

    /// @notice Authorize a contract that can interact with this reimbursment pool.
    ///         Can be authorized by the governance only.
    /// @param _contract Authorized contract.
    function authorize(address _contract) external onlyOwner {
        isAuthorized[_contract] = true;
    }

    /// @notice Unauthorize a contract that was previously authorized to interact
    ///         with this reimbursment pool. Can be unauthorized by the
    ///         governance only.
    /// @param _contract Authorized contract.
    function unauthorize(address _contract) external onlyOwner {
        delete isAuthorized[_contract];
    }

    /// @notice Setting a static gas cost for executing a transaction. Can be set
    ///         by the governance only.
    /// @param _staticGasCost Static gas cost.
    function setStaticGasCost(uint256 _staticGasCost) external onlyOwner {
        staticGasCost = _staticGasCost;
    }

    /// @notice Setting a max gas price for transactions. Can be set by the
    ///         governance only.
    /// @param _maxGasPrice Max gas price used to reimburse tx submitters.
    function setMaxGasPrice(uint256 _maxGasPrice) external onlyOwner {
        maxGasPrice = _maxGasPrice;
    }

    /// @notice Withdraws ETH amount from this pool which are sent to a given
    ///         address. Can be set by the governance only.
    /// @param amount Amount to withdraw from the pool.
    /// @param to An address where ETH is sent.
    function withdraw(uint256 amount, address to) external onlyOwner {
        require(
            address(this).balance >= amount,
            "Insufficient contract balance"
        );

        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = to.call{value: amount}("");
        require(sent, "Failed to send Ether");
    }

    /// @notice Withdraws all ETH from this pool which are sent to a given
    ///         address. Can be set by the governance only.
    /// @param to An address where ETH is sent.
    function withdrawAll(address to) external onlyOwner {
        require(address(this).balance > 0, "Nothing to withdraw");

        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = to.call{value: address(this).balance}("");
        require(sent, "Failed to send Ether");
    }

    /// @notice Receive ETH
    receive() external payable {}
}
