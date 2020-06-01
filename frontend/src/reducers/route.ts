import Immutable from "immutable";
import {
  CREATE_ROUTE_FULFILLED,
  HttpRoute,
  LOAD_ROUTES_FAILED,
  LOAD_ROUTES_FULFILLED,
  LOAD_ROUTES_PENDING,
  UPDATE_ROUTE_FULFILLED,
  DELETE_ROUTE_FULFILLED
} from "types/route";
import { Actions } from "../types";
import { ImmutableMap } from "../typings";

export type State = ImmutableMap<{
  isLoading: boolean;
  isFirstLoaded: boolean;
  httpRoutes: Immutable.List<HttpRoute>;
}>;

const initialState: State = Immutable.Map({
  isLoading: false,
  isFirstLoaded: false,
  httpRoutes: Immutable.List([])
});

const reducer = (state: State = initialState, action: Actions): State => {
  switch (action.type) {
    case LOAD_ROUTES_PENDING: {
      return state.set("isLoading", true);
    }
    case LOAD_ROUTES_FAILED: {
      return state.set("isLoading", false);
    }
    case LOAD_ROUTES_FULFILLED: {
      state = state.set("isLoading", false);
      state = state.set("isFirstLoaded", true);
      state = state.set("httpRoutes", action.payload.httpRoutes);
      break;
    }
    case CREATE_ROUTE_FULFILLED: {
      return state.set("httpRoutes", state.get("httpRoutes").push(action.payload.route));
    }
    case UPDATE_ROUTE_FULFILLED: {
      let routes = state.get("httpRoutes");
      const index = routes.findIndex(x => x.get("name") === action.payload.route.get("name"));

      if (index >= 0) {
        state = state.setIn(["httpRoutes", index], action.payload.route);
      }

      break;
    }
    case DELETE_ROUTE_FULFILLED: {
      let routes = state.get("httpRoutes");
      const index = routes.findIndex(x => x.get("name") === action.payload.name);

      if (index >= 0) {
        state = state.deleteIn(["httpRoutes", index]);
      }

      break;
    }
  }

  return state;
};

export default reducer;