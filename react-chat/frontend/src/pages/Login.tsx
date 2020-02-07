import React, {useState} from 'react';
import API from '../api';
import {useApolloClient} from '@apollo/react-hooks';

const Login = () => {
  const [state, setState] = useState({
    username: '',
    error: '',
  });
  const apolloClient = useApolloClient();

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    API.post('/login', {
      username: state.username,
    })
        .then((_)=>{
          apolloClient.resetStore();
        })
        .catch((_) => {
          setState({...state, error: 'Failed to Login'});
        });
  };

  return (
    <>
      <form onSubmit={onSubmit}>
        {state.error || null}
        <p>username: <input type="text" name="username" id="" onChange={
          (e: React.ChangeEvent<HTMLInputElement>) => {
            return setState({...state, username: e.target.value});
          }} /></p>
        <p><input type="submit" value="Login"/></p>
      </form>

    </>
  );
};

export default Login;
