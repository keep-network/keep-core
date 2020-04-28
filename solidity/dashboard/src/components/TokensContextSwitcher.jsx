import React, { useState } from 'react'
import * as Icons from './Icons'
import Dropdown from './Dropdown'
import SelectedGrantDropdown from './SelectedGrantDropdown'

const TokensContextSwitcher = ({ onContextChange }) => {
    const [context, setContext] = useState('grant')
    return (
        <div className="tokens-context-switcher-wrapper">
            <div className={`grants ${context === 'grant' ? 'active' : 'inactive'}`} onClick={() => setContext('grant')}>
                <div className="flex row">
                    <Icons.GrantContextIcon />
                    <div className="ml-1">
                        <h2 className="text-grey-70">Grants</h2>
                        <h4 className="balance">1000</h4>
                    </div>
                </div>
                <div className="grants-dropdown">
                    <Dropdown
                        onSelect={() => console.log('on select grant')}
                        options={[]}
                        valuePropertyName='id'
                        labelPropertyName='id'
                        selectedItem={{}}
                        labelPrefix='Grant ID'
                        noItemSelectedText='Select Grant'
                        label="Choose Grant"
                        selectedItemComponent={<SelectedGrantDropdown grant={{}} />}
                    />
                </div>
            </div>
            <div className={`owned ${context === 'owned' ? 'active' : 'inactive'}`} onClick={() => setContext('owned')}>
                <Icons.MoneyWalletOpen />
                <h2 className="text-grey-70">Owned</h2>
                <h4 className="balance">1000</h4>
            </div>
        </div>
    )
}

export default TokensContextSwitcher
