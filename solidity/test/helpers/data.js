const {web3} = require("@openzeppelin/test-environment")

const blsData = {
  // data generated using master secret key 123
  secretKey: 123,
  
  // altbn128 public key for secret key 123
  groupPubKey: "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",

  // initial beacon entry from KeepRandomBeaconServiceImplV1.sol
  previousEntry: "0x15c30f4b6cf6dbbcbdcc10fe22f54c8170aea44e198139b776d512d8f027319a1b9e8bfaf1383978231ce98e42bafc8129f473fc993cf60ce327f7d223460663",

  // group signature over previousEntry
  groupSignature: "0x112d462728e89432b0fe40251eeb6608aed4560f3dc833a9877f5010ace9b1312006dbbe2f30c6e0e3e7ec47dc078b7b6b773379d44d64e44ec4e017bfa7375c",
  
  // group signature over groupSignature
  nextGroupSignature: "0x144b0508140c2c63fac298ee0cdd98571971a2d958f7c97d2bab82f3e1e727542d08314d6f087aca6ec2173b9a1d928cb80ff45258984a8929977a58d8b2fc26",
  
  // group signature over nextGroupSignature 
  nextNextGroupSignature: "0x10bbc10ee3e5509fffa43f797a5967c46f94b19a01d40360f88d0b13a2a5dc491112a6540550709ae6673af3d4e90ef96c7d7593b71ba6e335276d0b5fc5f3ae",

  // group signature over nextNextGroupSignature 
  nextNextNextGroupSignature: "0x2ed17b6237a4b9b9389f6964c4a07017d8461cb2602c367fbce7f9414585f6fc0ba88e1264f783386fb8493d31d733a100bf21f46acf0b3bf99c89834517907f",

  // uint256(keccak256(groupSignature))
  groupSignatureNumber: web3.utils.toBN("33136845259729814081977577432759716433762925710284529101628804651946005705295"),

  // uint256(keccak256(nextGroupSignature))
  nextGroupSignatureNumber: web3.utils.toBN("63728602218731680462391130646366490115686215934918152495059999747427209618626"),

  // uint256(keccak256(nextNextGroupSignatureNumber))
  nextNextGroupSignatureNumber: web3.utils.toBN("9775714160560684692317854436734528040997044179054191762446020971841873971471"),

  // uint256(keccak256(nextNextNextGroupSignatureNumber))
  nextNextNextGroupSignatureNumber: web3.utils.toBN("76547413554296705539222469062348207443894285352667024772941555864069963020436"),
};

module.exports = blsData