# data CandidateGroup = CandidateGroup
#   { members   :: [Staker]
#   , createdAt :: Int    -- blockheight
#   , minH      :: Int    -- minimum number of honest parties required
#   , maxM      :: Int }  -- maximum number of misbehavers before failure
CandidateGroup = namedtuple(
    'CandidateGroup',
    ['members', 'createdAt', 'minH', 'maxM']
)


# data Staker = Staker
#   { ecdsaPubkey  :: EcdsaPubkey
#   , blsPubkey    :: BlsPubkey
#   , stakedTokens :: TokenAmount }
Staker = namedtuple(
    'Staker',
    ['ecdsaPubkey', 'blsPubkey', 'stakedTokens']
)


# data Result = Result
#   { groupPubkey  :: Maybe BeaconPubkey
#   , disqualified :: [Bool]
#   , inactive     :: [Bool] }
Result = namedtuple(
    'Result',
    ['groupPubkey', 'disqualified', 'inactive']
)


# data SignedResult = SignedResult
#   { result    :: Result
#   , signature :: BlsSignature }
SignedResult = namedtuple(
    'SignedResult',
    ['result', 'signature']
)


# data ResultSubmission = ResultSubmission
#   { result    :: Result
#   , signature :: BlsSignature
#   , signers   :: [Bool] }
ResultSubmission = namedtuple(
    'ResultSubmission',
    ['result', 'signature', 'signers']
)


def Failure(disqualified = [False] * groupN):
    return Result(False, disqualified, [False] * groupN)


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

# type BlsPrivkey   = Bn256.Scalar
# type BlsPubkey    = Bn256.G2
# type BlsSignature = Bn256.G1


# blsPubkey :: Bn256.Scalar -> Bn256.G2
def blsPubkey(privkey):
    return Bn256.ecmul(Bn256.generator, privkey)


# blsSign :: Bytes -> Bn256.Scalar -> Bn256.G1
def blsSign(message, privkey):
    # h :: G1
    h = Bn256.asPoint(sha3(message))
    return Bn256.ecmul(h, privkey)


# blsVerify :: Bn256.G1 -> Bytes -> Bn256.G2 -> Bool
def blsVerify(signature, message, pubkey):
    h = Bn256.asPoint(sha3(message))
    sigPairing = Bn256.pairing(signature, Bn256.generator)
    keyPairing = Bn256.pairing(h, pubkey)
    return sigPairing == keyPairing


# blsCombinePubkeys :: [BlsPubkey] -> BlsPubkey
def blsCombinePubkeys(pubkeys):
    return reduce(Bn256.ecadd, pubkeys)


# getResult :: Dict Staker SignedResult
#           -> CandidateGroup
#           -> Maybe ResultSubmission
def getResult(signedResults, groupInfo):

    # preferredResults :: Dict Result [(Staker, BlsSignature)]
    # -- for each result, list its supporters and their signatures
    preferredResults = {}
    for (staker, signedResult) in signedResults:
        result = signedResult.result
        signature = signedResult.signature

        if result in preferredResults:
            preferredResults[result].append((staker, signature))
        else:
            preferredResults[result] = [(staker, signature)]

    # orderedResults :: Dict Result Int
    # -- for each result, counts how many members support it
    orderedResults = {}
    for (result, supporters) in preferredResults:
        orderedResults[result] = len(supporters)

    leadingResult = keyByMaxValue(orderedResults)

    participantPlurality = preferredResults[leadingResult]

    if len(participantPlurality) < groupInfo.minH:
        return False
    else:
        participants = groupInfo.members

        # submissionSigners :: [Bool]
        # -- set up the bit array showing who signed the result
        submissionSigners = [False] * len(participants)

        submissionSignatures = []

        # for each supporter of the plurality resul, set their index in the
        # signers array to True and add their signature to the group signature
        for (signer, signature) in participantPlurality:
            submissionSigners[participants.index(signer)] = True
            submissionSignatures.append(signature)

        # groupSignature :: BlsSignature
        groupSignature = reduce(Bn256.ecadd, submissionSignatures)

        return ResultSubmission(
            leadingResult,
            groupSignature,
            submissionSigners
        )


# validate :: ResultSubmission -> CandidateGroup -> Bool
def validate(submission, candidateGroup):

    # signature :: BlsSignature
    # a group signature for the result
    signature = submission.signature

    # signers :: [Bool]
    # A one-indexed array whose length is N,
    # signers[i] = does the group signature include
    # the signature of the i-th member
    signers = submission.signers

    # result :: Result
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


# isEligible :: Staker -> BlockHeight -> [Staker] -> Bool
def isEligible(sender, t_diff, participants):
    i = participants.index(sender)
    maxEligible = 1 + max(0, (t_diff - t_dkg) / t_step)

    return i <= maxEligible



# receiveResult :: CandidateGroup -> IO Result
def receiveResult(candidateGroup):
    finalized = False
    validResult = False

    t_init = candidateGroup.createdAt

    while not finalized:
        t_now = currentBlockHeight()
        t_diff = t_now - t_init

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
