pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";


contract TestTokenGrantLookupVesting {

    // Create KEEP token.
    KeepToken t = new KeepToken();

    // Create token grant contract with 30 days withdrawal delay.
    TokenGrant c = new TokenGrant(t, 0, 30 days);

    uint id;
    address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;
    uint start = now;
    uint duration = 10 days;
    uint cliff = 0;

    function beforeAll() public {
        t.approve(address(c), 100);
        id = c.grant(100, beneficiary, duration, start, cliff, false);
    }

    function testCanGetGrantVestingScheduleByGrantID() public {
        address _owner;
        uint _duration;
        uint _start;
        uint _cliff;
        (_owner, _duration, _start, _cliff) = c.getGrantVestingSchedule(id);
        Assert.equal(_owner, address(this), "Grant should maintain a record of the creator.");
        Assert.equal(_duration, duration, "Grant should have vesting schedule duration.");
        Assert.equal(_start, start, "Grant should have start time.");
        Assert.equal(_cliff, start+cliff, "Grant should have vesting schedule cliff duration.");
    }
}
