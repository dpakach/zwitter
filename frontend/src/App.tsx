import * as React from 'react';
import { useState } from 'react';
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom';
import Posts from './Posts';
import Login from './Login';
import Signup from './Signup';
import ProfilePage from './ProfilePage';
import { post } from './helpers/request';
import SinglePost from './SinglePost';
import Docs from './Docs';

import './styles/main.scss';
import { Tokens } from './types/types';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faSignOutAlt, faHome, faUser, faUserPlus } from '@fortawesome/free-solid-svg-icons'

export default function App() {
  const tokensString = window.localStorage.getItem('tokens');

  const [tokens, setTokens]: [Tokens, (Tokens: Tokens) => void] = useState<Tokens>(JSON.parse(tokensString) as Tokens);
  const [loggedIn, setLoggedIn] = useState(tokensString != null);
  const [message, setMessage]: [string, (messages: string) => void] = useState('');

  function handleLogout() {
    window.localStorage.removeItem('tokens');
    setLoggedIn(false);
    return post('/auth/logout', { headers: { token: tokens.token } })
      .then((res) => res.json())
      .then(
        () => {
          setTokens({} as Tokens);
          setLoggedIn(false);
          setMessage('Logged out successfully');
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }

  if (window.location.pathname.startsWith('/docs')) {
    return <Docs />;
  }

  return (
    <Router>
      <div className="main">
        <div className="container">
          <div className="navigation">
            <h1 className="navigation__title">ZWITTER</h1>
            <ul className="navigation__list">
              <Link to="/">
                <li className="navigation__item navigation__item--active">
                  <FontAwesomeIcon className="navigation__icon" icon={faHome} />
                  <p className="navigation__label">Home</p>
                </li>
              </Link>
              {loggedIn ? (
                <>
                  <Link to="/profile">
                    <li className="navigation__item navigation__item--active">
                      <FontAwesomeIcon className="navigation__icon" icon={faUser} />
                      <p className="navigation__label">Profile</p>
                    </li>
                  </Link>

                  <Link to="#" onClick={handleLogout}>
                    <li className="navigation__item navigation__item--active">
                      <FontAwesomeIcon className="navigation__icon" icon={faSignOutAlt} />
                      <p className="navigation__label">Logout</p>
                    </li>
                  </Link>
                </>
              ) : (
                <>
                  <Link to="/login">
                    <li className="navigation__item navigation__item--active">
                      <FontAwesomeIcon className="navigation__icon" icon={faUser} />
                      <p className="navigation__label">Login</p>
                    </li>
                  </Link>
                  <Link to="/signup">
                    <li className="navigation__item navigation__item--active">
                      <FontAwesomeIcon className="navigation__icon" icon={faUserPlus} />
                      <p className="navigation__label">Signup</p>
                    </li>
                  </Link>
                </>
              )}              
            </ul>
          </div>
          
          {loggedIn && (
            <p>
              Logged in as <b>@{tokens.user.username}</b>
            </p>
          )}
          <div className="main-content">
            <Switch>
              <Route
                exact
                path="/post/:id"
                render={(props) => <SinglePost {...props} loggedIn={loggedIn} tokens={tokens} />}
              />
              <Route exact path="/login">
                <Login loggedIn={loggedIn} setTokens={setTokens} setLoggedIn={setLoggedIn} />
              </Route>
              <Route exact path="/signup">
                <Signup loggedIn={loggedIn} setLoggedIn={setLoggedIn} />
              </Route>
              <Route exact path="/profile">
                <ProfilePage loggedIn={loggedIn} tokens={tokens} self={true} />
              </Route>
              <Route exact path="/profile/:user">
                <ProfilePage loggedIn={loggedIn} tokens={tokens} />
              </Route>
              <Route
                exact
                path="/"
                render={(props) => <Posts {...props} loggedIn={loggedIn} tokens={tokens} showInput={true} />}
              />
            </Switch>
          </div>
        </div>
      </div>
    </Router>
  );
}
