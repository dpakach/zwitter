import * as React from "react"
import {useState, useEffect} from "react"
import { RedocStandalone } from "redoc";
import {baseUrl} from "./const";

function Docs() {
  const [service, setService]: [string, (service: string) => void] = useState("auth")

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    setService(e.target.value)
  } 

  return (
    <div>
      <select value={service} onChange={handleChange}>
        <option value="auth">Auth Service</option>
        <option value="posts">Posts Service</option>
        <option value="users">Users Service</option>
      </select> 
      <RedocStandalone specUrl={`${baseUrl}/${service}/swagger.json`} />
    </div>
  )
}

export default Docs
