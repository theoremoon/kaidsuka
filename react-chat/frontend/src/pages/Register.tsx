import React, {useState} from 'react';
import API from '../api';

const Register = () => {
  const [state, setState] = useState({
    username: '',
    error: '',
  });
  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    API.post('/register', {
      username: state.username,
    }).catch((_) => {
      setState({...state, error: 'Failed to Register'});
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
        <p><input type="submit" value="Register"/></p>
      </form>

    </>
  );
};

export default Register;
