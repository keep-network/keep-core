import React from "react"
import Skeleton from "./Skeleton"

const DataTableSkeleton = ({
  columns = 5,
  rows = 3,
  titleWidth = "300px",
  subtitleWidth = "164px",
}) => {
  return (
    <section style={{ padding: "2rem" }}>
      <Skeleton
        shining
        tag="h3"
        color="grey-20"
        className="mb-1"
        width={titleWidth}
      />
      <Skeleton shining tag="h4" color="grey-20" width={subtitleWidth} />
      <table className="table__skeleton">
        <thead>
          <tr>
            {Array.from(Array(columns)).map((_, index) => (
              <th key={index} />
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from(Array(rows)).map((_, rowIndex) => (
            <tr key={`data-row-${rowIndex}`} className="table__skeleton__row">
              {Array.from(Array(columns)).map((_, index) => (
                <td key={`data-row-cell-${index}`}>
                  <Skeleton shining color="grey-20" tag="h5" width="70%" />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}

export default DataTableSkeleton
