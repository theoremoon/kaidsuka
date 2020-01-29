import * as React from "react";
import Register from "./pages/Register";
import Login from "./pages/Login";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import Index from "./pages/Index";

const App = () => (
    <Router>
        <Switch>
            <Route exact path="/" component={Index} />
            <Route exact path="/login" component={Login} />
            <Route exact path="/register" component={Register} />
        </Switch>
    </Router>

);
export default App;