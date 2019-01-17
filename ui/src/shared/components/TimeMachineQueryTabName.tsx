// Libraries
import React, {PureComponent, KeyboardEvent, ChangeEvent} from 'react'

// Components
import {Input, InputType, ComponentSize} from 'src/clockface'

// Styles
import 'src/shared/components/TimeMachineQueryTab.scss'

interface Props {
  isActive: boolean
  isEditing: boolean
  queryIndex: number
  name?: string
  onEdit: () => void
  onCancelEdit: () => void
  onUpdate: (newName: string) => void
}

interface State {
  newName?: string
}

class TimeMachineQueryTabName extends PureComponent<Props, State> {
  public state: State = {newName: null}

  public render() {
    const {queryIndex, name, isEditing, onCancelEdit} = this.props
    const queryName = !!name ? name : `Query ${queryIndex + 1}`

    if (isEditing) {
      return (
        <Input
          type={InputType.Text}
          value={this.state.newName || ''}
          onChange={this.handleChange}
          onBlur={onCancelEdit}
          onKeyUp={this.handleEnterKey}
          size={ComponentSize.ExtraSmall}
          autoFocus={true}
        />
      )
    }

    return (
      <div
        className="time-machine-query-tab--name"
        onDoubleClick={this.handleDoubleClick}
      >
        {queryName}
      </div>
    )
  }

  private handleDoubleClick = () => {
    if (this.props.isActive) {
      this.props.onEdit()

      this.setState({newName: this.props.name || ''})
    }
  }

  private handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    this.setState({newName: e.target.value})
  }

  private handleEnterKey = (e: KeyboardEvent<HTMLInputElement>) => {
    switch (e.key) {
      case 'Enter':
        return this.handleUpdate()
      case 'Escape':
        return this.props.onCancelEdit()
    }
  }

  private handleUpdate() {
    const {onUpdate, onCancelEdit} = this.props
    const {newName} = this.state

    if (newName !== null) {
      onUpdate(newName)
      onCancelEdit()
    }

    this.setState({newName: null})
  }
}

export default TimeMachineQueryTabName
