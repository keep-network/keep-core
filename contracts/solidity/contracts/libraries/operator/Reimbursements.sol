pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "../../TokenStaking.sol";

library Reimbursements {
    using SafeMath for uint256;
    using BytesLib for bytes;

    /**
     * @dev Reimburses callback execution cost and surplus based on actual gas
     * usage to the submitter's beneficiary address and if necessary to the
     * callback requestor (surplus recipient).
     * @param stakingContract Staking contract to get the address of the beneficiary
     * @param gasLimit Gas limit set for the callback
     * @param gasSpent Gas spent by the submitter on the callback
     * @param callbackFee Fee paid for the callback by the requestor
     * @param callbackReturnData Data containing surplus recipient address
     */
    function reimburseCallback(
        TokenStaking stakingContract,
        uint256 priceFeedEstimate,
        uint256 gasLimit,
        uint256 gasSpent,
        uint256 callbackFee,
        bytes memory callbackReturnData
    ) public {
        uint256 gasPrice = tx.gasprice < priceFeedEstimate ? tx.gasprice : priceFeedEstimate;

        // Obtain the actual callback gas expenditure and refund the surplus.
        //
        // In case of heavily underpriced transactions, EVM may wrap the call
        // with additional opcodes. In this case gasSpent > gasLimit.
        // The worst scenario cost is included in entry verification fee.
        // If this happens we return just the gasLimit here.
        uint256 actualCallbackGas = gasSpent < gasLimit ? gasSpent : gasLimit;
        uint256 actualCallbackFee = actualCallbackGas.mul(gasPrice);

        // Get the beneficiary.
        address payable magpie = stakingContract.magpieOf(msg.sender);

        // If we spent less on the callback than the customer transferred for the
        // callback execution, we need to reimburse the difference.
        if (actualCallbackFee < callbackFee) {
            uint256 callbackSurplus = callbackFee.sub(actualCallbackFee);
            // Reimburse submitter with his actual callback cost.
            magpie.call.value(actualCallbackFee)("");

            // Return callback surplus to the requestor.
            // Expecting 32 bytes data containing 20 byte address
            if (callbackReturnData.length == 32) {
                address surplusRecipient = callbackReturnData.toAddress(12);
                surplusRecipient.call.gas(8000).value(callbackSurplus)("");
            }
        } else {
            // Reimburse submitter with the callback payment sent by the requestor.
            magpie.call.value(callbackFee)("");
        }
    }
}
