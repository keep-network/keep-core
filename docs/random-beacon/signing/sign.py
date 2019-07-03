class Request(NamedTuple):
    seedCommitment:  Commitment[SeedValue]
    requesterPubkey: EcdsaPubkey
    stakeMultiplier: Float
    placedAt:        Blockheight
    openTimeout:     Blockheight
    previousOutput:  BlsSignature

def openCommitment(signingId, seed_i):
    request_i = Requests[signingId] # abort if no such request
    senderPubkey = getSenderPubkey()
    T_open = getCurrentBlockheight()

    # ignore commitment openings that arrive too late
    if T_open > request_i.openTimeout:
        abort()

    # ignore commitment openings by parties other than original requester
    if not request_i.requesterPubkey == senderPubkey:
        abort()

    # ignore invalid commitment openings
    if not request_i.seedCommitment == sha3(seed_i, senderPubkey):
        abort()

    Block_kPlus1 = getBlockByHeight(request_i.placedAt + 1)
    rseed_i = Block_kPlus1.blockhash

    Group_i = select(AllGroups, rseed_i)
    v_iMinus1 = request_i.previousOutput
    beaconInput = sha3(seed_i, rseed_i, v_iMinus1)

    outputWaiting = OpenOutput(
        startedAt    = T_open,
        signingGroup = Group_i,
        signingInput = beaconInput
    )

    OutputInProgress[signingId] = outputWaiting


def receiveOutput(signingId, outputSignature):
    outputWaiting = OutputInProgress[signingId] # abort if not found
    request_i = Requests[signingId]

    pubkey_Group_i = outputWaiting.signingGroup.groupPubkey
    input_i = outputWaiting.signingInput

    submitter = getSenderPubkey()

    signatureValid = blsVerify(
        outputSignature,
        input_i,
        pubkey_Group_i
    )

    if signatureValid:
        T_output = getCurrentBlockHeight() - outputWaiting.startedAt

        rewardGroup(
            submitter,
            group_i,
            T_output,
            request_i.stakeMultiplier
        )

    else:
        punish(
            submitter,
            INVALID_SIGNATURE_PENALTY * request_i.stakeMultiplier
        )

