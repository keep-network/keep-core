const { BN } = require('@openzeppelin/test-helpers');
const ethWallet = require('ethereumjs-wallet').default;

function toBN (x) {
    return new BN(x);
}

function generateSalt () {
    return ethWallet.generate().getPrivateKeyString().substr(0, 34);
}

module.exports = {
    toBN,
    generateSalt,
};
