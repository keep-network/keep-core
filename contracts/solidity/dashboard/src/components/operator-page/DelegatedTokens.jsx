import React from 'react'
import AddressShortcut from '../AddressShortcut'
import InlineForm from '../InlineForm'

const DelegatedTokens = (props) => {
  return (
    <section id="delegated-tokens" className="tile">
      <h5>Total Delegated Tokens</h5>
      <div className="flex flex-row">
        <div className="delegated-tokens-summary flex flex-column" style={{ flex: '1' }} >
          <h2 className="balance">
            80,000 K
          </h2>
          <div>
            <h6 className="text-darker-grey">OWNER&nbsp;
              <AddressShortcut address='0x345CFE101B5860d5F8ed775C21Ad5e55806C6D9B' classNames='text-big text-darker-grey' />
            </h6>
            <h6 className="text-darker-grey">BENEFICIARY&nbsp;
              <AddressShortcut address='0x345CFE101B5860d5F8ed775C21Ad5e55806C6D9B' classNames='text-big text-darker-grey' />
            </h6>
          </div>
        </div>
        <InlineForm inputProps={{ placeholder: 'Amount' }} classNames="undelegation-form" />
      </div>
    </section>
  )
}

export default DelegatedTokens
