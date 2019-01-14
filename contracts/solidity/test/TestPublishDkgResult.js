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
  
  it("should know if member is inactive of disqualified", async function() {
    let DQbytes = "0x000000000000000101"
    let IAbytes = "0x000000000000010000"
    let groupPubKey = "0x100101011";
    await keepGroupImplV1.submitDkgResult(1814, true, groupPubKey , DQbytes, IAbytes);
    let DQ = [false, false, false, false, false, false, false, true, true];
    let IA = [false, false, false, false, false, false, true, false, false];
    let resDQ = [];
    let resIA = [];
    for(let i=0; i< (DQbytes.length -2)/2; i++){
        resDQ.push(await keepGroupImplV1.isDisqualified(groupPubKey, i));
        resIA.push(await keepGroupImplV1.isInactive(groupPubKey, i));
        }
    
    assert.equal(JSON.stringify(DQ), JSON.stringify(resDQ), "did not correctly return disqualified members");
    assert.equal(JSON.stringify(IA), JSON.stringify(resIA), "did not correctly return inactive members");
 
});
})



