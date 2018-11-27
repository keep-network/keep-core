BitArray = List[bool]
StakerSignature = Tuple[Staker, BlsSignature]

BlockHeight = NewType('BlockHeight', int)

class CandidateGroup(NamedTuple):
    members:   List[Staker]
    createdAt: Blockheight
    minH:      int # minimum number of honest parties required
    maxM:      int # maximum number of misbehaving parties before failure


class Staker(NamedTuple):
    ecdsaPubkey:  EcdsaPubkey
    blsPubkey:    BlsPubkey
    stakedTokens: TokenAmount


class Result(NamedTuple):
    groupPubkey:  Optional[BeaconPubkey]
    disqualified: Optional[BitArray]
    inactive:     Optional[BitArray]


class SignedResult(NamedTuple):
    result:    Result
    signature: BlsSignature


class ResultSubmission(NamedTuple):
    result:    Result
    signature: BlsSignature
    signers:   BitArray


def Failure(disqualified = None):
    return Result(False, disqualified, None)


def Success(pubkey, disqualified, inactive):
    return Result(pubkey, disqualified, inactive)


def isFailure(result):
    return result.pubkey == False


def isSuccess(result):
    return result.pubkey != False


# Bn256 interface:

# -- functions applying to each of G1 and G2:
# class CurvePoint gx where
#   generator :: gx
#   asPoint   :: Bytes -> gx
#   ecmul     :: gx -> Scalar -> gx
#   ecadd     :: gx -> gx -> gx
# instance CurvePoint G1
# instance CurvePoint G2

# -- calculate the pairing of two curve points
# pairing :: G1 -> G2 -> GT


BlsPrivkey   = Bn256.Scalar
BlsPubkey    = Bn256.G2
BlsSignature = Bn256.G1


def blsPubkey(
        privkey: BlsPrivkey
) -> BlsPubkey:
    return Bn256.ecmul(Bn256.generator, privkey)


def blsSign(
        message: Bytes,
        privkey: BlsPrivkey
) -> BlsSignature:
    h = Bn256.asPoint(sha3(message))
    return Bn256.ecmul(h, privkey)


def blsVerify(
        signature: BlsSignature,
        message:   Bytes,
        pubkey:    BlsPubkey
) -> bool:
    h = Bn256.asPoint(sha3(message))
    sigPairing = Bn256.pairing(signature, Bn256.generator)
    keyPairing = Bn256.pairing(h, pubkey)
    return sigPairing == keyPairing


def blsCombinePubkeys(
        pubkeys: List[BlsPubkey]
) -> BlsPubkey:
    return reduce(Bn256.ecadd, pubkeys)


def getResult(
        signedResults: Dict[Staker, SignedResult],
        groupInfo:     CandidateGroup
) -> Optional[ResultSubmission]:

    # for each result, list its supporters and their signatures
    preferredResults: Dict[Result, StakerSignature] = {}

    for (staker, signedResult) in signedResults:
        result = signedResult.result
        signature = signedResult.signature

        if result in preferredResults:
            preferredResults[result].append((staker, signature))
        else:
            preferredResults[result] = [(staker, signature)]

    # for each result, counts how many members support it
    orderedResults: Dict[Result int] = {}

    for (result, supporters) in preferredResults:
        orderedResults[result] = len(supporters)

    leadingResult = keyByMaxValue(orderedResults)

    participantPlurality = preferredResults[leadingResult]

    if len(participantPlurality) < groupInfo.minH:
        return None
    else:
        participants = groupInfo.members

        # set up the bit array showing who signed the result
        submissionSigners = [False] * len(participants)

        submissionSignatures = []

        # for each supporter of the plurality resul, set their index in the
        # signers array to True and add their signature to the group signature
        for (signer, signature) in participantPlurality:
            submissionSigners[participants.index(signer)] = True
            submissionSignatures.append(signature)

        groupSignature = reduce(Bn256.ecadd, submissionSignatures)

        return ResultSubmission(
            leadingResult,
            groupSignature,
            submissionSigners
        )


def validate(
        submission:     ResultSubmission,
        candidateGroup: CandidateGroup
) -> bool:

    # a group signature for the result
    signature = submission.signature

    # A one-indexed array whose length is N,
    # signers[i] = does the group signature include
    # the signature of the i-th member
    signers = submission.signers

    result = submission.result

    # require at least the honest majority number of signers
    signersRequired = candidateGroup.minH
    participants = candidateGroup.members

    # too few signatures?
    if signers.count(True) < signersRequired:
        return False
    else:
        # roll together the pubkeys of those who have signed the result
        signerPubkeys = []
        for s in signers, p in participants:
            if s == True:
                signerPubkeys.append(p.blsPubkey)
        creationPubkey = reduce(Bn256.ecadd, signerPubkeys)

        return blsVerify(
            signature,
            result,
            creationPubkey
        )


def isEligible(
        sender:       Staker,
        t_diff:       BlockHeight,
        participants: List[Staker]
) -> bool:
    i = participants.index(sender)
    maxEligible = 1 + max(0, (t_diff - t_dkg) / t_step)

    return i <= maxEligible


def receiveResult(
        candidateGroup: CandidateGroup
) -> Result:
    finalized = False
    validResult = False

    t_init = candidateGroup.createdAt

    while not finalized:
        t_now = currentBlockHeight()
        t_diff = BlockHeight(t_now - t_init)

        # has the timeout been reached?
        if t_diff > t_timeout:
            finalized = True
        else:
            # get latest submission
            (submission, sender) = latestSubmission()

            eligible = isEligible(sender, t_diff, candidateGroup.members)

            valid = validate(submission, candidateGroup)

            if eligible and valid:
                validResult = submission.result

                # rewardSubmitter :: Staker -> CandidateGroup -> IO ()
                # -- give the submitter a reward and penalize everyone whose
                # -- index is smaller than the submitter's for being late
                rewardSubmitter(sender, candidateGroup)

                finalized = True

    # if no valid result received -> no fault failure
    if not validResult:
        return Failure()
    else:
        return validResult
