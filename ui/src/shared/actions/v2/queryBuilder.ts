// APIs
import {
  QueryBuilderFetcher,
  CancellationError,
} from 'src/shared/apis/v2/queryBuilder'

import {bucketsAPI} from 'src/utils/api'

// Utils
import {
  getActiveQuerySource,
  getActiveQuery,
  getActiveTimeMachine,
} from 'src/shared/selectors/timeMachines'

// Types
import {Dispatch} from 'redux-thunk'
import {GetState} from 'src/types/v2'
import {RemoteDataState} from 'src/types'

const fetcher = new QueryBuilderFetcher()

export type Action =
  | SetBuilderBucketSelectionAction
  | SetBuilderBucketsAction
  | SetBuilderBucketsStatusAction
  | SetBuilderTagKeysAction
  | SetBuilderTagKeysStatusAction
  | SetBuilderTagValuesAction
  | SetBuilderTagValuesStatusAction
  | SetBuilderTagKeySelectionAction
  | SetBuilderTagValuesSelectionAction
  | AddTagSelectorAction
  | RemoveTagSelectorAction
  | SelectFunctionAction
  | SetValuesSearchTermAction
  | SetKeysSearchTermAction

interface SetBuilderBucketsStatusAction {
  type: 'SET_BUILDER_BUCKETS_STATUS'
  payload: {bucketsStatus: RemoteDataState}
}

const setBuilderBucketsStatus = (
  bucketsStatus: RemoteDataState
): SetBuilderBucketsStatusAction => ({
  type: 'SET_BUILDER_BUCKETS_STATUS',
  payload: {bucketsStatus},
})

interface SetBuilderBucketsAction {
  type: 'SET_BUILDER_BUCKETS'
  payload: {buckets: string[]}
}

const setBuilderBuckets = (buckets: string[]): SetBuilderBucketsAction => ({
  type: 'SET_BUILDER_BUCKETS',
  payload: {buckets},
})

interface SetBuilderBucketSelectionAction {
  type: 'SET_BUILDER_BUCKET_SELECTION'
  payload: {bucket: string}
}

const setBuilderBucket = (bucket: string): SetBuilderBucketSelectionAction => ({
  type: 'SET_BUILDER_BUCKET_SELECTION',
  payload: {bucket},
})

interface SetBuilderTagKeysAction {
  type: 'SET_BUILDER_TAG_KEYS'
  payload: {index: number; keys: string[]}
}

const setBuilderTagKeys = (
  index: number,
  keys: string[]
): SetBuilderTagKeysAction => ({
  type: 'SET_BUILDER_TAG_KEYS',
  payload: {index, keys},
})

interface SetBuilderTagKeysStatusAction {
  type: 'SET_BUILDER_TAG_KEYS_STATUS'
  payload: {index: number; status: RemoteDataState}
}

const setBuilderTagKeysStatus = (
  index: number,
  status: RemoteDataState
): SetBuilderTagKeysStatusAction => ({
  type: 'SET_BUILDER_TAG_KEYS_STATUS',
  payload: {index, status},
})

interface SetBuilderTagValuesAction {
  type: 'SET_BUILDER_TAG_VALUES'
  payload: {index: number; values: string[]}
}

const setBuilderTagValues = (
  index: number,
  values: string[]
): SetBuilderTagValuesAction => ({
  type: 'SET_BUILDER_TAG_VALUES',
  payload: {index, values},
})

interface SetBuilderTagValuesStatusAction {
  type: 'SET_BUILDER_TAG_VALUES_STATUS'
  payload: {index: number; status: RemoteDataState}
}

const setBuilderTagValuesStatus = (
  index: number,
  status: RemoteDataState
): SetBuilderTagValuesStatusAction => ({
  type: 'SET_BUILDER_TAG_VALUES_STATUS',
  payload: {index, status},
})

interface SetBuilderTagKeySelectionAction {
  type: 'SET_BUILDER_TAG_KEY_SELECTION'
  payload: {index: number; key: string}
}

const setBuilderTagKeySelection = (
  index: number,
  key: string
): SetBuilderTagKeySelectionAction => ({
  type: 'SET_BUILDER_TAG_KEY_SELECTION',
  payload: {index, key},
})

interface SetBuilderTagValuesSelectionAction {
  type: 'SET_BUILDER_TAG_VALUES_SELECTION'
  payload: {index: number; values: string[]}
}

const setBuilderTagValuesSelection = (
  index: number,
  values: string[]
): SetBuilderTagValuesSelectionAction => ({
  type: 'SET_BUILDER_TAG_VALUES_SELECTION',
  payload: {index, values},
})

interface AddTagSelectorAction {
  type: 'ADD_TAG_SELECTOR'
}

const addTagSelectorSync = (): AddTagSelectorAction => ({
  type: 'ADD_TAG_SELECTOR',
})

interface RemoveTagSelectorAction {
  type: 'REMOVE_TAG_SELECTOR'
  payload: {index: number}
}

const removeTagSelectorSync = (index: number): RemoveTagSelectorAction => ({
  type: 'REMOVE_TAG_SELECTOR',
  payload: {index},
})

interface SelectFunctionAction {
  type: 'SELECT_BUILDER_FUNCTION'
  payload: {name: string}
}

export const selectFunction = (name: string): SelectFunctionAction => ({
  type: 'SELECT_BUILDER_FUNCTION',
  payload: {name},
})

interface SetValuesSearchTermAction {
  type: 'SET_BUILDER_VALUES_SEARCH_TERM'
  payload: {index: number; searchTerm: string}
}

interface SetKeysSearchTermAction {
  type: 'SET_BUILDER_KEYS_SEARCH_TERM'
  payload: {index: number; searchTerm: string}
}

