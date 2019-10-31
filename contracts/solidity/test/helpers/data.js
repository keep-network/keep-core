export const bls = {
  // data generated using master secret key 123
  
  // compressed altbn128 public key for secret key 123
  groupPubKey: "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0",
  
  previousEntry: web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862'),
  seed: web3.utils.toBN('27182818284590452353602874713526624977572470936999595749669676277240766303535'),

  // compressed group signature for combined previousEntry | seed
  groupSignature: web3.utils.toBN('10920102476789591414949377782104707130412218726336356788412941355500907533021'),
  
  // compressed group signature for combined groupSignature | seed
  nextGroupSignature: web3.utils.toBN('14584315380385966635411228245329978667313108710498038331780910721312248967343'),
  
  // compressed group signature for combined nextGroupSignature | seed 
  nextNextGroupSignature: web3.utils.toBN('72012275566908154022988289911738818616836097778504726833078217965357610052464'),
};
