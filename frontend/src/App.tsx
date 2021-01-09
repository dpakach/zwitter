import * as React from "react"
import {useState} from "react"
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom'
import Posts from "./Posts";
import Login from './Login'
import Signup from './Signup'
import ProfilePage from './ProfilePage'
import {post} from "./helpers/request"
import SinglePost from "./SinglePost";
import Docs from "./Docs";

import "./styles/main.scss"
import { Tokens } from "./types/types";

export default function App() {
  const tokensString = window.localStorage.getItem("tokens")

  const [tokens, setTokens]: [Tokens, (Tokens: Tokens) => void] = useState<Tokens>(JSON.parse(tokensString) as Tokens)
  const [loggedIn, setLoggedIn] = useState(tokensString != null)
  const [message, setMessage]: [string, (messages: string) => void] = useState("")

  function handleLogout() {
    window.localStorage.removeItem("tokens")
    setLoggedIn(false)
    return post("/auth/logout", {headers: {"token": tokens.token}})
      .then(res => res.json())
      .then(() => {
        setTokens({} as Tokens)
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
                <Link to="/profile">Profile </Link>
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
                <SinglePost {...props} loggedIn={loggedIn} tokens={tokens} />
              )}/>
              <Route exact path="/login">
                <Login loggedIn={loggedIn} setTokens={setTokens} setLoggedIn={setLoggedIn} />
              </Route>
              <Route exact path="/signup">
                <Signup loggedIn={loggedIn} setLoggedIn={setLoggedIn} />
              </Route>
              <Route exact path="/profile">
                <ProfilePage loggedIn={loggedIn} tokens={tokens} self={true}/>
              </Route>
              <Route exact path="/profile/:user">
                <ProfilePage loggedIn={loggedIn} tokens={tokens}/>
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
