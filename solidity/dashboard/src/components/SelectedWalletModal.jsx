import React from 'react'

const SelectedWalletModal = ({
  walletName,
  icon,
  description,
  iconDescription,
  btnText,
  btnLink,
  coinbaseQR,
}) => {
  return (
    <div className="flex column center">
      <div className= "flex full-center mb-3">
        {icon}
        <h3 className="ml-1">{walletName}</h3>
      </div>
      { iconDescription &&
        <img src={iconDescription} className="mb-3" alt={walletName} />
      }
      <span className="text-center">
        {description}
      </span>
      {walletName === 'COINBASE' && coinbaseQR}

      {(btnLink && btnText) &&
        <a
          href={btnLink}
          className="btn bt-lg btn-primary mt-3"
          target="_blank"
          rel="noopener noreferrer"
        >
          {btnText}
        </a>
      }
    </div>
  )
}

export default React.memo(SelectedWalletModal)
