// Actions
import {loadBuckets} from 'src/shared/actions/v2/queryBuilder'

// Types
import {Dispatch} from 'redux-thunk'
import {TimeMachineState} from 'src/shared/reducers/v2/timeMachines'
import {Action as QueryBuilderAction} from 'src/shared/actions/v2/queryBuilder'
import {TimeRange, ViewType} from 'src/types/v2'
import {
  Axes,
  DecimalPlaces,
  XYViewGeom,
  FieldOption,
  TableOptions,
} from 'src/types/v2/dashboards'
import {TimeMachineTab} from 'src/types/v2/timeMachine'
import {Color} from 'src/types/colors'

export type Action =
  | QueryBuilderAction
  | SetActiveTimeMachineAction
  | SetActiveTabAction
  | SetNameAction
  | SetTimeRangeAction
  | SetTypeAction
  | SetActiveQueryText
  | SubmitScriptAction
  | SetIsViewingRawDataAction
  | SetGeomAction
  | SetDecimalPlaces
  | SetBackgroundThresholdColoringAction
  | SetTextThresholdColoringAction
  | SetAxes
  | SetStaticLegend
  | SetColors
  | SetYAxisLabel
  | SetYAxisMinBound
  | SetYAxisMaxBound
  | SetYAxisPrefix
  | SetYAxisSuffix
  | SetYAxisBase
  | SetYAxisScale
  | SetPrefix
  | SetSuffix
  | IncrementSubmitToken
  | SetActiveQueryIndexAction
  | AddQueryAction
  | RemoveQueryAction
  | ToggleQueryAction
  | EditActiveQueryAsFluxAction
  | EditActiveQueryAsInfluxQLAction
  | EditActiveQueryWithBuilderAction
  | UpdateActiveQueryNameAction
  | SetFieldOptionsAction
  | SetTableOptionsAction
  | SetTimeFormatAction

interface SetActiveTimeMachineAction {
  type: 'SET_ACTIVE_TIME_MACHINE'
  payload: {
    activeTimeMachineID: string
    initialState: Partial<TimeMachineState>
  }
}

export const setActiveTimeMachine = (
  activeTimeMachineID: string,
  initialState: Partial<TimeMachineState> = {}
): SetActiveTimeMachineAction => ({
  type: 'SET_ACTIVE_TIME_MACHINE',
  payload: {activeTimeMachineID, initialState},
})

interface SetActiveTabAction {
  type: 'SET_ACTIVE_TAB'
  payload: {activeTab: TimeMachineTab}
}

export const setActiveTab = (
  activeTab: TimeMachineTab
): SetActiveTabAction => ({
  type: 'SET_ACTIVE_TAB',
  payload: {activeTab},
})

interface SetNameAction {
  type: 'SET_VIEW_NAME'
  payload: {name: string}
}

export const setName = (name: string): SetNameAction => ({
  type: 'SET_VIEW_NAME',
  payload: {name},
})

interface SetTimeRangeAction {
  type: 'SET_TIME_RANGE'
  payload: {timeRange: TimeRange}
}

export const setTimeRange = (timeRange: TimeRange): SetTimeRangeAction => ({
  type: 'SET_TIME_RANGE',
  payload: {timeRange},
})

interface SetTypeAction {
  type: 'SET_VIEW_TYPE'
  payload: {type: ViewType}
}

export const setType = (type: ViewType): SetTypeAction => ({
  type: 'SET_VIEW_TYPE',
  payload: {type},
})

interface SetActiveQueryText {
  type: 'SET_ACTIVE_QUERY_TEXT'
  payload: {text: string}
}

export const setActiveQueryText = (text: string): SetActiveQueryText => ({
  type: 'SET_ACTIVE_QUERY_TEXT',
  payload: {text},
})

interface SubmitScriptAction {
  type: 'SUBMIT_SCRIPT'
}

export const submitScript = (): SubmitScriptAction => ({
  type: 'SUBMIT_SCRIPT',
})

interface SetIsViewingRawDataAction {
  type: 'SET_IS_VIEWING_RAW_DATA'
  payload: {isViewingRawData: boolean}
}

export const setIsViewingRawData = (
  isViewingRawData: boolean
): SetIsViewingRawDataAction => ({
  type: 'SET_IS_VIEWING_RAW_DATA',
  payload: {isViewingRawData},
})

