import React, {PureComponent, MouseEvent} from 'react'
import classnames from 'classnames'

import Cells from 'src/shared/components/cells/Cells'
import FancyScrollbar from 'src/shared/components/fancy_scrollbar/FancyScrollbar'
import DashboardEmpty from 'src/dashboards/components/dashboard_empty/DashboardEmpty'

import {Dashboard, Cell} from 'src/types/v2'
import {TimeRange} from 'src/types'

import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props {
  dashboard: Dashboard
  timeRange: TimeRange
  autoRefresh: number
  manualRefresh: number
  inPresentationMode: boolean
  inView: (cell: Cell) => boolean
  onDeleteCell: (cell: Cell) => void
  onCloneCell: (cell: Cell) => void
  onZoom: (range: TimeRange) => void
  onPositionChange: (cells: Cell[]) => void
  setScrollTop: (e: MouseEvent<HTMLElement>) => void
  onEditView: (cellID: string) => void
  onAddCell: () => void
}

@ErrorHandling
class DashboardComponent extends PureComponent<Props> {
  public render() {
    const {
      onZoom,
      dashboard,
      timeRange,
      autoRefresh,
      manualRefresh,
      onDeleteCell,
      onCloneCell,
      onEditView,
      onPositionChange,
      inPresentationMode,
      setScrollTop,
      onAddCell,
    } = this.props

    return (
      <FancyScrollbar
        className={classnames('page-contents', {
          'presentation-mode': inPresentationMode,
        })}
        setScrollTop={setScrollTop}
      >
        <div className="dashboard container-fluid full-width">
          {dashboard.cells.length ? (
            <Cells
              onZoom={onZoom}
              timeRange={timeRange}
              autoRefresh={autoRefresh}
              manualRefresh={manualRefresh}
              cells={dashboard.cells}
              onCloneCell={onCloneCell}
              onDeleteCell={onDeleteCell}
              onPositionChange={onPositionChange}
              onEditView={onEditView}
            />
          ) : (
            <DashboardEmpty onAddCell={onAddCell} />
          )}
          {/* This element is used as a portal container for note tooltips in cell headers */}
          <div className="cell-header-note-tooltip-container" />
        </div>
      </FancyScrollbar>
    )
  }
}

export default DashboardComponent
