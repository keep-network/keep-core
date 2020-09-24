import React, { useState, useEffect } from "react"
import * as Icons from "./Icons"
import Button from "./Button"
import ChooseWalletAddress from "./ChooseWalletAddress"
import { isEmptyArray } from "../utils/array.utils"
import { useShowMessage, messageType } from "./Message"
import { KeepLoadingIndicator } from "./Loadable"

const LedgerModal = ({ connector, connectAppWithWallet, closeModal }) => {
  const [accounts, setAccounts] = useState([])

  const [isFetching, setIsFetching] = useState(false)
  const [accountsOffSet, setAccountsOffSet] = useState(0)
  const [ledgerVersion, setLedgerVersion] = useState("")
  const showMessage = useShowMessage()

  useEffect(() => {
    let shoudlSetState = true
    if (ledgerVersion === "LEDGER_LIVE" || ledgerVersion === "LEDGER_LEGACY") {
      setIsFetching(true)
      connector[ledgerVersion]
        .getAccounts(5, accountsOffSet)
        .then((accounts) => {
          if (shoudlSetState) {
            setAccounts(accounts)
            setIsFetching(false)
          }
        })
        .catch((error) => {
          if (shoudlSetState) {
            setIsFetching(false)
          }
          showMessage({ type: messageType.ERROR, title: error.message })
        })
    }

    return () => {
      shoudlSetState = false
    }
  }, [ledgerVersion, connector, showMessage, accountsOffSet])

  const onSelectAccount = async (account) => {
    connector[ledgerVersion].defaultAccount = account
    await connectAppWithWallet(connector[ledgerVersion], "LEDGER")
    closeModal()
  }

  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <Icons.Ledger />
        <h3 className="ml-1">Ledger</h3>
      </div>
      <Icons.LedgerDevice className="mb3" />
      <span className="text-center">Plug in Ledger device and unlock.</span>
      <div
        className="flex column mt-1"
        style={{
          alignSelf: "normal",
          justifyContent: "space-around",
        }}
      >
        <Button
          onClick={() => setLedgerVersion("LEDGER_LIVE")}
          className="btn btn-primary btn-md mb-1"
        >
          ledger live
        </Button>
        <Button
          onClick={() => setLedgerVersion("LEDGER_LEGACY")}
          className="btn btn-primary btn-md"
        >
          ledger legacy
        </Button>
      </div>
      {isFetching ? (
        <KeepLoadingIndicator />
      ) : (
        !isEmptyArray(accounts) && (
          <ChooseWalletAddress
            onSelectAccount={onSelectAccount}
            addresses={accounts}
            withPagination
            renderPrevBtn={accountsOffSet > 0}
            onNext={() => setAccountsOffSet((prevOffset) => prevOffset + 5)}
            onPrev={() => setAccountsOffSet((prevOffset) => prevOffset - 5)}
          />
        )
      )}
    </div>
  )
}

export default LedgerModal
