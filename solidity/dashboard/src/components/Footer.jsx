import React from "react"
import { LINK } from "../constants/constants";

const links = [
  { label: "Join Discord", url: LINK.discord },
  { label: "About Keep", url: LINK.keepWebsite },
  {
    label: "User Guide",
    url: LINK.stakingDocumentation,
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
      <a
        href="https://github.com/keep-network/keep-core/releases"
        className="footer__app-version"
        rel="noopener noreferrer"
        target="_blank"
      >
        {`Version ${process.env.REACT_APP_VERSION}`}
      </a>
    </footer>
  )
}

const FooterLinkItem = ({ label, url }) => (
  <li className="footer__links__item">
    <a href={url} rel="noopener noreferrer" target="_blank">
      {label}
    </a>
  </li>
)

const renderFooterLinkItem = (item) => (
  <FooterLinkItem key={item.label} {...item} />
)

export default Footer
