import React, { useMemo } from 'react'
import Undelegations from '../components/Undelegations'
import DelegatedTokensTable from '../components/DelegatedTokensTable'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { useTokensPageContext } from '../contexts/TokensPageContext'
import { formatDate } from '../utils/general.utils'
import moment from 'moment'

const DelegationOverview = () => {
    const { 
        undelegations,
        delegations,
        refreshData,
        tokensContext,
        selectedGrant,
    } = useTokensPageContext()

    const ownedDelegations = useMemo(() => {
        return delegations.filter(delegation => !delegation.gratnId)
    }, [delegations])

    const ownedUndelegations = useMemo(() => {
        return undelegations.filter(undelegation => !undelegation.gratnId)
    }, undelegations)

    const grantDelegations = useMemo(() => {
        return delegations.filter(delegation => delegation.gratnId === selectedGrant.id)
    }, [delegations, selectedGrant])

    const grantUndelegations = useMemo(() => {
        return undelegations.filter(undelegation => undelegation.gratnId === selectedGrant.id)
    }, [undelegations, selectedGrant])

    const getDelegations = () => {
        if (tokensContext === 'granted') {
            return grantDelegations
        }

        return ownedDelegations
    }

    const getUndelegations = () => {
        if (tokensContext === 'granted') {
            return grantUndelegations
        }

        return ownedUndelegations
    }
 
    return (
        <section>
            <div className="flex wrap self-center mt-3 mb-2">
                <h2 className="text-grey-60">{`${tokensContext === 'granted' ? 'Grant ' : ''}Delegation Overview`}</h2>
                {tokensContext === 'granted' && 
                    <>
                        <span className="flex self-center ml-2">
                            <StatusBadge
                                className="self-center"
                                status={BADGE_STATUS.DISABLED}
                                text="grant id"
                            />
                            <span className="self-center h4 text-grey-50 ml-1">
                                {selectedGrant.id}
                            </span>
                        </span>
                        <span className="flex self-center ml-2">
                            <StatusBadge
                                className="self-center"
                                status={BADGE_STATUS.DISABLED}
                                text="issued"
                            />
                            <span className="h4 text-grey-50 ml-1">
                                {formatDate(moment.unix(selectedGrant.start))}
                            </span>
                        </span>
                    </>
                }
            </div>
            <DelegatedTokensTable
                delegatedTokens={getDelegations()}
                cancelStakeSuccessCallback={refreshData}
            />
            <Undelegations
                undelegations={getUndelegations()}
            />
        </section>
    )
}

export default DelegationOverview
