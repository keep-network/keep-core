var assert = require('chai').assert

module.exports = async (promise, message) => {
  try {
    await promise;
  } catch (error) {
    assert.include(error.message, message);

     return;
  }
  assert.fail('Expected throw not received');
}; 