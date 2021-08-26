const { BN } = require('@openzeppelin/test-helpers');

function toBN (x) {
    return new BN(x);
}

module.exports = {
    toBN,
};