interface SetGeomAction {
  type: 'SET_GEOM'
  payload: {geom: XYViewGeom}
}

export const setGeom = (geom: XYViewGeom): SetGeomAction => ({
  type: 'SET_GEOM',
  payload: {geom},
})

interface SetAxes {
  type: 'SET_AXES'
  payload: {axes: Axes}
}

export const setAxes = (axes: Axes): SetAxes => ({
  type: 'SET_AXES',
  payload: {axes},
})

interface SetYAxisLabel {
  type: 'SET_Y_AXIS_LABEL'
  payload: {label: string}
}

export const setYAxisLabel = (label: string): SetYAxisLabel => ({
  type: 'SET_Y_AXIS_LABEL',
  payload: {label},
})

interface SetYAxisMinBound {
  type: 'SET_Y_AXIS_MIN_BOUND'
  payload: {min: string}
}

export const setYAxisMinBound = (min: string): SetYAxisMinBound => ({
  type: 'SET_Y_AXIS_MIN_BOUND',
  payload: {min},
})

interface SetYAxisMaxBound {
  type: 'SET_Y_AXIS_MAX_BOUND'
  payload: {max: string}
}

export const setYAxisMaxBound = (max: string): SetYAxisMaxBound => ({
  type: 'SET_Y_AXIS_MAX_BOUND',
  payload: {max},
})

interface SetYAxisPrefix {
  type: 'SET_Y_AXIS_PREFIX'
  payload: {prefix: string}
}

export const setYAxisPrefix = (prefix: string): SetYAxisPrefix => ({
  type: 'SET_Y_AXIS_PREFIX',
  payload: {prefix},
})

interface SetYAxisSuffix {
  type: 'SET_Y_AXIS_SUFFIX'
  payload: {suffix: string}
}

export const setYAxisSuffix = (suffix: string): SetYAxisSuffix => ({
  type: 'SET_Y_AXIS_SUFFIX',
  payload: {suffix},
})

interface SetYAxisBase {
  type: 'SET_Y_AXIS_BASE'
  payload: {base: string}
}

export const setYAxisBase = (base: string): SetYAxisBase => ({
  type: 'SET_Y_AXIS_BASE',
  payload: {base},
})

interface SetYAxisScale {
  type: 'SET_Y_AXIS_SCALE'
  payload: {scale: string}
}

export const setYAxisScale = (scale: string): SetYAxisScale => ({
  type: 'SET_Y_AXIS_SCALE',
  payload: {scale},
})

interface SetPrefix {
  type: 'SET_PREFIX'
  payload: {prefix: string}
}

export const setPrefix = (prefix: string): SetPrefix => ({
  type: 'SET_PREFIX',
  payload: {prefix},
})

interface SetSuffix {
  type: 'SET_SUFFIX'
  payload: {suffix: string}
}

export const setSuffix = (suffix: string): SetSuffix => ({
  type: 'SET_SUFFIX',
  payload: {suffix},
})

interface SetStaticLegend {
  type: 'SET_STATIC_LEGEND'
  payload: {staticLegend: boolean}
}

export const setStaticLegend = (staticLegend: boolean): SetStaticLegend => ({
  type: 'SET_STATIC_LEGEND',
  payload: {staticLegend},
})

interface SetColors {
  type: 'SET_COLORS'
  payload: {colors: Color[]}
}

export const setColors = (colors: Color[]): SetColors => ({
  type: 'SET_COLORS',
  payload: {colors},
})

interface SetDecimalPlaces {
  type: 'SET_DECIMAL_PLACES'
  payload: {decimalPlaces: DecimalPlaces}
}

export const setDecimalPlaces = (
  decimalPlaces: DecimalPlaces
): SetDecimalPlaces => ({
  type: 'SET_DECIMAL_PLACES',
  payload: {decimalPlaces},
})

interface SetBackgroundThresholdColoringAction {
  type: 'SET_BACKGROUND_THRESHOLD_COLORING'
}

export const setBackgroundThresholdColoring = (): SetBackgroundThresholdColoringAction => ({
  type: 'SET_BACKGROUND_THRESHOLD_COLORING',
})

interface SetTextThresholdColoringAction {
  type: 'SET_TEXT_THRESHOLD_COLORING'
}

export const setTextThresholdColoring = (): SetTextThresholdColoringAction => ({
  type: 'SET_TEXT_THRESHOLD_COLORING',
})

