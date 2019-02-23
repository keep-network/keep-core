pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";

contract TestTokenGrantRevoke {  
    
    // Create KEEP token.
    KeepToken t = new KeepToken();

    // Create token grant contract with 30 days withdrawal delay.
    TokenGrant c = new TokenGrant(address(t), address(0), 30 days);

    uint id;
    address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;

    // Grant owner can revoke revocable token grant.
    function testCanFullyRevokeGrant() public {
        uint balance = t.balanceOf(address(this));
    
        // Create revocable token grant.
        t.approve(address(c), 100);
        id = c.grant(100, beneficiary, 10 days, now, 0, true);
        
        Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be taken out from grant creator main balance.");
        Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
        
        c.revoke(id);

        Assert.equal(t.balanceOf(address(this)), balance, "Amount should be returned to token grant owner.");
        Assert.equal(c.balanceOf(beneficiary), 0, "Amount should be removed from beneficiary's grant balance.");
    }

    // Token grant creator can revoke the grant but no amount 
    // is refunded since duration of the vesting is over.
    function testCanZeroRevokeGrant() public {
        uint balance = t.balanceOf(address(this));
    
        // Create revocable token grant with 0 duration.
        t.approve(address(c), 100);
        id = c.grant(100, beneficiary, 0, now, 0, true);
        
        Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be removed from grant creator main balance.");
        Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
        
        c.revoke(id);

        Assert.equal(t.balanceOf(address(this)), balance - 100, "No amount to be returned to grant creator since vesting duration is over.");
        Assert.equal(c.balanceOf(beneficiary), 100, "Amount should stay at beneficiary's grant balance.");
    }

}
