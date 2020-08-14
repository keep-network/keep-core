pragma solidity 0.5.17;

import "../GasPriceOracle.sol";

contract GasPriceOracleConsumerStub is GasPriceOracleConsumer {
    GasPriceOracle gasPriceOracle;

    uint256 public gasPrice;

    constructor(GasPriceOracle _gasPriceOracle) public {
        gasPriceOracle = _gasPriceOracle;
    }

    function refreshGasPrice() public {
        gasPrice = gasPriceOracle.gasPrice();
    }
}