interface IncrementSubmitToken {
  type: 'INCREMENT_SUBMIT_TOKEN'
}

export const incrementSubmitToken = (): IncrementSubmitToken => ({
  type: 'INCREMENT_SUBMIT_TOKEN',
})

interface EditActiveQueryWithBuilderAction {
  type: 'EDIT_ACTIVE_QUERY_WITH_BUILDER'
}

export const editActiveQueryWithBuilder = (): EditActiveQueryWithBuilderAction => ({
  type: 'EDIT_ACTIVE_QUERY_WITH_BUILDER',
})

interface EditActiveQueryAsFluxAction {
  type: 'EDIT_ACTIVE_QUERY_AS_FLUX'
}

export const editActiveQueryAsFlux = (): EditActiveQueryAsFluxAction => ({
  type: 'EDIT_ACTIVE_QUERY_AS_FLUX',
})

interface EditActiveQueryAsInfluxQLAction {
  type: 'EDIT_ACTIVE_QUERY_AS_INFLUXQL'
}

export const editActiveQueryAsInfluxQL = (): EditActiveQueryAsInfluxQLAction => ({
  type: 'EDIT_ACTIVE_QUERY_AS_INFLUXQL',
})

interface SetActiveQueryIndexAction {
  type: 'SET_ACTIVE_QUERY_INDEX'
  payload: {activeQueryIndex: number}
}

export const setActiveQueryIndexSync = (
  activeQueryIndex: number
): SetActiveQueryIndexAction => ({
  type: 'SET_ACTIVE_QUERY_INDEX',
  payload: {activeQueryIndex},
})

export const setActiveQueryIndex = (activeQueryIndex: number) => (
  dispatch: Dispatch<Action>
) => {
  dispatch(setActiveQueryIndexSync(activeQueryIndex))
  dispatch(loadBuckets())
}

interface AddQueryAction {
  type: 'ADD_QUERY'
}

export const addQuerySync = (): AddQueryAction => ({
  type: 'ADD_QUERY',
})

export const addQuery = () => (dispatch: Dispatch<Action>) => {
  dispatch(addQuerySync())
  dispatch(loadBuckets())
}

interface RemoveQueryAction {
  type: 'REMOVE_QUERY'
  payload: {queryIndex: number}
}

export const removeQuerySync = (queryIndex: number): RemoveQueryAction => ({
  type: 'REMOVE_QUERY',
  payload: {queryIndex},
})

interface ToggleQueryAction {
  type: 'TOGGLE_QUERY'
  payload: {queryIndex: number}
}

export const toggleQuerySync = (queryIndex: number): ToggleQueryAction => ({
  type: 'TOGGLE_QUERY',
  payload: {queryIndex},
})

export const removeQuery = (queryIndex: number) => (
  dispatch: Dispatch<Action>
) => {
  dispatch(removeQuerySync(queryIndex))
  dispatch(loadBuckets())
}

export const toggleQuery = (queryIndex: number) => (
  dispatch: Dispatch<Action>
) => {
  dispatch(toggleQuerySync(queryIndex))
  dispatch(submitScript())
}

interface UpdateActiveQueryNameAction {
  type: 'UPDATE_ACTIVE_QUERY_NAME'
  payload: {queryName: string}
}

export const updateActiveQueryName = (
  queryName: string
): UpdateActiveQueryNameAction => ({
  type: 'UPDATE_ACTIVE_QUERY_NAME',
  payload: {queryName},
})

interface SetFieldOptionsAction {
  type: 'SET_FIELD_OPTIONS'
  payload: {
    fieldOptions: FieldOption[]
  }
}

export const setFieldOptions = (
  fieldOptions: FieldOption[]
): SetFieldOptionsAction => ({
  type: 'SET_FIELD_OPTIONS',
  payload: {fieldOptions},
})

interface SetTableOptionsAction {
  type: 'SET_TABLE_OPTIONS'
  payload: {
    tableOptions: TableOptions
  }
}

export const setTableOptions = (
  tableOptions: TableOptions
): SetTableOptionsAction => ({
  type: 'SET_TABLE_OPTIONS',
  payload: {tableOptions},
})

interface SetTimeFormatAction {
  type: 'SET_TIME_FORMAT'
  payload: {
    timeFormat: string
  }
}

export const setTimeFormat = (timeFormat: string): SetTimeFormatAction => ({
  type: 'SET_TIME_FORMAT',
  payload: {timeFormat},
})
