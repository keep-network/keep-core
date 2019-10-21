pragma solidity ^0.5.4;

contract MyTestContractOperator {

  
  
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

  event MyEvent(uint256 value);

  uint256 public averagePrice = 20*1e9;

  function functionFour() public {
      uint256 gasPrice = tx.gasprice < averagePrice ? tx.gasprice : averagePrice;
      emit MyEvent(gasPrice);
  }

  function functionFive() public {
      uint256 gasPrice = tx.gasprice;
      emit MyEvent(gasPrice);
  }
}
