pragma solidity ^0.5.4;

/**
 * @title Keep Random Beacon
 *
 * @notice Keep Random Beacon generates verifiable randomness that is resistant
 * to bad actors both in the relay network and on the anchoring blockchain.
 */
interface IRandomBeacon {
    /**
     * @notice Event is emitted when a new relay entry is generated for specific
     * relay request ID.
     *
     * @param requestId Relay request ID for which entry was generated.
     * @param entry Generated relay entry.
     */
    event RelayEntryGenerated(uint256 requestId, uint256 entry);

    /**
     * @notice Provides the customer with an estimated entry fee in wei to use
     * in the request. The fee estimate is only valid for the transaction it is
     * called in, so the customer must make the request immediately after
     * obtaining the estimate. Insufficient payment will lead to the request
     * being rejected and the transaction reverted.
     *
     * The customer may decide to provide more ether for an entry fee than
     * estimated by this function. This is especially helpful when callback gas
     * cost fluctuates. Any surplus between the passed fee and the actual cost
     * of producing an entry and executing a callback is returned back to the
     * customer.
     * @param callbackGas Gas required for the callback.
     */
    function entryFeeEstimate(uint256 callbackGas)
        external
        view
        returns (uint256);

    /**
     * @notice Submits a request to generate a new relay entry. Executes the
     * provided callback with the generated entry and emits
     * `RelayEntryGenerated(uint256 requestId, uint256 entry)` event.
     *
     * @dev Beacon does not support concurrent relay requests. No new requests
     * should be made while the beacon is already processing another request.
     * Requests made while the beacon is busy will be rejected and the
     * transaction reverted.

     * @param callbackContract Callback contract address. Callback is called
     * once a new relay entry has been generated.
     * @param callbackMethod Callback contract method signature. String
     * representation of your method with a single
     * uint256 input parameter i.e. "relayEntryCallback(uint256)".
     * @param callbackGas Gas required for the callback.
     * The customer needs to ensure they provide a sufficient callback gas
     * to cover the gas fee of executing the callback. Any surplus is returned
     * to the customer. If the callback gas amount turns to be not enough to
     * execute the callback, callback execution is skipped.
     * @return An uint256 representing uniquely generated relay request ID
     */
    function requestRelayEntry(
        address callbackContract,
        string calldata callbackMethod,
        uint256 callbackGas
    ) external payable returns (uint256);

    /**
     * @notice Submits a request to generate a new relay entry. Emits
     * `RelayEntryGenerated(uint256 requestId, uint256 entry)` event for the
     * generated entry.
     *
     * @dev Beacon does not support concurrent relay requests. No new requests
     * should be made while the beacon is already processing another request.
     * Requests made while the beacon is busy will be rejected and the
     * transaction reverted.
     *
     * @return An uint256 representing uniquely generated relay request ID
     */
    function requestRelayEntry() external payable returns (uint256);
}
