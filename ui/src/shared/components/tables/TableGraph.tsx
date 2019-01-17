import React, {PureComponent} from 'react'
import _ from 'lodash'
import {ErrorHandling} from 'src/shared/decorators/errors'

import {
  ASCENDING,
  DESCENDING,
  DEFAULT_SORT_DIRECTION,
} from 'src/shared/constants/tableGraph'
import {FluxTable} from 'src/types'
import {TableView, SortOptions} from 'src/types/v2/dashboards'
import TableGraphTransform from 'src/shared/components/tables/TableGraphTransform'
import TableGraphTable from 'src/shared/components/tables/TableGraphTable'

interface Props {
  table: FluxTable
  properties: TableView
}

interface State {
  sortOptions: SortOptions
}

@ErrorHandling
class TableGraph extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props)
    const sortField = _.get(
      props,
      'properties.tableOptions.sortBy.internalName'
    )

    this.state = {
      sortOptions: {
        field: sortField,
        direction: ASCENDING,
      },
    }
  }

  public render() {
    const {table, properties} = this.props

    return (
      <TableGraphTransform
        data={table.data}
        properties={properties}
        sortOptions={this.sortOptions}
      >
        {transformedDataBundle => (
          <TableGraphTable
            transformedDataBundle={transformedDataBundle}
            onSort={this.handleSetSort}
            properties={properties}
          />
        )}
      </TableGraphTransform>
    )
  }

  public handleSetSort = (fieldName: string) => {
    const {sortOptions} = this.state

    if (fieldName === sortOptions.field) {
      sortOptions.direction =
        sortOptions.direction === ASCENDING ? DESCENDING : ASCENDING
    } else {
      sortOptions.field = fieldName
      sortOptions.direction = DEFAULT_SORT_DIRECTION
    }
    this.setState({sortOptions})
  }

  private get sortOptions(): SortOptions {
    const {sortOptions} = this.state
    const {table} = this.props
    const headerSet = new Set(table.data[0])

    if (headerSet.has(sortOptions.field)) {
      return sortOptions
    } else if (headerSet.has('_time')) {
      return {...sortOptions, field: '_time'}
    } else if (headerSet.has('_start')) {
      return {...sortOptions, field: '_start'}
    } else if (headerSet.has('_stop')) {
      return {...sortOptions, field: '_stop'}
    } else {
      const headers = table.data[0]
      return {...sortOptions, field: headers[0]}
    }
  }
}

export default TableGraph
