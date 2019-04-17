const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const KeepRandomBeaconProxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');
const KeepGroupProxy = artifacts.require('./KeepGroup.sol');
const KeepGroupImplV1 = artifacts.require('./KeepGroupImplV1.sol');

async function mineOneBlock() {
  await web3.currentProvider.send({
    jsonrpc: '2.0',
    method: 'evm_mine',
    id: new Date().getTime()
  }, function(err, _) {
    if (err) console.log("Error mining a block.", err);
  });
}

async function mineBlocks(blocks) {
  for (let i = 0; i < blocks; i++)
    await mineOneBlock();
}

contract('TestKeepGroupExpiration', function(accounts) {

  let stakingProxy, minimumStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge,
    groupExpirationTimeout, numberOfActiveGroups, rounds, totalGas, groupsBatch,
    keepRandomBeaconImplV1, keepRandomBeaconProxy,
    keepGroupImplV1, keepGroupProxy, keepGroupImplViaProxy

  beforeEach(async () => {
    // Initialize staking contract under proxy
    stakingProxy = await StakingProxy.new();
    
    keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.new();
    keepRandomBeaconProxy = await KeepRandomBeaconProxy.new(keepRandomBeaconImplV1.address);

    // Initialize Keep Group contract
    minimumStake = 200000;
    groupThreshold = 15;
    groupSize = 20;
    timeoutInitial = 20;
    timeoutSubmission = 50;
    timeoutChallenge = 60;

    groupExpirationTimeout = 300; // time in blocks after which group expires
    numberOfActiveGroups = 5;     // number of minimal number of active groups
    
    totalGas = 0;       // set total gas to zero before each test

    groupsBatch = 10;   // number of groups added in one round/batch
    rounds = 100;       // number of rounds of groupAdd and selectGroup batches

    keepGroupImplV1 = await KeepGroupImplV1.new();
    keepGroupProxy = await KeepGroupProxy.new(keepGroupImplV1.address);
    keepGroupImplViaProxy = await KeepGroupImplV1.at(keepGroupProxy.address);
    await keepGroupImplViaProxy.initialize(
      stakingProxy.address, keepRandomBeaconProxy.address, minimumStake, groupThreshold, groupSize, timeoutInitial, timeoutSubmission, timeoutChallenge, groupExpirationTimeout, numberOfActiveGroups
    );
  });

  /* 
   * Function checking the gas price of running selectGroup after adding
   * groupsBatch of groups at once. selectGroup is staticaly initialized.
   * Without block mining.
   * 
   * Not a realistic scenario.
   */
  async function checkGasPerBatch(select, rounds, groupsBatch) {
    let after = 0;
    let gsGas = 0;
    let gsTotalGas = 0;
    let gstx = 0;
  
    for (var j = 1; j <= rounds; j++) {
      for (var i = 1; i <= groupsBatch; i++) {
        let gtx = await keepGroupImplViaProxy.groupAdd([0]);
        totalGas += gtx.receipt.gasUsed;
      }
      switch(select) {
        case 0: gstx = await keepGroupImplViaProxy.selectGroupV0("1"); break;
        case 1: gstx = await keepGroupImplViaProxy.selectGroupV1("1"); break;
        case 2: gstx = await keepGroupImplViaProxy.selectGroupV2("1"); break;
        case 3: gstx = await keepGroupImplViaProxy.selectGroupV3("1", 100); break;
      }
      gsGas = gstx.receipt.gasUsed;
      gsTotalGas += gsGas;
      after += groupsBatch;
      switch(select) {
        case 0: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "groupAdd -> ", Number(totalGas), "selectGroupV0() ->", Number(gsGas), "total ->", Number(gsTotalGas)); break;
        case 1: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "groupAdd -> ", Number(totalGas), "selectGroupV1() ->", Number(gsGas), "total ->", Number(gsTotalGas)); break;
        case 2: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "groupAdd -> ", Number(totalGas), "selectGroupV2() ->", Number(gsGas), "total ->", Number(gsTotalGas)); break;
        case 3: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "groupAdd -> ", Number(totalGas), "selectGroupV3() ->", Number(gsGas), "total ->", Number(gsTotalGas)); break;
      }
    }
  }

  /*
   * Function checking the gas price of running selectGroup after adding
   * a single group. selectGroup is staticaly initialized. Without block mining.
   * 
   * Not a realistic scenario.
   */
  async function checkGasPerGroupAdd(select, rounds, groupsBatch) {
    let after = 0;
    let gsGas = 0;
    let gsTotalGas = 0;
    let gstx = 0;
  
    for (var j = 1; j <= rounds; j++) {
      let gsMinGas = 99999999999999;
      let gsMaxGas = 0;
      for (var i = 1; i <= groupsBatch; i++) {
        await keepGroupImplViaProxy.groupAdd([0]);
        switch(select) {
          case 0: gstx = await keepGroupImplViaProxy.selectGroupV0("1"); break;
          case 1: gstx = await keepGroupImplViaProxy.selectGroupV1("1"); break;
          case 2: gstx = await keepGroupImplViaProxy.selectGroupV2("1"); break;
          case 3: gstx = await keepGroupImplViaProxy.selectGroupV3("1", 100); break;
        }
        gsGas = gstx.receipt.gasUsed;
        if (gsGas > gsMaxGas)
          gsMaxGas = gsGas;
        if (gsGas < gsMinGas)
          gsMinGas = gsGas;
        gsTotalGas += gsGas;
      }
      after += groupsBatch;
      switch(select) {
        case 0: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV0() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 1: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV1() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 2: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV2() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 3: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV3() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
      }
    }
  }

  /*
   * Function checking the gas price of running selectGroup after adding
   * a single group. selectGroup is staticaly initialized. With increasing block
   * mining.
   * 
   * Not a realistic scenario.
   */
  async function checkGasPerGroupAddWithMining(select, rounds, groupsBatch) {
    let after = 0;
    let gsGas = 0;
    let gsTotalGas = 0;
    let gstx = 0;
  
    for (var j = 1; j <= rounds; j++) {
      let gsMinGas = 99999999999999;
      let gsMaxGas = 0;
      for (var i = 1; i <= groupsBatch*j; i++) {
        await keepGroupImplViaProxy.groupAdd([0]);
        switch(select) {
          case 0: gstx = await keepGroupImplViaProxy.selectGroupV0("1"); break;
          case 1: gstx = await keepGroupImplViaProxy.selectGroupV1("1"); break;
          case 2: gstx = await keepGroupImplViaProxy.selectGroupV2("1"); break;
          case 3: gstx = await keepGroupImplViaProxy.selectGroupV3("1", 100); break;
        }
        gsGas = gstx.receipt.gasUsed;
        if (gsGas > gsMaxGas)
          gsMaxGas = gsGas;
        if (gsGas < gsMinGas)
          gsMinGas = gsGas;
        gsTotalGas += gsGas;
      }
      await mineBlocks(groupExpirationTimeout*j);
      after += groupsBatch;
      switch(select) {
        case 0: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV0() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 1: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV1() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 2: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV2() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 3: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV3() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
      }
    }
  }

  /*
   * Pseudorandom function used for introducing randomness in folowing tests.
   */
  function random(i) {
    return 99 % i;
  }

  /*
   * Function checking the gas price of running selectGroup after adding
   * a single group. selectGroup is pseudorandomly initialized. With increasing
   * block mining.
   * 
   * Not a realistic scenario.
   */
  async function checkGasPerGroupAddRandomlyWithMining(select, rounds, groupsBatch) {
    let after = 0;
    let gsGas = 0;
    let gsTotalGas = 0;
    let gstx = 0;
  
    for (var j = 1; j <= rounds; j++) {
      let gsMinGas = 99999999999999;
      let gsMaxGas = 0;
      for (var i = 1; i <= groupsBatch*j; i++) {
        await keepGroupImplViaProxy.groupAdd([0]);
        switch(select) {
          case 0: gstx = await keepGroupImplViaProxy.selectGroupV0(random(i)); break;
          case 1: gstx = await keepGroupImplViaProxy.selectGroupV1(random(i)); break;
          case 2: gstx = await keepGroupImplViaProxy.selectGroupV2(random(i)); break;
          case 3: gstx = await keepGroupImplViaProxy.selectGroupV3(random(i), groupsBatch); break;
        }
        gsGas = gstx.receipt.gasUsed;
        if (gsGas > gsMaxGas)
          gsMaxGas = gsGas;
        if (gsGas < gsMinGas)
          gsMinGas = gsGas;
        gsTotalGas += gsGas;
      }
      await mineBlocks(groupExpirationTimeout*j);
      after += groupsBatch;
      switch(select) {
        case 0: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV0() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 1: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV1() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 2: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV2() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 3: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV3() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
      }
    }
  }

  /*
   * Function checking the gas price of running selectGroup after adding
   * a single group. selectGroup is pseudorandomly initialized. With
   * pseudorandom block mining.
   * 
   * A realistic scenario.
   */
  async function checkGasPerGroupAddRandomlyWithRandomMining(select, rounds, groupsBatch) {
    let after = 0;
    let gsGas = 0;
    let gsTotalGas = 0;
    let gstx = 0;
  
    for (var j = 1; j <= rounds; j++) {
      let gsMinGas = 99999999999999;
      let gsMaxGas = 0;
      for (var i = 1; i <= groupsBatch*j; i++) {
        await keepGroupImplViaProxy.groupAdd([0]);
        switch(select) {
          case 0: gstx = await keepGroupImplViaProxy.selectGroupV0(random(i)); break;
          case 1: gstx = await keepGroupImplViaProxy.selectGroupV1(random(i)); break;
          case 2: gstx = await keepGroupImplViaProxy.selectGroupV2(random(i)); break;
          case 3: gstx = await keepGroupImplViaProxy.selectGroupV3(random(i), groupsBatch); break;
        }
        gsGas = gstx.receipt.gasUsed;
        if (gsGas > gsMaxGas)
          gsMaxGas = gsGas;
        if (gsGas < gsMinGas)
          gsMinGas = gsGas;
        gsTotalGas += gsGas;
        await mineBlocks(35*random(j));
      }
      after += groupsBatch;
      switch(select) {
        case 0: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV0() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 1: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV1() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 2: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV2() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
        case 3: console.log(Number(after), Number(await keepGroupImplViaProxy.numberOfGroups()), "selectGroupV3() ->", Number(gsGas), "total ->", Number(gsTotalGas), "MinMax", gsMinGas, gsMaxGas); break;
      }
    }
  }

  /*
   * Per batch gas checking section.
   */
   it("Checking gas price for selectGroupV0 per batch", async function() {
    this.timeout(0);
    await checkGasPerBatch(0, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV1 per batch", async function() {
    this.timeout(0);
    await checkGasPerBatch(1, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV2 per batch", async function() {
    this.timeout(0);
    await checkGasPerBatch(2, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV3 per batch", async function() {
    this.timeout(0);
    await checkGasPerBatch(3, rounds, groupsBatch);
  });

  /*
   * Per groupAdd gas checking section without mining.
   */
  it("Checking gas price for selectGroupV0 without mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAdd(0, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV1 without mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAdd(1, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV2 without mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAdd(2, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV3 without mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAdd(3, rounds, groupsBatch);
  });

  /*
   * Per groupAdd gas checking section with increasing mining.
   */
  it("Checking gas price for selectGroupV0 with increasing mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddWithMining(0, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV1 with increasing mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddWithMining(1, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV2 with increasing mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddWithMining(2, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV3 with increasing mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddWithMining(3, rounds, groupsBatch);
  });

  /*
   * Per groupAdd random gas checking section with increasing mining.
   */
  it("Checking gas price for random selectGroupV0 with mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithMining(0, rounds, groupsBatch);
  });

  it("Checking gas price for random selectGroupV1 with mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithMining(1, rounds, groupsBatch);
  });

  it("Checking gas price for random selectGroupV2 with mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithMining(2, rounds, groupsBatch);
  });

  it("Checking gas price for random selectGroupV3 with mining", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithMining(3, rounds, groupsBatch);
  });

  /*
   * Per groupAdd random gas checking section with random mining.
   */
  it("Checking gas price for selectGroupV0 with random mineblocks and random", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithRandomMining(0, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV1 with random mineblocks and random", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithRandomMining(1, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV2 with random mineblocks and random", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithRandomMining(2, rounds, groupsBatch);
  });

  it("Checking gas price for selectGroupV3 with random mineblocks and random", async function() {
    this.timeout(0);
    await checkGasPerGroupAddRandomlyWithRandomMining(3, rounds, groupsBatch);
  });
});