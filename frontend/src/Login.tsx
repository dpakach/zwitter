import * as React from "react"
import {useState} from "react"
import { Redirect } from 'react-router-dom'
import {post} from "./helpers/request"
import { Tokens } from "./types/types"

type LoginProps = {
  loggedIn: boolean,
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>,
  setTokens: React.Dispatch<React.SetStateAction<Tokens>>
}

export default function Login(props: LoginProps) {
  const [username, setUsername]: [string, (username: string) => void] = useState("")
  const [password, setPassword]: [string, (password: string) => void] = useState("")
  const [message, setMessage]: [string, (message: string) => void] = useState("")
  const [completed, setCompleted]: [boolean, (completed: boolean) => void] = useState<boolean>(false)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    return post("/auth/login", {body: {username, password}})
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
