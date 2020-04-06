import React from 'react'

export class DataTable extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      headers: [],
      fields: [],
    }
  }

  componentDidMount() {
    this.initializeDataTable()
  }

  initializeDataTable = () => {
    const headers = []
    const fields = []
    React.Children.forEach(this.props.children, (children) => {
      headers.push(children.props.header)
      fields.push({
        fieldName: children.props.field,
        renderContent: children.props.renderContent,
      })
    })
    this.setState({ headers, fields })
  }

  renderItemRow = (item) => {
    const columns = this.state.fields.map((field, index) => {
      return (
        <td key={`${item[this.props.itemFieldId]}-${field.fieldName}-${item[field.fieldName]}`}>
          <span className="responsive-header">{this.state.headers[index]}</span>
          {field.renderContent ?
            field.renderContent(item) :
            item[field.fieldName]
          }
        </td>
      )
    })

    return (
      <tr key={item[this.props.itemFieldId]}>
        {columns}
      </tr>
    )
  }

  renderHeader = (header) => (
    <th key={header}>
      {header}
    </th>
  )

  render() {
    return (
      <table>
        <thead>
          <tr>
            {this.state.headers.map(this.renderHeader)}
          </tr>
        </thead>
        <tbody>
          {this.props.data.map(this.renderItemRow)}
        </tbody>
      </table>
    )
  }
}

export const Column = ({ header, field, renderContent }) => null
