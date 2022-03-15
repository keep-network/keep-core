const ecdsaData = {
  group1: {
    // ecdsa private key
    privateKey:
      "0x937ffe93cfc943d1a8fc0cb8bad44a978090a4623da81eefdff5380d0a290b41",

    // ecdsa public key
    publicKey:
      "0x9a0544440cc47779235ccb76d669590c2cd20c7e431f97e17a1093faf03291c473e661a208a8a565ca1e384059bd2ff7ff6886df081ff1229250099d388c83df",
    publicKeyX:
      "0x9a0544440cc47779235ccb76d669590c2cd20c7e431f97e17a1093faf03291c4",
    publicKeyY:
      "0x73e661a208a8a565ca1e384059bd2ff7ff6886df081ff1229250099d388c83df",

    // digest to sign
    digest1:
      "0x8bacaa8f02ef807f2f61ae8e00a5bfa4528148e0ae73b2bd54b71b8abe61268e",

    // group signature over `digest`
    signature1: {
      r: "0xedc074a86380cc7e2e4702eaf1bec87843bc0eb7ebd490f5bdd7f02493149170",
      s: "0x3f5005a26eb6f065ea9faea543e5ddb657d13892db2656499a43dfebd6e12efc",
      v: 28,
    },
  },
  group2: {
    // ecdsa private key
    privateKey:
      "0x212bd2787dd6b0b9064f0341ff279aaa08eab74077f4b8c5ece576d256b01514",

    // ecdsa public key
    publicKey:
      "0xadba3062ac8cd30319b03c88637fb8c20c868905fad7c86faf3b28791212e8854bc55fecc3a40a1fbd5213de0993604e79a99e234d2cef0b33523e63383574c2",

    // digest to sign
    digest1: "",

    // group signature over `digest`
    signature1: "",
  },
}

export default ecdsaData
