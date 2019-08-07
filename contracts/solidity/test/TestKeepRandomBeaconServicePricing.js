import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServicePricing', function(accounts) {

  let operatorContract, serviceContract, callbackContract;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
  });

  it("should successfully refund callback gas surplus to the requestor", async function() {

    // Set higher gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(200), 'gwei'));

    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](
      bls.seed,
      callbackContract.address,
      "callback(uint256)",
      {value: minimumPayment, from: accounts[1]}
    );

    let minimumCallbackPayment = await serviceContract.minimumCallbackPayment()
    let requestorBalance = await web3.eth.getBalance(accounts[1]);

    await operatorContract.relayEntry(bls.nextGroupSignature);

    // Put back the default gas price
    await serviceContract.setMinimumGasPrice(web3.utils.toWei(web3.utils.toBN(20), 'gwei'));

    let updatedMinimumCallbackPayment = await serviceContract.minimumCallbackPayment()
    let updatedRequestorBalance = await web3.eth.getBalance(accounts[1])

    let surplus = web3.utils.toBN(minimumCallbackPayment).sub(web3.utils.toBN(updatedMinimumCallbackPayment))
    let refund = web3.utils.toBN(updatedRequestorBalance).sub(web3.utils.toBN(requestorBalance))

    assert.isTrue(refund.eq(surplus), "Callback gas surplus should be refunded to the requestor.");

  });

});
