// Libraries
import React, {PureComponent} from 'react'
import {withRouter, WithRouterProps} from 'react-router'
import {connect} from 'react-redux'
import _ from 'lodash'

// Components
import {ErrorHandling} from 'src/shared/decorators/errors'
import {
  WizardFullScreen,
  WizardProgressHeader,
  ProgressBar,
} from 'src/clockface'
import OnboardingStepSwitcher from 'src/onboarding/components/OnboardingStepSwitcher'

// Actions
import {notify as notifyAction} from 'src/shared/actions/notifications'
import {setSetupParams, setStepStatus, setupAdmin} from 'src/onboarding/actions'

// Constants
import {StepStatus} from 'src/clockface/constants/wizard'

// Types
import {Links} from 'src/types/v2/links'
import {SetupParams} from 'src/onboarding/apis'
import {Notification, NotificationFunc} from 'src/types'
import {AppState} from 'src/types/v2'

export interface OnboardingStepProps {
  links: Links
  currentStepIndex: number
  onSetCurrentStepIndex: (stepNumber: number) => void
  onIncrementCurrentStepIndex: () => void
  onDecrementCurrentStepIndex: () => void
  onSetStepStatus: (index: number, status: StepStatus) => void
  onSetSubstepIndex: (index: number, subStep: number | 'streaming') => void
  stepStatuses: StepStatus[]
  stepTitles: string[]
  setupParams: SetupParams
  handleSetSetupParams: (setupParams: SetupParams) => void
  notify: (message: Notification | NotificationFunc) => void
  onCompleteSetup: () => void
  onExit: () => void
}

interface OwnProps {
  startStep?: number
  stepStatuses?: StepStatus[]
  onCompleteSetup: () => void
  currentStepIndex: number
  onIncrementCurrentStepIndex: () => void
  onDecrementCurrentStepIndex: () => void
  onSetCurrentStepIndex: (stepNumber: number) => void
  onSetSubstepIndex: (stepNumber: number, substep: number | 'streaming') => void
}

interface DispatchProps {
  notify: (message: Notification | NotificationFunc) => void
  onSetSetupParams: typeof setSetupParams
  onSetStepStatus: typeof setStepStatus
  onSetupAdmin: typeof setupAdmin
}

interface StateProps {
  links: Links
  stepStatuses: StepStatus[]
  setupParams: SetupParams
}

type Props = OwnProps & StateProps & DispatchProps & WithRouterProps

@ErrorHandling
class OnboardingWizard extends PureComponent<Props> {
  public stepTitles = ['Welcome', 'Initial User Setup', 'Complete']

  public stepSkippable = [true, false, false]

  constructor(props: Props) {
    super(props)
  }

  public render() {
    const {currentStepIndex, setupParams, onSetupAdmin} = this.props

    return (
      <WizardFullScreen>
        {this.progressHeader}
        <div className="wizard-contents">
          <div className="wizard-step--container">
            <OnboardingStepSwitcher
              currentStepIndex={currentStepIndex}
              onboardingStepProps={this.onboardingStepProps}
              setupParams={setupParams}
              onSetupAdmin={onSetupAdmin}
            />
          </div>
        </div>
      </WizardFullScreen>
    )
  }

  private get progressHeader(): JSX.Element {
    const {stepStatuses, currentStepIndex, onSetCurrentStepIndex} = this.props

    if (currentStepIndex === 0) {
      return <div className="wizard--progress-header hidden" />
    }

    return (
      <WizardProgressHeader>
        <ProgressBar
          currentStepIndex={currentStepIndex}
          handleSetCurrentStep={onSetCurrentStepIndex}
          stepStatuses={stepStatuses}
          stepTitles={this.stepTitles}
          stepSkippable={this.stepSkippable}
        />
      </WizardProgressHeader>
    )
  }

  private handleExit = () => {
    const {router, onCompleteSetup} = this.props
    onCompleteSetup()
    router.push(`/`)
  }

  private get onboardingStepProps(): OnboardingStepProps {
    const {
      stepStatuses,
      links,
      notify,
      onCompleteSetup,
      setupParams,
      currentStepIndex,
      onSetStepStatus,
      onSetSetupParams,
      onSetCurrentStepIndex,
      onSetSubstepIndex,
      onDecrementCurrentStepIndex,
      onIncrementCurrentStepIndex,
    } = this.props

    return {
      stepStatuses,
      stepTitles: this.stepTitles,
      currentStepIndex,
      onSetCurrentStepIndex,
      onSetSubstepIndex,
      onIncrementCurrentStepIndex,
      onDecrementCurrentStepIndex,
      onSetStepStatus,
      links,
      setupParams,
      handleSetSetupParams: onSetSetupParams,
      notify,
      onCompleteSetup,
      onExit: this.handleExit,
    }
  }
}

const mstp = ({
  links,
  onboarding: {stepStatuses, setupParams},
}: AppState): StateProps => ({
  links,
  stepStatuses,
  setupParams,
})

const mdtp: DispatchProps = {
  notify: notifyAction,
  onSetSetupParams: setSetupParams,
  onSetStepStatus: setStepStatus,
  onSetupAdmin: setupAdmin,
}

export default connect<StateProps, DispatchProps, OwnProps>(
  mstp,
  mdtp
)(withRouter(OnboardingWizard))
