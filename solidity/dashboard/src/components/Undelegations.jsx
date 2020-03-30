import React from 'react'
import { displayAmount } from '../utils/general.utils'
import AddressShortcut from './AddressShortcut'
import SpeechBubbleInfo from './SpeechBubbleInfo'
import RecoverStakeButton from './RecoverStakeButton'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'
import { DataTable, Column } from './DataTable'


const Undelegations = ({ undelegations }) => {
  return (
    <section className="tile">
      <h3 className="text-grey-60">Undelegations</h3>
      <SpeechBubbleInfo className="mt-1 mb-1">
        <span className="text-bold">Recover</span>&nbsp;undelegated tokens to return them to your token balance.
      </SpeechBubbleInfo>
      <DataTable data={undelegations} itemFieldId="operatorAddress">
        <Column header="amount" field="amount" renderContent={({ amount }) => `${displayAmount(amount)} KEEP`} />
        <Column header="status" field="undelegationStatus" renderContent={(undelegation) => {
          const undelegationStatus = undelegation.canRecoverStake ? COMPLETE_STATUS : PENDING_STATUS
          return (
            <StatusBadge
              status={BADGE_STATUS[undelegationStatus]}
              text={undelegationStatus.toLowerCase()}
            />
          )
        }}/>
        <Column
          header="beneficiary"
          field="beneficiary"
          renderContent={({ beneficiary }) => <AddressShortcut address={beneficiary} />}
        />
        <Column
          header="operator"
          field="operatorAddress"
          renderContent={({ operatorAddress }) => <AddressShortcut address={operatorAddress} />}
        />
        <Column
          header="authorizer"
          field="authorizerAddress"
          renderContent={({ authorizerAddress }) => <AddressShortcut address={authorizerAddress} />}
        />
        <Column
          header=""
          field=""
          renderContent={(undelegation) => undelegation.canRecoverStake &&
            <RecoverStakeButton
              isFromGrant={undelegation.isFromGrant}
              operatorAddress={undelegation.operatorAddress}
            />
          }
        />
      </DataTable>
    </section>
  )
}

export default Undelegations
