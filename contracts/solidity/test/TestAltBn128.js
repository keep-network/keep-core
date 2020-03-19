import {bls} from './helpers/data'

const BLS = artifacts.require('./cryptography/BLS.sol');
const AltBn128 = artifacts.require('./cryptography/AltBn128.sol');
const AltBn128Stub = artifacts.require('./stubs/AltBn128Stub.sol');

contract('AltBn128', function() {
  let blsLibrary, altBn128Library, altBn128;
  
  before(async () => {
    blsLibrary = await BLS.new();
    altBn128Library = await AltBn128.new();
    
    await AltBn128Stub.link("AltBn128", altBn128Library.address);
    
    altBn128 = await AltBn128Stub.new();
  });

  it("should be able to sign a message and verify it", async function() {
    let message = web3.utils.stringToHex("A bear walks into a bar 123...")
    let signature = await altBn128.sign(message, bls.secretKey);

    let actualMessage = await altBn128.g1HashToPoint(message);

    let actual = await blsLibrary.verify(bls.groupPubKey, actualMessage, signature)
    assert.isTrue(actual, "Should be able to verify valid BLS signature.");
  })

});
