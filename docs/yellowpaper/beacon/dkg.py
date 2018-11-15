# i = always the player whose perspective we're in

self.sharePolyCoeffs = [0..T].map(randomScalar)
self.blindingFactors = [0..T].map(randomScalar)


def f_i(z):
    return evaluateAt(z, self.sharePolyCoeffs)


def g_i(z):
    return evaluateAt(z, self.blindingFactors)


z_i = self.sharePolyCoeffs[0]
# assert(z_i == f_i(0))


self.commitments = map(ecCommit, self.sharePolyCoeffs, self.blindingFactors)


encryptedShares = []
for j in goodParticipants[3]:
    s_ij = f_i(j)
    t_ij = g_i(j)

    pointsBytes = marshalPoints(s_ij, t_ij)
    payload_ij = encrypt(symkey_ij, pointsBytes)

    encryptedShares[j] = payload_ij


#
# PHASE 4
#




shareComplaints = []

for j in goodParticipants[4]:
    privkey_ij = self.ephemeralKey[j]
    pubkey_ji = ephemeralPubkey(j, i)

    symkey_ij = ecdh(privkey_ij, pubkey_ji)

    shares = decryptAndValidateShares(
        senderIndex = j,
        recipientIndex = i,
        symkey = symkey_ij
     )

    if not validShares:
        shareComplaints.append(shareComplaint(j))
    else:
        (s_ji, t_ji) = validShares
        self.shares[j] = (s_ji, t_ji)

broadcast(shareComplaints)

#
# PHASE 5
#

messages.receive(4)

for complaint in messages[4]:
    j = c.senderIndex
    m = c.accusedIndex
    privkey_jm = c.privkey

    if not validatePrivkey(
        senderIndex = j,
        recipientIndex = m,
        privkey = privkey_jm
    ):
        disqualify(5, j)
    else:
        pubkey_mj = ephemeralPubkey(m, j)

        symkey_jm = ecdh(privkey_jm, pubkey_mj)

        sharesValid = decryptAndValidateShares(
            senderIndex = m,
            recipientIndex = j,
            symkey = symkey_jm
        )

        if not sharesValid:
            disqualify(5, m)
        else:
            disqualify(5, j)


#
# PHASE 6
#

# q = curveOrder

self.x = sum([
    self.shares[j].share_S for j in goodParticipants[6]
]) % BLS_CURVE_ORDER

self.x_prime = sum([
    self.shares[j].share_T for j in goodParticipants[6]
]) % BLS_CURVE_ORDER

#
# PHASE 7
#

self.pubkeyCoeffs = [
    G.ecMul(A_ik) for A_ik in self.sharePolyCoeffs
]


#
# PHASE 8
#

pubkeyComplaints = []

for j in goodParticipants[8]:
    pubkeyShareValid = validatePubkeyCoeffs(
        senderIndex = j,
        recipientIndex = i,
        share_S = self.shares[j].share_S
    )

    if not pubkeyShareValid:
        pubkeyComplaints.append(pubkeyComplaint(j))

#
# PHASE 9
#

for complaint in messages[8]:
    j = complaint.senderIndex
    m = complaint.accusedIndex
    privkey_jm = complaint.privkey

    if not validatePrivkey(
            senderIndex = j,
            recipientIndex = m,
            privkey = privkey_jm
    ):
        disqualify(9, j)
    else:
        pubkey_mj = ephemeralPubkeys[m][j]

        symkey = ecdh(privkey_jm, pubkey_mj)

        badActor = resolvePubkeyComplaint(
            senderIndex = m,
            recipientIndex = j,
            symkey = symkey
        )

        if badActor == "accused" or badActor == "both":
            disqualify(9, m)
        if badActor == "complainer" or badActor == "both":
            disqualify(9, j)

#
# PHASE 10
#

disqualifiedKeys = []

for m in disqualifiedInPhase[9]:
    keyPackage = (m, self.ephemeralKey[m])
    disqualifiedKeys.append(keyPackage)

broadcast(disqualifiedKeys)


#
# PHASE 11
#

for keys_j in messages[10]:
    j = keys_j.sender
    for keyPackage in keys_j.keyPackages:
        m = keyPackage.index
        privkey_jm = keyPackage.ephemeralKey

        if not disqualifiedInPhase[9].contains(m):
            # P_j broadcast the wrong keys
            disqualify(11, j)

        if not validatePrivkey(
            senderIndex = j,
            recipientIndex = m,
            privkey = privkey_jm
        ):
            # P_j broadcast invalid keys
            disqualify(11, j)
        else:
            pubkey_mj = ephemeralPubkey(m, j)
            symkey_jm = ecdh(privkey_jm, pubkey_mj)

            validShares = decryptAndValidateShares(
                senderIndex = m,
                recipientIndex = j,
                symkey = symkey_jm
            )

            if not validShares:
                # P_j failed to complain earlier
                disqualify(11, j)
            else:
                (s_mj, t_mj) = validShares
                self.knownShares[m][j] = (s_mj, t_mj)


# TODO

#
# PHASE 12
#

Y = ecSum([
    pubkeyCoeffs(j)[0] for j in goodParticipants[6]
])
