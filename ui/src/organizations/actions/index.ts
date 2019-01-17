// Libraries
import {Dispatch} from 'redux'

// APIs
import {
  getOrganizations as getOrganizationsAPI,
  createOrg as createOrgAPI,
  deleteOrg as deleteOrgAPI,
  updateOrg as updateOrgAPI,
} from 'src/organizations/apis'

// Types
import {Organization} from 'src/api'

export enum ActionTypes {
  SetOrgs = 'SET_ORGS',
  AddOrg = 'ADD_ORG',
  RemoveOrg = 'REMOVE_ORG',
  EditOrg = 'EDIT_ORG',
}

export interface SetOrganizations {
  type: ActionTypes.SetOrgs
  payload: {
    organizations: Organization[]
  }
}

export type Actions = SetOrganizations | AddOrg | RemoveOrg | EditOrg

export const setOrgs = (organizations: Organization[]): SetOrganizations => {
  return {
    type: ActionTypes.SetOrgs,
    payload: {organizations},
  }
}

export interface AddOrg {
  type: ActionTypes.AddOrg
  payload: {
    org: Organization
  }
}

export const addOrg = (org: Organization): AddOrg => ({
  type: ActionTypes.AddOrg,
  payload: {org},
})

export interface RemoveOrg {
  type: ActionTypes.RemoveOrg
  payload: {
    org: Organization
  }
}

export const removeOrg = (org: Organization): RemoveOrg => ({
  type: ActionTypes.RemoveOrg,
  payload: {org},
})

export interface EditOrg {
  type: ActionTypes.EditOrg
  payload: {
    org: Organization
  }
}

export const editOrg = (org: Organization): EditOrg => ({
  type: ActionTypes.EditOrg,
  payload: {org},
})

// Async Actions

export const getOrganizations = () => async (
  dispatch: Dispatch<SetOrganizations>
): Promise<void> => {
  try {
    const organizations = await getOrganizationsAPI()
    dispatch(setOrgs(organizations))
  } catch (e) {
    console.error(e)
  }
}

export const createOrg = (org: Organization) => async (
  dispatch: Dispatch<AddOrg>
): Promise<void> => {
  try {
    const createdOrg = await createOrgAPI(org)
    dispatch(addOrg(createdOrg))
  } catch (e) {
    console.error(e)
  }
}

export const deleteOrg = (org: Organization) => async (
  dispatch: Dispatch<RemoveOrg>
): Promise<void> => {
  try {
    await deleteOrgAPI(org)
    dispatch(removeOrg(org))
  } catch (e) {
    console.error(e)
  }
}

export const updateOrg = (org: Organization) => async (
  dispatch: Dispatch<EditOrg>
) => {
  try {
    const updatedOrg = await updateOrgAPI(org)
    dispatch(editOrg(updatedOrg))
  } catch (e) {
    console.error(e)
  }
}
