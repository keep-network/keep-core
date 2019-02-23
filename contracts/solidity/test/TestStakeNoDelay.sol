pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenStaking.sol";

contract TestStakeNoDelay {
    // Create KEEP token
    KeepToken token = new KeepToken();

    // Create staking contract with no withdrawal delay
    TokenStaking stakingContract = new TokenStaking(address(token), address(0), 0);

    uint withdrawalId;

    // Token holder should be able to stake it's tokens
    function testCanStake() public {
        uint balance = token.balanceOf(address(this));

        token.approveAndCall(address(stakingContract), 100, "0x00");

        Assert.equal(token.balanceOf(address(this)), balance - 100, "Stake amount should be taken out from token holder's main balance.");
        Assert.equal(stakingContract.stakeBalanceOf(address(this)), 100, "Stake amount should be added to token holder's stake balance.");
    }


    // Token holder should be able to initiate unstake of it's tokens
    function testCanInitiateUnstake() public {
        uint balance = token.balanceOf(address(this));

        withdrawalId = stakingContract.initiateUnstake(100);

        // Inspect created withdrawal request
        address owner;
        uint amount;
        uint start;
        (owner, amount, start) = stakingContract.withdrawals(withdrawalId);
        Assert.equal(owner, address(this), "Withdrawal request should maintain a record of the owner.");
        Assert.equal(amount, 100, "Withdrawal request should maintain a record of the amount.");
        Assert.equal(start, now, "Withdrawal request should maintain a record of when it was initiated.");

        Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance."); 
        Assert.equal(token.balanceOf(address(this)), balance, "Unstake amount should not be added to token holder main balance.");
    }

    // Should be able to finish unstake of it's tokens when withdrawal delay is over
    function testCanFinishUnstake() public {
        uint balance = token.balanceOf(address(this));

        stakingContract.finishUnstake(withdrawalId);
        Assert.equal(token.balanceOf(address(this)), balance + 100, "Unstake amount should be added to token holder main balance.");
        Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Stake balance should be empty.");

        // Inspect changes in withdrawal request
        address owner;
        uint amount;
        uint start;
        (owner, amount, start) = stakingContract.withdrawals(withdrawalId);
        Assert.isZero(owner, "Withdrawal request should be cleared.");
        Assert.isZero(amount, "Withdrawal request should be cleared.");
        Assert.isZero(start, "Withdrawal request should be cleared.");
    }
}
