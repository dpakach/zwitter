import React, {useState, useEffect} from "react"
import { BrowserRouter as Router, Switch, Routes, Route, Link } from 'react-router-dom'
import Posts from "./Posts";
import Login from './Login'
import Signup from './Signup'
import {sendRequest} from "./helpers/request"
import SinglePost from "./SinglePost";

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

  return (
    <Router>
      <h1>Zwitter</h1>
      {
        loggedIn && <p>Logged in as <b>@{"<username>"}</b></p>
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
      </div>

      <Switch>
        <Route path="/post/:id" render={(props) => (
          <SinglePost {...props} loggedIn={loggedIn} tokens={tokens}/>
        )}/>
        <Route path="/login">
          <Login loggedIn={loggedIn} setTokens={setTokens} setLoggedIn={setLoggedIn} />
        </Route>
        <Route path="/signup">
          <Signup loggedIn={loggedIn} />
        </Route>
        <Route path="/" render={(props) => (
          <Posts {...props} loggedIn={loggedIn} tokens={tokens} />
        )} />
      </Switch>
    </Router>
  )
}
