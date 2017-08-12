import { SET_USER } from '../actions/auth'

const auth = (state = {}, action) => {
  switch (action.type) {
    case SET_USER:
      return {
        ...state,
        user: action.user
      }
    default:
      return state;
  }
}

export default auth
