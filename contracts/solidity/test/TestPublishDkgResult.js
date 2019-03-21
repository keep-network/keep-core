const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

contract('TestPublishDkgResult', function (accounts) {
  let keepGroupImplV1;

  beforeEach(async () => {
    keepGroupImplV1 = await KeepGroupImplV1.new();
  })

  it("should know that result has not been published yet", async function () {
  });

  it("should publish DKG result", async function () {
  });

  it("should now that result has been already published", async function () {
  });
})
