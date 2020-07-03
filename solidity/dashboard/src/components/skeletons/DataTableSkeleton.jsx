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
      <Skeleton className="mb-1" style={{ width: titleWidth }} />
      <Skeleton className="text-small" styles={{ width: subtitleWidth }} />
      <table className="table__skeleton">
        <thead>
          {Array.from(Array(columns)).map((_, index) => (
            <th key={index} />
          ))}
        </thead>
        <tbody>
          {Array.from(Array(rows)).map((_, rowIndex) => (
            <tr key={`data-row-${rowIndex}`} className="table__skeleton__row">
              {Array.from(Array(columns)).map((_, index) => (
                <td key={`data-row-cell-${index}`}>
                  <Skeleton styles={{ width: "70%" }} />
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
