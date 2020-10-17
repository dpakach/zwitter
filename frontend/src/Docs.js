import React, {useState, useEffect} from "react"
import { RedocStandalone } from "redoc";
import {baseUrl} from "./const";

function Docs() {
  const [service, setService] = useState("auth")

  function handleChange(e) {
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
