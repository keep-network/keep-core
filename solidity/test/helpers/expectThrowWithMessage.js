export default async (promise, message) => {
  try {
    await promise;
  } catch (error) {
    assert.include(error.message, message);

     return;
  }
  assert.fail('Expected throw not received');
}; 