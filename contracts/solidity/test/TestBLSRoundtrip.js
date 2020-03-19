import {bls} from './helpers/data'

const BLS = artifacts.require('./cryptography/BLS.sol');
const AltBn128 = artifacts.require('./cryptography/AltBn128.sol');
const AltBn128Stub = artifacts.require('./stubs/AltBn128Stub.sol');

contract('TestBLSRoundtrip', function(accounts) {
  let blsLibrary, altBn128Library, altBn128StubLibrary, from;
  
  before(async () => {
    from = accounts[1]
    blsLibrary = await BLS.new();
    altBn128Library = await AltBn128.new();
    
    await AltBn128Stub.link("AltBn128", altBn128Library.address);
    
    altBn128StubLibrary = await AltBn128Stub.new();
  });

  it("should be able to sign a message and verify it", async function() {
    let signature = await altBn128StubLibrary.sign(bls.secretKey, {from: from});

    let message = await altBn128StubLibrary.g1HashToPoint(from);

    let actual = await blsLibrary.verify(bls.groupPubKey, message, signature)
    assert.isTrue(actual, "Should be able to verify valid BLS signature.");
  })

});
