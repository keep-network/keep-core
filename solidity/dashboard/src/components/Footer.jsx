import React from "react"

const links = [
  { label: "Join Discord", url: "https://discordapp.com/invite/wYezN7v" },
  { label: "About Keep", url: "https://keep.network/" },
  {
    label: "User Guide",
    url: "https://keep-network.gitbook.io/staking-documentation/",
  },
]

const Footer = ({ className }) => {
  return (
    <footer className={`footer ${className}`}>
      <ul className="footer__links">{links.map(renderFooterLinkItem)}</ul>
      <div className="footer__signature">
        <p>A Thesis* Build</p>
        <p>&#169; 2020 Keep, SEZC</p>
        <p>All Rights Reserved.</p>
      </div>
      <div className="footer__app-version">{`Version ${process.env.REACT_APP_VERSION}`}</div>
    </footer>
  )
}

const FooterLinkItem = ({ label, url }) => (
  <li className="footer__links__item">
    {/* TODO add arrow link to the external link*/}
    <a href={url} rel="noopener noreferrer" target="_blank">
      {label}
    </a>
  </li>
)

const renderFooterLinkItem = (item) => (
  <FooterLinkItem key={item.label} {...item} />
)

export default Footer
