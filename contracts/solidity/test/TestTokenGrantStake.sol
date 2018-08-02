pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";


contract TestTokenGrantStake {

    // Create KEEP token
    KeepToken t = new KeepToken();

    // Create token grant contract with 30 days withdrawal delay.
    TokenGrant c = new TokenGrant(t, 0, 30 days);

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

        bool _locked;
        (, , _locked, , , , , , ,) = c.grants(id);
        Assert.equal(_locked, true, "Token grant should become locked.");
    }

    // Token grant beneficiary should be able to initiate unstake of the token grant
    function testCanInitiateUnstakeTokenGrant() public {
        c.initiateUnstake(id);
        Assert.equal(c.stakeWithdrawalStart(id), now, "Stake withdrawal start should be set.");
        Assert.equal(c.stakeBalances(beneficiary), 0, "Stake balance should change immediately after unstake initiation.");
    }
}
