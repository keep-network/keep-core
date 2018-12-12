pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantStakeNoDelay {  
    
    // Create KEEP token
    KeepToken t = new KeepToken();

    // Create token grant contract with no withdrawal delay
    TokenGrant c = new TokenGrant(address(t), address(0), 0);

    uint id;
    address beneficiary = address(this); // For test simplicity set beneficiary the same as sender.
    uint start = now;
    uint duration = 10 days;
    uint cliff = 0;


    // Token grant beneficiary should be able to stake unreleased granted balance.
    function testCanStakeTokenGrant() public {

        // Approve transfer of tokens to the token grant contract.
        t.approve(address(c), 100);
        // Create new token grant.
        id = c.grant(100, beneficiary, duration, start, cliff, false);
        // Stake token grant.
        c.stake(id);

        Assert.equal(c.stakeBalances(beneficiary), 100, "Token grant balance should be added to beneficiary grant stake balance.");

        bool _staked;
        (, , _staked, , , , , , ,) = c.grants(id);
        Assert.equal(_staked, true, "Token grant should become staked.");
    }

    // Token grant beneficiary should be able to initiate unstake of the token grant
    function testCanInitiateUnstakeTokenGrant() public {
        c.initiateUnstake(id);
        Assert.equal(c.stakeWithdrawalStart(id), now, "Stake withdrawal start should be set.");
        Assert.equal(c.stakeBalances(beneficiary), 0, "Stake balance should change immediately after unstake initiation.");
    }

    // Token grant beneficiary can finish unstake of token grant when delay is over
    function testCanFinishUnstakeTokenGrant() public {
        c.finishUnstake(id);
        Assert.equal(c.stakeBalances(beneficiary), 0, "Stake balance should stay unchanged.");
        bool _staked;
        (, , _staked, , , , , , ,) = c.grants(id);
        Assert.equal(_staked, false, "Grant should become unstaked.");
    }
}
