import React, {useState, useEffect} from "react"
import { BrowserRouter as Router, Switch, Routes, Route, Link } from 'react-router-dom'
import Posts from "./Posts";
import Login from './Login'
import Signup from './Signup'
import {sendRequest} from "./helpers/request"

export default function App() {
  const [loggedIn, setLoggedIn] = useState(false)
  const [tokens, setTokens] = useState({})
  const [message, setMessage] = useState("")

  useEffect(() => {
    const tokensString = window.localStorage.getItem("tokens")
    if (tokensString) {
      setTokens(JSON.parse(tokensString))
      setLoggedIn(true)
    }
  }, [])

  function handleLogout() {
    return sendRequest("/auth/logout", {"token": tokens.token})
      .then(res => res.json())
      .then(() => {
        setMessage("Logged out successfully")
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  return (
    <Router>
      <h1>Zwitter</h1>
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
        <Route path="/login">
          <Login loggedIn={loggedIn} setTokens={setTokens} setLoggedIn={setLoggedIn} />
        </Route>
        <Route path="/signup">
          <Signup />
        </Route>
      </Switch>
      <Posts loggedIn={loggedIn} tokens={tokens}/>
    </Router>
  )
}
