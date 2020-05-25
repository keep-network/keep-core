import React from "react"
import Tile from "../Tile"
import * as Icons from "../Icons"

// TODO link
const docs = [
  {
    title: "Run Random Beacon",
    link:
      "https://github.com/keep-network/keep-core/blob/master/docs/run-random-beacon.adoc",
  },
  { title: "GitBook Staking Documentation", link: "https://github.com/" },
  { title: "Staking Overview", link: "https://github.com/" },
]

const DocumentationSection = () => {
  return (
    <Tile
      title="Documentation"
      titleClassName={"h3 text-grey-70"}
      id="documentation"
    >
      <ul className="documentation-list">{docs.map(renderDocItem)}</ul>
    </Tile>
  )
}

const renderDocItem = (item) => <DocItem key={item.title} {...item} />

const DocItem = ({ title, link }) => (
  <li>
    <Icons.DocumentWithBg />
    <a
      href={link}
      className="arrow-link"
      rel="noopener noreferrer"
      target="_blank"
    >
      {title}
    </a>
  </li>
)

export default DocumentationSection
