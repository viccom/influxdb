// Libraries
import React, {PureComponent} from 'react'
import {withRouter, WithRouterProps} from 'react-router'
import _ from 'lodash'

// apis
import {getSetupStatus} from 'src/onboarding/apis'

// Components
import {ErrorHandling} from 'src/shared/decorators/errors'
import SplashPage from 'src/shared/components/splash_page/SplashPage'
import SigninForm from 'src/onboarding/components/SigninForm'
import {Spinner} from 'src/clockface'
import {RemoteDataState} from 'src/types'

interface State {
  status: RemoteDataState
}

@ErrorHandling
class SigninPage extends PureComponent<WithRouterProps, State> {
  constructor(props) {
    super(props)

    this.state = {
      status: RemoteDataState.Loading,
    }
  }
  public async componentDidMount() {
    const isSetupAllowed = await getSetupStatus()

    if (isSetupAllowed) {
      this.props.router.push('/onboarding/0')
    }

    this.setState({status: RemoteDataState.Done})
  }

  public render() {
    return (
      <Spinner loading={this.state.status}>
        <SplashPage panelWidthPixels={300}>
          <SplashPage.Panel>
            <SplashPage.Logo />
            <SplashPage.Header title="InfluxData" />
            <SigninForm />
          </SplashPage.Panel>
        </SplashPage>
      </Spinner>
    )
  }
}

export default withRouter(SigninPage)
