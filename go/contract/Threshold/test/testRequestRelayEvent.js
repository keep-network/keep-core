// Specifically request an abstraction for KStart
var KStart = artifacts.require("KStart");

contract('KStart', function(accounts) {
	it("KStart.requestRelay must return an event:T-00001: ", function() {
		return KStart.deployed().then(function(instance) {
			return instance.requestRelay(12,12,12);
		}).then(function(events) {
			// console.log ( 'events=', events );
			assert.equal(events.logs.length, 1);
			assert.equal(events.logs[0].event, "RequestRelayEvent");
			//
			// TODO: should check some other stuff at this point too!
			//
		});
	});
});
