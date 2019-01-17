// Libraries
import React, {PureComponent} from 'react'
import {connect} from 'react-redux'

// Components
import {
  Button,
  ComponentColor,
  ComponentSize,
  ComponentStatus,
} from 'src/clockface'

// Actions
import {submitScript} from 'src/shared/actions/v2/timeMachines'

// Utils
import {getActiveQuery} from 'src/shared/selectors/timeMachines'

// Types
import {RemoteDataState} from 'src/types'
import {AppState} from 'src/types/v2'

interface StateProps {
  submitButtonDisabled: boolean
}

interface DispatchProps {
  onSubmitScript: typeof submitScript
}

interface OwnProps {
  queryStatus: RemoteDataState
}

type Props = StateProps & DispatchProps & OwnProps

interface State {
  didClick: boolean
}

class SubmitQueryButton extends PureComponent<Props, State> {
  public state: State = {didClick: false}

  public componentDidUpdate(prevProps: Props) {
    if (
      prevProps.queryStatus === RemoteDataState.Loading &&
      this.props.queryStatus === RemoteDataState.Done
    ) {
      this.setState({didClick: false})
    }
  }

  public render() {
    return (
      <Button
        text="Submit"
        size={ComponentSize.Small}
        status={this.buttonStatus}
        onClick={this.handleClick}
        color={ComponentColor.Primary}
      />
    )
  }

  private get buttonStatus(): ComponentStatus {
    const {queryStatus, submitButtonDisabled} = this.props
    const {didClick} = this.state

    if (submitButtonDisabled) {
      return ComponentStatus.Disabled
    }

    if (queryStatus === RemoteDataState.Loading && didClick) {
      // Only show loading state for button if it was just clicked
      return ComponentStatus.Loading
    }

    return ComponentStatus.Default
  }

  private handleClick = (): void => {
    this.props.onSubmitScript()
    this.setState({didClick: true})
  }
}

const mstp = (state: AppState) => {
  const submitButtonDisabled = getActiveQuery(state).text === ''

  return {submitButtonDisabled}
}

const mdtp = {
  onSubmitScript: submitScript,
}

export default connect<StateProps, DispatchProps, OwnProps>(
  mstp,
  mdtp
)(SubmitQueryButton)
