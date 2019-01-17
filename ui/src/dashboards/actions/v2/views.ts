// Utils
import {
  readView as readViewAJAX,
  updateView as updateViewAJAX,
} from 'src/dashboards/apis/v2/'

// Types
import {RemoteDataState} from 'src/types'
import {Dispatch} from 'redux'
import {View} from 'src/types/v2'

export type Action = SetViewAction

export interface SetViewAction {
  type: 'SET_VIEW'
  payload: {
    id: string
    view: View
    status: RemoteDataState
  }
}

export const setView = (
  id: string,
  view: View,
  status: RemoteDataState
): SetViewAction => ({
  type: 'SET_VIEW',
  payload: {id, view, status},
})

export const readView = (dashboardID: string, cellID: string) => async (
  dispatch: Dispatch<Action>
): Promise<void> => {
  dispatch(setView(cellID, null, RemoteDataState.Loading))
  try {
    const view = await readViewAJAX(dashboardID, cellID)

    dispatch(setView(cellID, view, RemoteDataState.Done))
  } catch {
    dispatch(setView(cellID, null, RemoteDataState.Error))
  }
}

export const updateView = (dashboardID: string, view: View) => async (
  dispatch: Dispatch<Action>
): Promise<View> => {
  const viewID = view.cellID

  dispatch(setView(viewID, null, RemoteDataState.Loading))

  try {
    const newView = await updateViewAJAX(dashboardID, viewID, view)

    dispatch(setView(viewID, newView, RemoteDataState.Done))
    return newView
  } catch {
    dispatch(setView(viewID, null, RemoteDataState.Error))
  }
}
