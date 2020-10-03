import React, {useState} from "react"
import { Redirect } from 'react-router-dom'
import {sendRequest} from "./helpers/request"

export default function Login(props) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")
  const [completed, setCompleted] = useState(false)

  function handleSubmit(e) {
    e.preventDefault()
    console.log({username, password})
    return sendRequest("/auth/token/get", {username, password})
      .then(res => res.json())
      .then(json => {
        setMessage("Success: Logged In")
        window.localStorage.setItem("tokens", JSON.stringify(json))
        props.setTokens(json)
        setTimeout(() => {
          props.setLoggedIn(true)
          setCompleted(true)
        }, 2000)
      }, (error) => {
        console.log(error)
        setMessage("Error: " + error.message)
      })
  }
  return (
    <>
      {!(completed || props.loggedIn) || <Redirect to="/" />}

      <p>{message}</p>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Username:
          <input type="text" value={username} onChange={(e) => { setUsername(e.target.value)}} name="username" />
        </label>
        <label>
          Password:
          <input type="password" value={password} onChange={(e) => { setPassword(e.target.value)}} name="password" />
        </label>
        <input type="submit" value="Submit" />
      </form>
    </>
  )
}
