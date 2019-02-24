class GroupSelectionContext:
    stakers: Map[StakerID, Staker]
    tickets: TicketDataStructure
    V_i:     BlsSignature
    T_init:  BlockHeight


    def T_elapsed(T_amount: BlockHeight) -> bool:
        """
        Determine if a given amount of time has elapsed (in blocks)
        """
        T_now = getCurrentBlockHeight()
        return T_now - T_init >= T_amount


    def receiveTicket(t: Ticket):
        if t in tickets:
            ignore(t)
        elif cheapCheck(t):
            tickets.add(t)
        else:
            punish(t.sender, INVALID_TICKET_PENALTY)


    def cheapCheck(t: Ticket) -> bool:
        """
        Cheap check to perform on-chain on every ticket reception
        """
        t_staker = stakers[t.sender]

        valid_Q_j = t.proof.Q_j == t_staker.address
        valid_vs = t_staker.weight() > t.proof.vs >= 1

        return valid_Q_j && valid_vs


    def receiveChallenge(challenge: Challenge):
        if not challenge.ticket in tickets:
            # Challenged ticket not found
            # eg. it has already been successfully challenged and removed
            ignore(challenge)

        elif costlyCheck(challenge):
            punish(challenge.sender, INVALID_CHALLENGE_PENALTY)

        else:
            tickets.remove(challenge.ticket)

            punish(challenge.ticket.sender, INVALID_TICKET_PENALTY)

            reward(challenge.sender, TICKET_TATTLETALE_REWARD)


    def costlyCheck(challenge: Challenge) -> bool:
        """
        Costly check to perform on-chain when a ticket is challenged;
        assumes the ticket passes cheapCheck()

        Not very costly in current version (only sha3) but enables forward
        compatibility with future ZKP implementations
        """
        t = challenge.ticket
        p = t.proof

        ticket_valid = getTicketValue(V_i, p.Q_j, p.vs) == t.value
        return ticket_valid


    def runGroupSelection() -> List[StakerID]:
        """
        Run the entire group selection protocol:

        Start by receiving tickets or challenges,
        perform a cheapCheck on each received ticket

        After timeout for tickets is over, wait a bit for challenges

        For each received challenge, perform costlyCheck
        and punish the misbehaving party
        """
        T_init = getCurrentBlockHeight()

        while not T_elapsed(TICKET_SUBMISSION_TIMEOUT):
            t = receive()

            if isTicket(t):
                receiveTicket(t)
            elif isChallenge(t):
                receiveChallenge(t)
            else:
                ignore(t)

        while not T_elapsed(TICKET_CHALLENGE_TIMEOUT):
            c = receive()

            if isChallenge(c):
                receiveChallenge(c)
            else:
                ignore(c)

        # Now we have received the tickets and challenges
        # Time to select the group candidate

        # Get the N tickets with the lowest values
        bestTickets = tickets.query(N)
        # P is the virtual stakers corresponding to those tickets
        P = bestTickets.map(sender)

        # Staker S_i may be represented multiple times in P if they have many
        # tokens and get lucky. In this case we simply proceed normally, with
        # multiple participants P_j, P_k etc. corresponding to the same actual
        # staker S_i.
        return P


class TicketDataStructure:
    """
    TicketDataStructure is an abstraction which takes in tickets
    and keeps them sorted in order: lowest ticket value first
    """
    def add(t: Ticket):
    def remove(t: Ticket):
    def query(n: int) -> List[Ticket]:


class Ticket:
    """
    A Ticket contains a pseudorandom value which is used to determine whether
    a given virtual staker is in the candidate group P
    """
    value:  Sha3Digest # corresponds to W_k
    proof:  TicketProof
    sender: StakerID


class TicketProof:
    """
    A TicketProof is an abstraction for the information necessary to determine
    whether a ticket is valid
    """
    vs:  int
    Q_j: Address


class Challenge:
    """
    A Challenge is a signed claim that a given ticket is invalid
    """
    ticket:   Ticket
    sender:   StakerID


def getTicketValue(
        V_i: BlsSignature,
        Q_j: Address,
        vs:  int
) -> Sha3Digest:
    """
    Utility function to clarify how ticket values are determined with sha3
    """
    return sha3(asBytes(V_i) ++ asBytes(Q_j) ++ asBytes(vs))


class Staker:
    address:      Address
    ecdsaPubkey:  EcdsaPubkey
    blsPubkey:    BlsPubkey
    stakedTokens: TokenAmount


    def weight() -> int:
        """
        A staker's weight is how many minimum-stake stakers a given actual staker
        could form if they were to blitzpants their stake

        By creating a number of virtual stakers corresponding to the weight of each
        actual staker, we remove the incentive to blitzpants
        """
        return floor(stakedTokens / MINIMUM_STAKE)


    def tickets(V_i: BlsSignature) -> List[Ticket]:
        """
        Generate all the tickets of a given staker, in sorted order

        The leading tickets in the resulting list can be queried for promising ones
        """
        virtualStakers = range(1, weight())
        Q_j = s.address

        createdTickets = []
        for vs in virtualStakers:
            W_k = getTicketValue(V_i, Q_j, vs)
            newTicket = Ticket(W_k, TicketProof(vs, Q_j))
            createdTickets.append(newTicket)

        return createdTickets.sort_ascending(value)


    def runGroupSelection():
        unpublished = tickets()
        tNat = naturalThreshold()

        # Publish all tickets that fall under the natural threshold
        while unpublished.first().value < tNat:
            promisingTicket = unpublished.pop()
            publish(promisingTicket)

        while not T_elapsed(TICKET_INITIAL_TIMEOUT):
            wait()

        while not T_elapsed(TICKET_SUBMISSION_TIMEOUT):
            t = getLatestSubmission()

            if isTicket(t) and not isValid(t):
                c = Challenge(t, self)
                publish(c)

            # If it seems like we would have another ticket eligible for P,
            # publish the most promising ones. This prevents dumping all tickets
            # at once before others have had time to publish theirs, but should
            # be decently responsive anyway.
            currentBestTickets = getContract().tickets.query(N)
            currentBestThreshold = max(currentBestTickets.map(value))

            bestTicket = unpublished.first()
            bestTicketEligible = bestTicket.value < currentBestThreshold
            timeToPublish = bestTicket.value * TICKET_INITIAL_TIMEOUT / tNat

            if bestTicketEligible and T_elapsed(timeToPublish):
                unpublished.pop()
                publish(bestTicket)

        # Group selection done, see if we are in P
        P = getContract().tickets.query(N).map(sender)

        if P.contains(self):
            runDKG()
        else:
            obvserveDKG()


