import React, { useState, useEffect } from "react"
import Button from "./Button"
import { shortenAddress } from "../utils/general.utils"
import ViewAddressInBlockExplorer from "./ViewAddressInBlockExplorer"

const ChooseWalletAddress = ({ addresses, onSelectAccount }) => {
  const [selectedAccount, setAccount] = useState("")

  useEffect(() => {
    setAccount("")
  }, [addresses])

  return (
    <>
      <h4 className="mt-1 mb-1">Select account</h4>
      <ul className="choose-wallet-address">
        {addresses.map((address) => (
          <li key={address} onClick={() => setAccount(address)}>
            <label title={address}>
              <input
                type="radio"
                name="address"
                value={address}
                checked={address === selectedAccount}
              />
              {shortenAddress(address)}
            </label>
            <ViewAddressInBlockExplorer address={address} urlSuffix="" />
          </li>
        ))}
      </ul>
      <Button
        className="btn btn-primary btn-md mt-1"
        disabled={!selectedAccount}
        onClick={() => onSelectAccount(selectedAccount)}
      >
        select account
      </Button>
    </>
  )
}

export default ChooseWalletAddress
