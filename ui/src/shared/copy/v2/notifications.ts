import {Notification} from 'src/types'
import {FIVE_SECONDS, TEN_SECONDS, INFINITE} from 'src/shared/constants/index'

type NotificationExcludingMessage = Pick<
  Notification,
  Exclude<keyof Notification, 'message'>
>

const defaultErrorNotification: NotificationExcludingMessage = {
  type: 'error',
  icon: 'alert-triangle',
  duration: TEN_SECONDS,
}

const defaultSuccessNotification: NotificationExcludingMessage = {
  type: 'success',
  icon: 'checkmark',
  duration: FIVE_SECONDS,
}

export const taskNotCreated = (additionalMessage: string): Notification => ({
  ...defaultErrorNotification,
  message: `Failed to create new task: ${additionalMessage}`,
})

export const taskNotFound = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to find task',
})

export const tasksFetchFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to get tasks from server',
})

export const taskDeleteFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to delete task',
})

export const taskUpdateFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to update task',
})

export const taskImportFailed = (
  fileName: string,
  errorMessage: string
): Notification => ({
  ...defaultErrorNotification,
  duration: INFINITE,
  message: `Failed to import Task from file ${fileName}: ${errorMessage}.`,
})

export const taskImportSuccess = (fileName: string): Notification => ({
  ...defaultSuccessNotification,
  duration: FIVE_SECONDS,
  message: `Successfully imported file ${fileName}.`,
})

export const getTelegrafConfigFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to get telegraf config',
})

export const savingNoteFailed = (error: string): Notification => ({
  ...defaultErrorNotification,
  duration: FIVE_SECONDS,
  message: `Failed to save note: ${error}`,
})

export const writeLineProtocolFailed = (error: string): Notification => ({
  ...defaultErrorNotification,
  message: `Failed to write line protocol ${error}`,
})

export const labelCreateFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to create label',
})

export const labelDeleteFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to delete label',
})

export const labelUpdateFailed = (): Notification => ({
  ...defaultErrorNotification,
  message: 'Failed to update label',
})
