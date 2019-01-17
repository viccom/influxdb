// Libraries
import React, {PureComponent} from 'react'

// Components
import {Form, Grid, Columns, AutoInput} from 'src/clockface'

// Constants
import {MIN_DECIMAL_PLACES, MAX_DECIMAL_PLACES} from 'src/dashboards/constants'

// Types
import {DecimalPlaces} from 'src/types/v2/dashboards'

// Decorators
import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props extends DecimalPlaces {
  onDecimalPlacesChange: (decimalPlaces: DecimalPlaces) => void
}

@ErrorHandling
class DecimalPlacesOption extends PureComponent<Props> {
  constructor(props: Props) {
    super(props)
  }

  public render() {
    return (
      <Grid.Column widthXS={Columns.Six}>
        <Form.Element label="Decimal Places">
          <AutoInput
            name="decimal-places"
            inputPlaceholder="Enter a number"
            onChange={this.handleSetValue}
            value={this.value}
            min={MIN_DECIMAL_PLACES}
            max={MAX_DECIMAL_PLACES}
          />
        </Form.Element>
      </Grid.Column>
    )
  }

  public handleSetValue = (value: number): void => {
    const digits = Math.max(value, 0)
    const isEnforced = true

    this.props.onDecimalPlacesChange({digits, isEnforced})
  }

  private get value(): number {
    const {isEnforced, digits} = this.props
    if (!isEnforced) {
      return
    }

    return digits
  }
}

export default DecimalPlacesOption
