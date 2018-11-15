# i is generic and doesn't relate to perspective

def ecCommit(s, t):
    Gs = G.scalarMult(s)
    Ht = H.scalarMult(t)
    return ecAdd(Gs, Ht)


def ecSum(points):
    return reduce(ecAdd, points)


def evaluateAt(z, coeffs):
    Q = BLS_CURVE_ORDER
    return sum([
        coeffs[k] * z^k for k in [0..T]
    ]) % Q


def ephemeralPubkey(senderIndex, recipientIndex):
    return messages[1][senderIndex].pubkey[recipientIndex]


def validatePrivkey(senderIndex, recipientIndex, privkey):
    expectedPubkey = ephemeralPubkey(senderIndex, recipientIndex)
    return correspondingPubkey(privkey) == expectedPubkey


def encryptedShares(senderIndex, recipientIndex):
    return messages[3][senderIndex].encryptedShares[recipientIndex]


def commitments(senderIndex):
    return messages[3][senderIndex].commitments


def decryptShares(
    senderIndex,
    recipientIndex,
    symkey
):
    payload = encryptedShares(senderIndex, recipientIndex)

    return decrypt(payload, symkey)


def decryptAndValidateShares(
    senderIndex,
    recipientIndex,
    symkey
):
    plaintext = decryptShares(
        senderIndex,
        recipientIndex,
        symkey
    )

    if not plaintext:
        return False
    else:
        (share_S, share_T) = parsePolyPoints(plaintext)

        sharesValid = checkShareConsistency(
            senderIndex,
            recipientIndex,
            share_S,
            share_T
        )

        if sharesValid:
            return (share_S, share_T)
        else:
            return False


def checkShareConsistency(
    recipientIndex,
    senderIndex,
    share_S,
    share_T
):
    i = senderIndex
    j = recipientIndex

    C_i = commitments(j)

    C_ecSum = ecSum([
        C_i[k].scalarMult(j^k) for k in [0..T]
    ])

    sharesValid = ecCommit(share_S, share_T) == C_ecSum
    return sharesValid


def pubkeyCoeffs(senderIndex):
    return messages[7][senderIndex].pubkeyCoeffs


def reconstructMasterPubkey(QUAL):
    def A_(i): return pubkeyCoeffs(i)

    Y = ecSum([
        A_(i)[0] for i in QUAL
    ])

    return Y


def pubkeyShare(senderIndex, recipientIndex):
    i = senderIndex
    j = recipientIndex

    A_i = pubkeyCoeffs(i)

    pubkeyShare = ecSum([
        A_i[k].scalarMult(j^k) for k in [0..T]
    ])
    return pubkeyShare


def individualPublicKey(memberIndex, QUAL):
    pubkeyShares = [ pubkeyShare(i, memberIndex) for i in QUAL ]
    return ecSum(pubkeyShares)


def validatePubkeyCoeffs(
        senderIndex,
        recipientIndex,
        share_S
):
    return G.scalarMult(share_S) == pubkeyShare(senderIndex, recipientIndex)


def resolvePubkeyComplaint(
        senderIndex,
        recipientIndex,
        symkey
):
    plaintext = decryptShares(
        senderIndex,
        recipientIndex,
        symkey
    )

    if not plaintext:
        # only happens if the complainer failed to complain earlier
        # and thus both violated protocol
        return "both"
    else:
        (share_S, _) = parsePolyPoints(plaintext)

        pubkeyValid = validatePubkeyCoeffs(
            senderIndex,
            recipientIndex,
            share_S
        )

        if pubkeyValid:
            return "complainer"
        else:
            return "accused"
