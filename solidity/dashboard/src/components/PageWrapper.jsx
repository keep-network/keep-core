import React from "react"
import { Link } from "react-router-dom"
import * as Icons from "./Icons"

const PageWrapper = ({
  title,
  children,
  nextPageLink,
  nextPageIcon,
  nextPageTitle,
  ...titleProps
}) => {
  const IconComponent = nextPageIcon
  return (
    <>
      <header className="flex row wrap center space-between">
        <h1 className="mb-2" {...titleProps}>
          {title}
        </h1>
        {nextPageLink && nextPageTitle && (
          <nav style={{ marginLeft: "auto" }}>
            <Link to={nextPageLink} style={{ textDecoration: "none" }}>
              <div className="flex center">
                <IconComponent width={30} height={30} />
                <span className="text-black text-link-page mr-1 ml-1">
                  {nextPageTitle}
                </span>
                <Icons.ArrowRight className="arrow-right lg black self-center" />
              </div>
            </Link>
          </nav>
        )}
      </header>
      {children}
    </>
  )
}

export default PageWrapper
