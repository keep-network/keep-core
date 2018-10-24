class Bid(NamedTuple):
    amount:         Money
    expiresAt:      Blockheight
    seedCommitment: Commitment[SeedValue]


class BidPool(NamedTuple):
    bidTotal: Money
    allBids:  OrderedList[Bid]


def popBid(
        pool: BidPool
) -> Option[Bid]:
    if empty(pool.allBids):
        return False
    else:
        topBid = pool.allBids.head
        remainingBids = pool.allBids.tail

        pool.bidTotal -= topBid.amount
        pool.allBids = remainingBids

        return topBid


def pushBid(
        pool: BidPool,
        bid: Bid
) -> BidPool:
    pool.bidTotal += bid.amount
    pool.allBids = addToOrderedList(pool.allBids, bid)

    return pool


def filterExpired(
        pool: BidPool,
        currentTime: Blockheight
) -> BidPool:
    filteredPool = BidPool(0, [])

    # for conceptual simplicity, we create a new pool containing
    # only unexpired bids, then overwrite the old pool with it
    for bid in pool.allBids:
        if bid.expiresAt > currentTime:
            pushBid(filteredPool, bid)

    pool = filteredPool
    return filteredPool


def tick(
        pool: BidPool,
        newBids: List[Bid]
) -> Option[BeaconOutput]:

    # add in new bids
    for bid in newBids:
        pushBid(pool, bid)

    # remove expired bids
    filterExpired(pool, getCurrentBlockHeight())

    # the heavy lifting function
    # the current price should be determined by how many outputs the beacon
    # has recently generated; a larger number of outputs means higher price
    currentPrice = getCurrentOutputPrice()

    if pool.bidTotal >= currentPrice:
        return generateOutput(pool, currentPrice)
    else:
        return False


def generateOutput(
        pool: BidPool,
        currentPrice: Money
) -> Option[BeaconOutput]:
    usedBids = BidPool(0, [])

    while usedBids.bidTotal < currentPrice:
        nextBid = popBid(pool)
        pushBid(usedBids, nextBid)

    # determine the total stake for this output generation,
    # and then individual stakers' stake
    #
    # higher bids -> higher stake
    #
    # note that the actual rewards for stakers will be proportionally higher
    # when bids and stakes are higher, due to costs being constant
    #
    # thus highest stakes create the best risk:reward ratio
    totalStake = usedBids.bidTotal * BID_STAKE_MULTIPLIER
    memberStake = totalStake / N

    # an alternative algorithm, creating more variability in the stakes
    #
    # getCurrentGasPrice() is a simplification and should actually use
    # a smoother estimate to correct for price fluctuations
    totalProfit = usedBids.bidTotal - (getCurrentGasPrice() * OUTPUT_GAS_COST)
    totalStake = totalProfit * BID_STAKE_MULTIPLIER
    memberStake = totalStake / N

    # select the group to perform the beacon output
    # based on the hash of the block containing the transaction
    # triggering the output generation
    generationGroup = selectGroup(getLatestBlockHash())

    # use some defined function to determine the time we will
    # wait for any single bidder to reveal their seed
    requestTimeout = timeoutByBidN(len(usedBids.allBids))

    startingTime = getCurrentBlockHeight()
    timeoutBlock = startingTime + requestTimeout

    # initially, the top bid is eligible to reveal seed
    eligibleBid = popBid(usedBids)

    revealedSeed = False

    while revealedSeed == False:
        # reached timeout -> top bidder no longer eligible
        # next highest bidder becomes eligible
        if getCurrentBlockHeight() >= timeoutBlock:
            timeoutBlock += requestTimeout
            eligibleBid = popBid(usedBids)

            # nobody responds within timeout, no output produced
            if not eligibleBid:
                return False

        seedValue = receiveSeedRevealTx()

        if checkCommitment(seedValue, eligibleBid.seedCommitment):
            revealedSeed = seedValue

    # we have a revealed seed, time to proceed with generation
    output = generationGroup.sign(revealedSeed, v_previous)
    if not empty(output.misbehavingMembers):
        penalize(output.misbehavingMembers, memberStake)

    return output



