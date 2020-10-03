import React, {useState, useEffect} from "react"
import { BrowserRouter as Router, Switch, Routes, Route, Link } from 'react-router-dom'
import Posts from "./Posts";
import Login from './Login'
import Signup from './Signup'

export default function App() {
  const [loggedIn, setLoggedIn] = useState(false)
  const [tokens, setTokens] = useState({})

  useEffect(() => {
    const tokensString = window.localStorage.getItem("tokens")
    if (tokensString) {
      setTokens(JSON.parse(tokensString))
      setLoggedIn(true)
    }
  }, [])

  return (
    <Router>
      <h1>Zwitter</h1>
      <div>
        <Link to="/">Home </Link>
        {loggedIn ? <></> : (
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
