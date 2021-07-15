import { Address } from '@graphprotocol/graph-ts'
import { ManagedGrant } from "../../generated/templates"
import { Grant, Delegation } from '../../generated/schema'
import {  ManagedGrant as ManagedGrantContract, GranteeReassignmentConfirmed} from '../../generated/templates/ManagedGrant/ManagedGrant'
import { ManagedGrantCreated } from '../../generated/ManagedGrantFactory/ManagedGrantFactory'
import { TokenGrant, TokenGrantCreated, TokenGrantStaked, TokenGrantWithdrawn} from '../../generated/TokenGrant/TokenGrant'

export function managedGrantCreated(event: ManagedGrantCreated) : void {
    let managedGrantAddress = event.params.grantAddress
    // Start indexing the `ManagedGrant` contract instance at
    // `event.params.grantAddress`. `event.params.grantAddress`is the address of
    // the new `ManagedGrant` contract
    ManagedGrant.create(managedGrantAddress)

    let grant = getOrCreateManagedGrant(managedGrantAddress)
    grant.grantee = event.params.grantee.toHexString()

    grant.save()
}

export function granteeReassignmentConfirmed(event: GranteeReassignmentConfirmed): void {
    let grant = getOrCreateManagedGrant(event.address)
    grant.grantee = event.params.newGrantee.toHexString()

    grant.save()
}

export function getOrCreateManagedGrant(managedGrantAddress: Address): Grant{
    let managedGrantContract = ManagedGrantContract.bind(managedGrantAddress)
    let grantId = managedGrantContract.grantId().toString()

    let grant = Grant.load(grantId)

    if (grant === null) {
        grant = new Grant(grantId)   
    }

    grant.isManagedGrant = true
    grant.managedGrantContract = managedGrantAddress.toHexString()

    return grant as Grant
}

export function tokenGrantCreated(event: TokenGrantCreated): void {
    let grantId = event.params.id.toString()
    let grant = Grant.load(grantId)

    if(grant !== null) {
        // Grant exists.
        return
    }

    // Creating the new grant.
    grant = new Grant(event.params.id.toString())

    let tokenGrantContract = TokenGrant.bind(event.address)

    // The `getGrant` function returns the grant details in the following order: 
    // value0 -> amount
    // value1 -> withdrawn
    // value2 -> staked
    // value3 -> revokedAmount
    // value4 -> revokedAt
    // value5 -> grantee
    let details = tokenGrantContract.getGrant(event.params.id)
    grant.amount = details.value0
    grant.withdrawn = details.value1
    grant.staked = details.value2
    grant.revokedAmount = details.value3
    grant.revokedAt = details.value4
    grant.grantee = details.value5.toHexString()

     // The `getGrantUnlockingSchedule` function returns the unlocking schedule
     // details in the following order:
     // value0 -> grantManager
     // value1 -> duration
     // value2 -> start
     // value3 -> clff
     // value4 -> policy
    let unlockingSchedule = tokenGrantContract.getGrantUnlockingSchedule(event.params.id)
    grant.grantManager = unlockingSchedule.value0.toHexString()
    grant.duration = unlockingSchedule.value1
    grant.start = unlockingSchedule.value2
    grant.cliff = unlockingSchedule.value3
    grant.stakingPolicy = unlockingSchedule.value4.toHexString()

    grant.save()
}

export function tokenGrantStaked(event: TokenGrantStaked): void{
    let grantId = event.params.grantId.toString()
    let operator = event.params.operator.toHexString()

    let grant = Grant.load(grantId)
    grant.staked = grant.staked.plus(event.params.amount)
    grant.save()
    
    let delegation = Delegation.load(operator)
    delegation.grant = grant.id

    delegation.save()
}

export function tokenGrantWithdrawn(event: TokenGrantWithdrawn): void {
    let grantId = event.params.grantId.toString()

    let grant = Grant.load(grantId)
    grant.withdrawn = grant.withdrawn.plus(event.params.amount)

    grant.save()
}