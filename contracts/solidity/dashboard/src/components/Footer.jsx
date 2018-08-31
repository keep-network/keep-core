import React from 'react'

const Footer = () => {
  const d = new Date()
  const year = d.getFullYear()
  return (
    <footer class="footer">
      <div class="container text-center">
        <span class="text-muted"><small>© {year} Keep. All Rights Reserved.</small></span>
      </div>
    </footer>
  )
}

export default Footer;
