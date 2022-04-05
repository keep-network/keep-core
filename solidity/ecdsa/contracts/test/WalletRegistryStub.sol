pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@keep-network/random-beacon/contracts/ReimbursementPool.sol";
import "../WalletRegistry.sol";
import "../EcdsaDkgValidator.sol";

contract WalletRegistryStub is WalletRegistry {
    constructor(
        SortitionPool _sortitionPool,
        IStaking _staking,
        EcdsaDkgValidator _ecdsaDkgValidator,
        IRandomBeacon _randomBeacon,
        ReimbursementPool _reimbursementPool
    )
        WalletRegistry(
            _sortitionPool,
            _staking,
            _ecdsaDkgValidator,
            _randomBeacon,
            _reimbursementPool
        )
    {}

    function getDkgData() external view returns (EcdsaDkg.Data memory) {
        return dkg;
    }
}
