pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "../WalletFactory.sol";
import "../DKGValidator.sol";
import "../libraries/DKG.sol";

contract WalletFactoryStub is WalletFactory {
    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        Wallet _masterWallet
    )
        WalletFactory(
            _sortitionPool,
            _tToken,
            _staking,
            _dkgValidator,
            _masterWallet
        )
    {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }
}
