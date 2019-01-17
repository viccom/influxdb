import {
  DEFAULT_VERTICAL_TIME_AXIS,
  DEFAULT_FIX_FIRST_COLUMN,
} from 'src/shared/constants/tableGraph'
import {Cell, Dashboard} from 'src/types/v2'
import {DecimalPlaces} from 'src/types/v2/dashboards'
import {DEFAULT_TIME_FORMAT} from 'src/shared/constants'

export const UNTITLED_GRAPH: string = 'Untitled Graph'

export const TIME_FORMAT_TOOLTIP_LINK: string =
  'http://momentjs.com/docs/#/parsing/string-format/'

export const DEFAULT_DECIMAL_PLACES: DecimalPlaces = {
  isEnforced: true,
  digits: 2,
}

export interface TimeField {
  internalName: string
  displayName: string
  visible: boolean
}
export const DEFAULT_TIME_FIELD: TimeField = {
  internalName: '_time',
  displayName: 'time',
  visible: true,
}

export const DEFAULT_TABLE_OPTIONS = {
  verticalTimeAxis: DEFAULT_VERTICAL_TIME_AXIS,
  sortBy: DEFAULT_TIME_FIELD,
  wrapping: 'truncate',
  fixFirstColumn: DEFAULT_FIX_FIRST_COLUMN,
}

export const TIME_FORMAT_CUSTOM: string = 'Custom'

export const FORMAT_OPTIONS: Array<{text: string}> = [
  {text: DEFAULT_TIME_FORMAT},
  {text: 'MM/DD/YYYY HH:mm:ss.SSS'},
  {text: 'YYYY/MM/DD HH:mm:ss'},
  {text: 'HH:mm:ss'},
  {text: 'HH:mm:ss.SSS'},
  {text: 'MMMM D, YYYY HH:mm:ss'},
  {text: 'dddd, MMMM D, YYYY HH:mm:ss'},
  {text: TIME_FORMAT_CUSTOM},
]

export type NewDefaultCell = Pick<
  Cell,
  Exclude<keyof Cell, 'id' | 'links' | 'dashboardID' | 'name'>
>

export const NEW_DEFAULT_DASHBOARD_CELL: NewDefaultCell = {
  x: 0,
  y: 0,
  w: 4,
  h: 4,
}

export type EmptyDefaultDashboard = Pick<
  Dashboard,
  Exclude<keyof Dashboard, 'templates' | 'links' | 'organization' | 'cells'>
> & {
  cells: NewDefaultCell[]
}

export const EMPTY_DASHBOARD: EmptyDefaultDashboard = {
  id: '0',
  name: '',
  cells: [NEW_DEFAULT_DASHBOARD_CELL],
}

export const DashboardTemplate: EmptyDefaultDashboard = {
  id: '0',
  name: 'Create a New Dashboard',
  cells: [],
}

type NewDefaultDashboard = Pick<
  Dashboard,
  Exclude<keyof Dashboard, 'id' | 'templates' | 'organization' | 'cells'> & {
    cells: NewDefaultCell[]
  }
>
export const DEFAULT_CELL_NAME = 'Name this Cell'
export const DEFAULT_DASHBOARD_NAME = 'Name this Dashboard'
export const NEW_DASHBOARD: NewDefaultDashboard = {
  name: DEFAULT_DASHBOARD_NAME,
  cells: [NEW_DEFAULT_DASHBOARD_CELL],
}

export const TYPE_QUERY_CONFIG: string = 'queryConfig'
export const TYPE_SHIFTED: string = 'shifted queryConfig'
export const TYPE_FLUX: string = 'flux'
export const DASHBOARD_NAME_MAX_LENGTH: number = 50
export const CELL_NAME_MAX_LENGTH: number = 68

export enum CEOTabs {
  Queries = 'Queries',
  Vis = 'Visualization',
}

export const MIN_DECIMAL_PLACES = 0
export const MAX_DECIMAL_PLACES = 10

// used in importing dashboards and mapping sources
export const DYNAMIC_SOURCE = 'dynamic'
export const DYNAMIC_SOURCE_INFO = {
  name: 'Dynamic Source',
  id: DYNAMIC_SOURCE,
  link: '',
}

export const UPDATED_AT_TIME_FORMAT = 'YYYY-MM-DD HH:MM:ss'
