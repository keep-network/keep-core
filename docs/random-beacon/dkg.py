# i = always the player whose perspective we're in

# Because G1 and G2 in alt_bn128 are cyclic groups of prime order, this number
# can also be used as the size of the secret sharing finite field
q = G1.curveOrder

#
# PHASE 1
#

# Receive the DKG parameters from on-chain
dkgSetup = getDkgSetup()

# Presented from the perspective of P_i
i = dkgSetup.members.index(self.pubkey)

# Keep track of other qualified participants
#
# A compliant implementation won't get disqualified so we can assume P_i is
# always good
#
goodParticipants[1] = [1..N].remove(i)

# Record the blockheight at the start of the DKG
#
# Used later for calculating timeouts
#
T_dkgInit = getCurrentBlockHeight()

ephemeralPubkeys = []

for j in goodParticipants[1]:
    x_ij = genEcdhKeypair()

    self.ephemeralKey[j] = x_ij

    y_ij = x_ij.pubkey

    ephemeralPubkeys[j] = y_ij

broadcast(messagePhase1(ephemeralPubkeys))


#
# PHASE 2
#

# Receive messages from phase 1:
# - ephemeral public keys of other participants
messages.receive(1)

# `goodParticipants[P]` denotes the qualified participants in phase `P`
for j in goodParticipants[2]:
    privkey_ij = self.ephemeralKey[j]
    pubkey_ji = ephemeralPubkey(j, i)

    k_ij = ecdh(privkey_ij, pubkey_ji)
    self.symkey[j] = k_ij


#
# PHASE 3
#

# a_ij = sharePolyCoeffs[j]
# b_ij = blindingFactors[j]
self.sharePolyCoeffs = [0..M].map(G1.randomScalar)
self.blindingFactors = [0..M].map(G1.randomScalar)


def f_i(z):
    return evaluateAt(z, self.sharePolyCoeffs) % q


def g_i(z):
    return evaluateAt(z, self.blindingFactors) % q


z_i = self.sharePolyCoeffs[0]
# assert(z_i == f_i(0))


self.commitments = map(ecCommit, self.sharePolyCoeffs, self.blindingFactors)

encryptedShares = []

for j in goodParticipants[3]:
    s_ij = f_i(j)
    t_ij = g_i(j)

    pointsBytes = marshalPoints(s_ij, t_ij)
    payload_ij = encrypt(self.symkey[j], pointsBytes)

    encryptedShares[j] = payload_ij

broadcast(messagePhase3(encryptedShares, self.commitments))


#
# PHASE 4
#

# Receive messages from phase 3:
# - commitments to the secret sharing polynomials
# - encrypted share payloads
messages.receive(3)

shareComplaints = []

for j in goodParticipants[4]:
    k_ij = self.symkey[j]

    validShares = decryptAndValidateShares(
        senderIndex = j,
        recipientIndex = i,
        symkey = k_ij
     )

    if not validShares:
        X_ij = self.ephemeralKey[j]
        shareComplaints.append(shareComplaint(j, X_ij))
    else:
        (s_ji, t_ji) = validShares
        self.shares[j] = (s_ji, t_ji)

broadcast(messagePhase4(shareComplaints))


#
# PHASE 5
#

# Receive messages from phase 4:
# - complaints about inconsistent shares
messages.receive(4)

for complaint in messages[4]:
    j = complaint.senderIndex
    m = complaint.accusedIndex
    privkey_jm = complaint.privkey

    # Presented private key does not correspond to the published public key
    #
    # Disqualify accuser
    #
    if not validatePrivkey(
        senderIndex = j,
        recipientIndex = m,
        privkey = privkey_jm
    ):
        disqualify(5, j)
    else:
        pubkey_mj = ephemeralPubkey(m, j)

        k_jm = ecdh(privkey_jm, pubkey_mj)

        # Check whether the shares are consistent with the accused's commitments
        sharesValid = decryptAndValidateShares(
            senderIndex = m,
            recipientIndex = j,
            symkey = k_jm
        )

        # Shares inconsistent, disqualify accused
        if not sharesValid:
            disqualify(5, m)
        # Shares consistent, disqualify accuser
        else:
            disqualify(5, j)


#
# PHASE 6
#

# GJKR 2:
#
QUAL = goodParticipants[6].append(i)

# GJKR 3:
#
#   x_i = sum([ s_ji for j in QUAL ]) % q
#   x'_i = sum([ t_ji for j in QUAL ]) % q
#
x_i = sum(
    [ self.shares[j].share_S for j in QUAL ]
) % q

xprime_i = sum(
    [ self.shares[j].share_T for j in QUAL ]
) % q


#
# PHASE 7
#

# GJKR 4.(a):
#
#   A_ik = g^a_ik % p
#
self.pubkeyCoeffs = [
    P1.ecMul(A_ik) for A_ik in self.sharePolyCoeffs
]

broadcast(messagePhase7(self.pubkeyCoeffs))


#
# PHASE 8
#

# Receive messages from phase 7:
# - public key coefficients
messages.receive(7)

pubkeyComplaints = []

for j in goodParticipants[8]:
    pubkeyShareValid = validatePubkeyCoeffs(
        senderIndex = j,
        recipientIndex = i,
        share_S = self.shares[j].share_S
    )

    if not pubkeyShareValid:
        pubkeyComplaints.append(pubkeyComplaint(j))

broadcast(messagePhase8(pubkeyComplaints))


#
# PHASE 9
#

# Receive messages from phase 8:
# - complaints about invalid public key coefficients
messages.receive(8)

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

broadcast(messagePhase10(disqualifiedKeys))


#
# PHASE 11
#

# Receive messages from phase 10:
# - good participants' ephemeral private keys for each disqualified participant
messages.receive(10)

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
                self.revealedShares[m][j] = (s_mj, t_mj)

for m in disqualifiedInPhase[9]:
    shares_m = self.revealedShares[m].values
    indices_m = self.revealedShares[m].indices

    z_m = reconstruct(shares_m, indices_m)
    y_m = G.ecMul(z_m)
    self.reconstructed_Y_[m] = y_m


#
# PHASE 12
#

# GJKR 4.(c):
#
#   Y = product([ A_i0 for i in QUAL ]) % p
#
def A_(i):
    if not disqualifiedInPhase[9].contains(i):
        return pubkeyCoeffs(i)
    else:
        return [self.reconstructed_Y_[i]]

Y = ecSum(
    [ A_(i)[0] for i in QUAL ]
)