export const setValuesSearchTerm = (
  index: number,
  searchTerm: string
): SetValuesSearchTermAction => ({
  type: 'SET_BUILDER_VALUES_SEARCH_TERM',
  payload: {index, searchTerm},
})

export const setKeysSearchTerm = (
  index: number,
  searchTerm: string
): SetKeysSearchTermAction => ({
  type: 'SET_BUILDER_KEYS_SEARCH_TERM',
  payload: {index, searchTerm},
})

export const loadBuckets = () => async (
  dispatch: Dispatch<Action>,
  getState: GetState
) => {
  dispatch(setBuilderBucketsStatus(RemoteDataState.Loading))

  try {
    const {data} = await bucketsAPI.bucketsGet('')
    const buckets = data.buckets.map(b => b.name)
    const selectedBucket = getActiveQuery(getState()).builderConfig.buckets[0]

    dispatch(setBuilderBuckets(buckets))

    if (selectedBucket && buckets.includes(selectedBucket)) {
      dispatch(selectBucket(selectedBucket))
    } else {
      dispatch(selectBucket(buckets[0]))
    }
  } catch (e) {
    if (e instanceof CancellationError) {
      return
    }

    console.error(e)
    dispatch(setBuilderBucketsStatus(RemoteDataState.Error))
  }
}

export const selectBucket = (bucket: string) => async (
  dispatch: Dispatch<Action>
) => {
  dispatch(setBuilderBucket(bucket))
  dispatch(loadTagSelector(0))
}

export const loadTagSelector = (index: number) => async (
  dispatch: Dispatch<Action>,
  getState: GetState
) => {
  const {buckets, tags} = getActiveQuery(getState()).builderConfig

  if (!tags[index] || !buckets[0]) {
    return
  }

  const tagPredicates = tags.slice(0, index)
  const queryURL = getActiveQuerySource(getState()).links.query

  dispatch(setBuilderTagKeysStatus(index, RemoteDataState.Loading))

  try {
    const searchTerm = getActiveTimeMachine(getState()).queryBuilder.tags[index]
      .keysSearchTerm

    const keys = await fetcher.findKeys(
      index,
      queryURL,
      buckets[0],
      tagPredicates,
      searchTerm
    )

    const {key} = tags[index]

    if (!key) {
      dispatch(setBuilderTagKeySelection(index, keys[0]))
    } else if (!keys.includes(key)) {
      // Even if the selected key didn't come back in the results, let it be
      // selected anyway
      keys.unshift(key)
    }

    dispatch(setBuilderTagKeys(index, keys))
    dispatch(loadTagSelectorValues(index))
  } catch (e) {
    if (e instanceof CancellationError) {
      return
    }

    console.error(e)
    dispatch(setBuilderTagKeysStatus(index, RemoteDataState.Error))
  }
}

const loadTagSelectorValues = (index: number) => async (
  dispatch: Dispatch<Action>,
  getState: GetState
) => {
  const {buckets, tags} = getActiveQuery(getState()).builderConfig
  const tagPredicates = tags.slice(0, index)
  const queryURL = getActiveQuerySource(getState()).links.query

  dispatch(setBuilderTagValuesStatus(index, RemoteDataState.Loading))

  try {
    const key = getActiveQuery(getState()).builderConfig.tags[index].key
    const searchTerm = getActiveTimeMachine(getState()).queryBuilder.tags[index]
      .valuesSearchTerm
    const values = await fetcher.findValues(
      index,
      queryURL,
      buckets[0],
      tagPredicates,
      key,
      searchTerm
    )

    const {values: selectedValues} = tags[index]

    for (const selectedValue of selectedValues) {
      // Even if the selected values didn't come back in the results, let them
      // be selected anyway
      if (!values.includes(selectedValue)) {
        values.unshift(selectedValue)
      }
    }

    dispatch(setBuilderTagValues(index, values))
    dispatch(loadTagSelector(index + 1))
  } catch (e) {
    if (e instanceof CancellationError) {
      return
    }

    console.error(e)
    dispatch(setBuilderTagValuesStatus(index, RemoteDataState.Error))
  }
}

export const selectTagValue = (index: number, value: string) => async (
  dispatch: Dispatch<Action>,
  getState: GetState
) => {
  const tags = getActiveQuery(getState()).builderConfig.tags
  const values = tags[index].values

  let newValues: string[]

  if (values.includes(value)) {
    newValues = values.filter(v => v !== value)
  } else {
    newValues = [...values, value]
  }

  dispatch(setBuilderTagValuesSelection(index, newValues))

  if (index === tags.length - 1 && newValues.length) {
    dispatch(addTagSelector())
  } else {
    dispatch(loadTagSelector(index + 1))
  }
}

export const selectTagKey = (index: number, key: string) => async (
  dispatch: Dispatch<Action>
) => {
  dispatch(setBuilderTagKeySelection(index, key))
  dispatch(loadTagSelectorValues(index))
}

export const searchTagValues = (index: number) => async (
  dispatch: Dispatch<Action>
) => {
  dispatch(loadTagSelectorValues(index))
}

export const searchTagKeys = (index: number) => async (
  dispatch: Dispatch<Action>
) => {
  dispatch(loadTagSelector(index))
}

export const addTagSelector = () => async (
  dispatch: Dispatch<Action>,
  getState: GetState
) => {
  dispatch(addTagSelectorSync())

  const newIndex = getActiveQuery(getState()).builderConfig.tags.length - 1

  dispatch(loadTagSelector(newIndex))
}

export const removeTagSelector = (index: number) => async (
  dispatch: Dispatch<Action>
) => {
  fetcher.cancelFindValues(index)
  fetcher.cancelFindKeys(index)

  dispatch(removeTagSelectorSync(index))
  dispatch(loadTagSelector(index))
}
