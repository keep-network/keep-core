const BLS = artifacts.require('./cryptography/BLS.sol');

contract('TestBLS', function() {

  let bls;
  beforeEach(async () => {
    bls = await BLS.new();
  });

  it("should be able to verify threshold BLS recovered/reconstructed signature", async function() {

    // Data generated using client Go code with master secret key 123 and message "Hello!"
    let result = await bls.verify(
      "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0",
      "0x48656c6c6f21",
      "0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85"
    )
    assert.equal(result, true, "Should be able to verify valid BLS signature.");
  });

  it("should be able to verify BLS aggregated signature", async function() {

    // Data generated using client Go code with multiple random signers signing the same message "Hello!"
    let result = await bls.verify(
      "0x05c188c72f44373a42008c55499dfb4eb8944d89a62b10cc395c6cf11acff7c71f96fdd34a73284dec3126d19db9dfa2fb07752709818ce58a5a539c62ff09ae",
      "0x48656c6c6f21",
      "0xafd0185522d03e015e2165ad450af72a3b601673e8b41bc7f07014aa80892b24"
    )
    assert.equal(result, true, "Should be able to verify valid BLS signature.");
  });

  it("should fail to verify non valid BLS signature", async function() {

    let result = await bls.verify(
      "0x05c188c72f44373a42008c55499dfb4eb8944d89a62b10cc395c6cf11acff7c71f96fdd34a73284dec3126d19db9dfa2fb07752709818ce58a5a539c62ff09ae",
      "0x48656c6c6f21",
      "0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85"
    )
    assert.equal(result, false, "Should return false for failed verification.");
  });

  it("should fail to verify BLS signature without valid message", async function() {

    let result = await bls.verify(
      "0x05c188c72f44373a42008c55499dfb4eb8944d89a62b10cc395c6cf11acff7c71f96fdd34a73284dec3126d19db9dfa2fb07752709818ce58a5a539c62ff09ae",
      "0x123456789",
      "0xafd0185522d03e015e2165ad450af72a3b601673e8b41bc7f07014aa80892b24"
    )
    assert.equal(result, false, "Should return false for failed verification.");
  });

  it("should fail to verify BLS signature without valid public key", async function() {

    let result = await bls.verify(
      "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0",
      "0x48656c6c6f21",
      "0xafd0185522d03e015e2165ad450af72a3b601673e8b41bc7f07014aa80892b24"
    )
    assert.equal(result, false, "Should return false for failed verification.");
  });
});
