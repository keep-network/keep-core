import { Delegation } from '../../generated/schema'
import { StakeDelegated, OperatorStaked, TokensSeized, TokensSlashed, RecoveredStake } from '../../generated/TokenStaking/TokenStaking'

export function stakeDelegated(event: StakeDelegated): void {
    let delegation = new Delegation(event.params.operator.toHexString())
    
    delegation.owner = event.params.owner.toHexString()
    delegation.createdAt = event.block.timestamp

    delegation.save()
}

export function operatorStaked(event: OperatorStaked): void {
    let delegationId = event.params.operator.toHexString()
    let delegation = Delegation.load(delegationId)

    if (delegation == null) {
        delegation = new Delegation(delegationId)
    }

    delegation.authorizer = event.params.authorizer.toHexString()
    delegation.beneficiary = event.params.beneficiary.toHexString()
    delegation.amount = event.params.value

    delegation.save()
}

// export function topUpInitiated(event: TopUpInitiated) {
   // TODO: handle TopUpInitiated event
// }


// export function topUpCompleted(event: TopUpCompleted) {
    // TODO: handle TopUpCompleted event

// }

export function tokensSeized(event: TokensSeized): void {
    let delegationId = event.params.operator.toHexString()
    let delegation = Delegation.load(delegationId)

    if (delegation == null) {
      delegation = new Delegation(delegationId)
    }

    delegation.amount = delegation.amount.minus(event.params.amount)

    delegation.save()
}

export function tokensSlashed(event: TokensSlashed): void {
    let delegationId = event.params.operator.toHexString()
    let delegation = Delegation.load(delegationId)
    
    if (delegation == null) {
      delegation = new Delegation(delegationId)
    }

    delegation.amount = delegation.amount.minus(event.params.amount)

    delegation.save()
}

export function stakeRecovered(event: RecoveredStake): void {
    // TODO
}
