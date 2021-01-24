import * as React from 'react';
import { useState } from 'react';
import { Redirect } from 'react-router-dom';
import { post } from './helpers/request';
import { UserRequest } from './types/types';

type signupProps = {
  loggedIn: boolean;
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function Signup(props: signupProps) {
  const [username, setUsername]: [string, (username: string) => void] = useState('');
  const [password, setPassword]: [string, (password: string) => void] = useState('');
  const [message, setMessage]: [string, (message: string) => void] = useState('');
  const [completed, setCompleted]: [boolean, (completed: boolean) => void] = useState<boolean>(false);

  function handleSubmit(e: React.FormEvent) {
    const body: UserRequest = { username, password };
    e.preventDefault();
    return post('/users', { body })
      .then((res) => res.json())
      .then(
        () => {
          setMessage('Success: Created User');
          setTimeout(() => {
            setCompleted(true);
          }, 2000);
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }
  return (
    <>
      {!(completed || props.loggedIn) || <Redirect to="/" />}

      <p>{message}</p>
      <h2>Signup</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Username:
          <input
            type="text"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
            }}
            name="username"
          />
        </label>
        <label>
          Password:
          <input
            type="password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
            name="password"
          />
        </label>
        <input type="submit" value="Submit" />
      </form>
    </>
  );
}
