import React, {useState} from "react"
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom'
import Posts from "./Posts";
import Login from './Login'
import Signup from './Signup'
import {sendRequest} from "./helpers/request"
import SinglePost from "./SinglePost";
import Docs from "./Docs";

import "./styles/main.scss"

export default function App() {
  const tokensString = window.localStorage.getItem("tokens")
  const [tokens, setTokens] = useState(JSON.parse(tokensString))
  const [loggedIn, setLoggedIn] = useState(tokensString != null)
  const [message, setMessage] = useState("")

  function handleLogout() {
    window.localStorage.removeItem("tokens")
    setLoggedIn(false)
    return sendRequest("/auth/logout", {"token": tokens.token})
      .then(res => res.json())
      .then(() => {
        setTokens({})
        setLoggedIn(false)
        setMessage("Logged out successfully")
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  if (window.location.pathname.startsWith("/docs")) {
      return <Docs />
  }

  return (
    <Router>
      <div className="main">
        <div className="container">
          <h1 className="header">Zwitter</h1>
          {
            message && <p> {message} </p>
          }
          {
            loggedIn && <p>Logged in as <b>@{tokens.user.username}</b></p>
          }
          <div>
            <Link to="/">Home </Link>
            {loggedIn ? (
              <>
                <Link to="#" onClick={handleLogout}> Logout </Link>
              </>
            ) : (
              <>
              <Link to="/login">Login </Link>
              <Link to="/signup">Signup </Link>
              </>
            )}

            <Switch>
              <Route exact path="/post/:id" render={(props) => (
                <SinglePost {...props} loggedIn={loggedIn} tokens={tokens}/>
              )}/>
              <Route exact path="/login">
                <Login loggedIn={loggedIn} setTokens={setTokens} setLoggedIn={setLoggedIn} />
              </Route>
              <Route exact path="/signup">
                <Signup loggedIn={loggedIn} />
              </Route>
              <Route exact path="/" render={(props) => (
                <Posts {...props} loggedIn={loggedIn} tokens={tokens} />
              )} />
            </Switch>
          </div>
        </div>
      </div>
    </Router>
  )
}
