/**
 * Set the currently logged in user
 */
export const SET_USER = 'SET_USER'
export function setUser(user) {
  return {
    type: SET_USER,
    user
  }
}
