import * as React from "react";
import { Route, Redirect, RouteProps } from "react-router-dom";

export interface AuthRouteProps {
  isLoggedIn: boolean;
  loginPath: string;
  component: React.ComponentType;
}

function AuthRoute(props: AuthRouteProps & RouteProps): React.ReactElement {
  const { isLoggedIn, loginPath, component, ...remain } = props;
  return (
    <Route
      {...remain}
      render={({ ...props }: React.Attributes): React.ReactElement =>
        isLoggedIn ? (
          React.createElement(component, props)
        ) : (
          <Redirect to={loginPath} />
        )
      }
    />
  );
}
export default AuthRoute;
