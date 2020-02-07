import React from 'react';
import {BrowserRouter as Router, Route, Switch, Link} from 'react-router-dom';

import Login from './pages/Login';
import Register from './pages/Register';
import Chat from './pages/Chat';
import API from './api';
import {useApolloClient} from '@apollo/react-hooks';

const App = () => {
  const apolloClient = useApolloClient();
  const onLogout = () => {
    API.post('/logout').then(() => {
      apolloClient.resetStore();
    });
  };

  return (
    <Router>
      <>
        <p>
          <Link to="/login">login</Link>
          <Link to="/register">register</Link>
          <a href="#" onClick={onLogout}>logout</a>
        </p>
      </>
      <Switch>
        <Route exact path='/login' component={Login} />
        <Route exact path='/register' component={Register} />
      </Switch>
      <Chat />
    </Router>
  );
};

export default App;
