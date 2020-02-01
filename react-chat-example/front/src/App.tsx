import * as React from "react";
import { BrowserRouter, Switch, Route, RouteProps } from "react-router-dom";
import AuthRoute from "./AuthRoute";
import Index from "./Index";

type KeyExtract<T, U extends keyof T> = { [K in Extract<keyof T, U>]: T[K] };
type KeyExclude<T, U extends keyof T> = { [K in Exclude<keyof T, U>]: T[K] };

const Private: React.FunctionComponent = () => {
  return (
    <>
      <div>Private</div>
    </>
  );
};

const PrivateRoute = (props: RouteProps): React.ReactElement => {
  const { component, ...remain } = props;
  return (
    <AuthRoute
      isLoggedIn={true}
      component={component}
      loginPath="/login"
      {...remain}
    />
  );
};

const App: React.FunctionComponent = () => {
  return (
    <>
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component={Index} />
          <Route
            exact
            path="/login"
            component={(): React.ReactElement => (
              <>
                <div>Login</div>{" "}
              </>
            )}
          />
          <PrivateRoute exact path="/private" component={Private} />
        </Switch>
      </BrowserRouter>
    </>
  );
};

export default App;
