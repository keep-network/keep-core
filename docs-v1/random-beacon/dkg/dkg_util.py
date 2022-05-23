# i is generic and doesn't relate to perspective


# tag::phase-2[]
# Fetch the correct ephemeral pubkey from messages broadcast in phase 1
#
# The format for the message of `P_j` in phase `P` is: `messages[P][j]`
#
def ephemeralPubkey(senderIndex, recipientIndex):
    return messages[1][senderIndex].pubkey[recipientIndex]
# end::phase-2[]


# tag::phase-3[]
# Evaluate a polynomial given by `coeffs` at point `z`
#
# `coeffs` is little-endian; `ax^2 + bx + c` is expressed as `[c, b, a]`
#
# `evaluateAt(2, [6, 3, 4]) = 6 + (3 * 2^1) + (4 * 2^2) = 28`
#
def evaluateAt(z, coeffs):
    return sum(
        [ coeffs[k] * z^k for k in [0..M] ]
    )


# Pedersen commitment to secret value `s` and blinding factor `t`
# `G = P1` is the standard generator of the elliptic curve
# `H = G*a` is a custom generator where `a` is unknown
#
# C(s, t) = G*s + H*t
#
def ecCommit(s, t):
    Gs = P1.scalarMult(s)
    Ht = H.scalarMult(t)
    return ecAdd(Gs, Ht)
# end::phase-3[]


# tag::phase-4[]
# Calculate the sum of a list of elliptic curve points
def ecSum(points):
    return reduce(ecAdd, points)


# Fetch the correct encrypted shares from messages broadcast in phase 3
def encryptedShares(senderIndex, recipientIndex):
    return messages[3][senderIndex].encryptedShares[recipientIndex]


# Fetch a specific participant's commitments from messages broadcast in phase 3
def commitments(senderIndex):
    return messages[3][senderIndex].commitments


# Fetch the correct shares and try to decrypt them with the provided key
def decryptShares(
    senderIndex,
    recipientIndex,
    symkey
):
    payload = encryptedShares(senderIndex, recipientIndex)

    return decrypt(payload, symkey)


# Fetch the shares and validate them
#
# Check that shares decrypt correctly and are consistent with the sender's
# published commitments
#
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
        (share_S, share_T) = unmarshalPoints(plaintext)

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


# Check that equation 2 from GJKR holds for `share_S, share_T`
#
# P_i is the player whose shares are validated
# P_j is the perspective player performing the validation
#
# GJKR 1.(b):
#
#   g^s_ij * h^t_ij == product([ C_ik ^ (j^k) for k in [0..T] ]) % p
#
def checkShareConsistency(
    recipientIndex,
    senderIndex,
    share_S,
    share_T
):
    i = senderIndex
    j = recipientIndex

    C_i = commitments(i)

    C_ecSum = ecSum(
        [ C_i[k].scalarMult(j^k) for k in [0..M] ]
    )

    sharesValid = ecCommit(share_S, share_T) == C_ecSum
    return sharesValid
# end::phase-4[]


# tag::phase-5[]
# Check that a revealed private key matches previously broadcast public key
def validatePrivkey(senderIndex, recipientIndex, privkey):
    expectedPubkey = ephemeralPubkey(senderIndex, recipientIndex)
    return derivePubkey(privkey) == expectedPubkey
# end::phase-5[]


# tag::phase-8[]
# Fetch the sender's public key coeffs `A_ik` from messages broadcast in phase 7
def pubkeyCoeffs(senderIndex):
    return messages[7][senderIndex].pubkeyCoeffs


# P_i is the player whose public key share is calculated
# P_j is the perspective player
#
def pubkeyShare(senderIndex, recipientIndex):
    i = senderIndex
    j = recipientIndex

    A_i = pubkeyCoeffs(i)

    pubkeyShare = ecSum(
        [ A_i[k].scalarMult(j^k) for k in [0..M] ]
    )
    return pubkeyShare


# Check that equation 3 holds for `share_S`
#
# GJKR 4.(b):
#
#   g^s_ij == product([ A_ik ^ (j^k) for k in [0..T] ]) % p
#
def validatePubkeyCoeffs(
        senderIndex,
        recipientIndex,
        share_S
):
    return P1.scalarMult(share_S) == pubkeyShare(senderIndex, recipientIndex)
# end::phase-8[]


# tag::phase-9[]
# Check which party is at fault when a complaint is presented in phase 8
#
# Decrypt the shares the accused sent to the complainer in phase 3 and check
# the validity of the accused's `A_ik` values
#
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
        (share_S, _) = unmarshalPoints(plaintext)

        pubkeyValid = validatePubkeyCoeffs(
            senderIndex,
            recipientIndex,
            share_S
        )

        if pubkeyValid:
            return "complainer"
        else:
            return "accused"
# end::phase-9[]


# tag::phase-11[]
def reconstruct(shares, indices):
    secret = sum(
        [ share_k * lagrange(k, indices) for share_k in shares, k in indices ]
    )
    return secret % q
# end::phase-11[]


# tag::phase-12[]
# Calculate the individual public key of a specific participant
#
# P_i is each qualified participant in turn
#
# GJKR (C1'):
#
#   g^x_j
#     = g^( sum([ s_ij for i in QUAL ]) ) % p
#     = product([ g^s_ij for i in QUAL ]) % p
#     = product([ product([ A_ik ^ (j^k) for k in [0..T] ]) for i in QUAL ]) % p
#
def individualPublicKey(memberIndex, QUAL):
    pubkeyShares = [ pubkeyShare(i, memberIndex) for i in QUAL ]
    return ecSum(pubkeyShares)
# end::phase-12[]
