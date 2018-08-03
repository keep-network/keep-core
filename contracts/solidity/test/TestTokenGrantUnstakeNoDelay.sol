pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";


contract TestTokenGrantUnstakeNoDelay {

    // Create KEEP token
    KeepToken t = new KeepToken();

    // Create token grant contract with no withdrawal delay
    TokenGrant c = new TokenGrant(t, 0, 0);

    uint id;
    address beneficiary = address(this); // For test simplicity set beneficiary the same as sender.
    uint start = now;
    uint duration = 10 days;
    uint cliff = 0;

    // Token grant beneficiary can finish unstake of token grant when delay is over
    function testCanFinishUnstakeTokenGrant() public {

        // Approve transfer of tokens to the token grant contract.
        t.approve(address(c), 100);
        // Create new token grant.
        id = c.grant(100, beneficiary, duration, start, cliff, false);
        // Stake token grant.
        c.stake(id);

        c.initiateUnstake(id);
        c.finishUnstake(id);
        Assert.equal(c.stakeBalances(beneficiary), 0, "Stake balance should stay unchanged.");
        bool _locked;
        (, , _locked, , , , , , ,) = c.grants(id);
        Assert.equal(_locked, false, "Grant should become unlocked.");
    }
}
