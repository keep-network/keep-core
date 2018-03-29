// Specifically request an abstraction for KeepRelayBeacon
var KeepRelayBeacon = artifacts.require("KeepRelayBeacon");

contract('KeepRelayBeacon', function(accounts) {
	it("KeepRelayBeacon.requestRelay must return an event:T-00001: ", function() {
		return KeepRelayBeacon.deployed().then(function(instance) {
			return instance.requestRelay(12,12,12);
		}).then(function(events) {
			// console.log ( 'events=', events );
			assert.equal(events.logs.length, 1);
			assert.equal(events.logs[0].event, "RelayEntryRequested");
			//
			// TODO: should check some other stuff at this point too!
			//
		});
	});
});
