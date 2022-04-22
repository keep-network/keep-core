pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@keep-network/random-beacon/contracts/ReimbursementPool.sol";
import "../WalletRegistry.sol";
import "../EcdsaDkgValidator.sol";

contract WalletRegistryStub is WalletRegistry {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(SortitionPool _sortitionPool, IStaking _staking)
        WalletRegistry(_sortitionPool, _staking)
    {}

    function getDkgData() external view returns (EcdsaDkg.Data memory) {
        return dkg;
    }
}
