// Libraries
import React, {PureComponent} from 'react'

// Components
import LineOptions from 'src/shared/components/view_options/LineOptions'
import GaugeOptions from 'src/shared/components/view_options/GaugeOptions'
import SingleStatOptions from 'src/shared/components/view_options/SingleStatOptions'
import TableOptions from 'src/shared/components/view_options/TableOptions'

// Types
import {ViewType, View, NewView} from 'src/types/v2'

interface Props {
  view: View | NewView
}

class OptionsSwitcher extends PureComponent<Props> {
  public render() {
    const {view} = this.props

    switch (view.properties.type) {
      case ViewType.XY:
      case ViewType.LinePlusSingleStat:
        return <LineOptions {...view.properties} />
      case ViewType.Gauge:
        return <GaugeOptions {...view.properties} />
      case ViewType.SingleStat:
        return <SingleStatOptions />
      case ViewType.Table:
        return <TableOptions />
      default:
        return <div />
    }
  }
}

export default OptionsSwitcher
