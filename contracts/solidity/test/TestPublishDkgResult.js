const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

contract('TestPublishDkgResult', function(accounts) {
  let keepGroupImplV1;

  beforeEach(async () => {
      keepGroupImplV1 = await KeepGroupImplV1.new();
  })

  it("should know that result has not been published yet", async function() {
      let published = await keepGroupImplV1.isDkgResultSubmitted(1811);
    
      assert.equal(published, false, "result has not been published yet")
  });

  it("should publish DKG result", async function() {
      let event = keepGroupImplV1.DkgResultPublishedEvent()

      await keepGroupImplV1.submitDkgResult(1812, true, "0x100101011", "0x000000011", "0x000000100");

      event.get(function(error, result) {
          assert.equal(result[0].event, "DkgResultPublishedEvent", "DkgResultPublishedEvent should occur")
      });     
  }); 

  it("should now that result has been already published", async function() {
      await keepGroupImplV1.submitDkgResult(1813, true, "0x100101011", "0x000000011", "0x000000100");

      let published = await keepGroupImplV1.isDkgResultSubmitted(1813);
    
      assert.equal(published, true, "result has been already published")
  });
})



