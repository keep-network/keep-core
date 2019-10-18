pragma solidity ^0.5.4;

contract MyTestContractOperator {

  event MyEvent(uint256 value);
  
  function functionOne() public {
      uint256 val = tx.gasprice;
      emit MyEvent(val);
  }

  function functionTwo() public {
      uint256 val = 20 * 1e9;
      emit MyEvent(val);
  }

  function functionThree() public {
      uint256 val = 500;
      emit MyEvent(val);
  }

  uint256 public priceFeedEstimate = 20*1e9;

  uint256 public dkgGasEstimate = 2260000;

  function functionFour() public {
      uint256 gasPrice = tx.gasprice < priceFeedEstimate ? tx.gasprice : priceFeedEstimate;
      emit MyEvent(gasPrice);
  }
}
