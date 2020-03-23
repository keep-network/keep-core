import React from 'react'
import * as Icons from './Icons'

const socialMedia = [
  { label: 'Twitter', url: 'https://twitter.com/keep_project' },
  { label: 'Telegram', url: 'https://t.me/KeepNetworkOfficial' },
  { label: 'Reddit', url: 'https://www.reddit.com/r/KeepNetwork' },
]

const aboutUs = [
  { label: 'Whitepaper', url: 'https://keep.network/whitepaper' },
  { label: 'Team', url: 'https://keep.network/#team' },
  { label: 'Advisors', url: 'https://keep.network/#advisors' },
  { label: 'Blog', url: 'https://blog.keep.network/' },
]

const Footer = () => {
  return (
    <footer>
      <Icons.KeepCircle color='#F2F2F2' />
      <ul>
        {aboutUs.map(renderFooterLinkItem)}
      </ul>
      <ul>
        {socialMedia.map(renderFooterLinkItem)}
      </ul>
      <div className="signature text-smaller text-grey-70">
        <div>
          A Thesis* Build
        </div>
        <div>
          © 2020 Keep, SEZC. All Rights Reserved.
        </div>
      </div>
    </footer>
  )
}

const FooterLinkItem = ({ label, url }) => (
  <li>
    <a href={url} rel="noopener noreferrer" target="_blank" >
      <h5>{label}</h5>
    </a>
  </li>
)

const renderFooterLinkItem = (item) => <FooterLinkItem key={item.label} {...item} />

export default Footer
