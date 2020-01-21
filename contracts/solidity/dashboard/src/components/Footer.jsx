import React from 'react'

const Footer = () => {
  const d = new Date()
  const year = d.getFullYear()
  return (
    <footer className="footer">
      <div className="container text-center">
        <span className="text-muted"><small>© {year} Keep. All Rights Reserved.</small></span>
      </div>
    </footer>
  )
}

export default Footer
