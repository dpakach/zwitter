import React, {useState, useEffect} from "react"
import { Redirect, useParams } from 'react-router-dom'
import {get, post} from "./helpers/request"
import {getDate, getTimeStamp} from "./helpers/date"

export default function ProfilePage(props) {
  const {tokens, loggedIn, self} = props

  const genders = Object.freeze({
    "NOT_SPECIFIED": 0,
    "MALE": 1,
    "FEMALE": 2,
    "OTHER": 3,
  })

  function getGenderName(id) {
    const gender = Object.keys(genders).find(key => genders[key] === id)
    if (gender) {
      return gender.toLowerCase()
    } else {
      return "not specified"
    }
  }

  const params = useParams();

  const [profileUser] = useState(params.user)
  const [message, setMessage] = useState("")
  const [displayName, setDisplayName] = useState("")
  const [gender, setGender] = useState(0)
  const [birthday, setBirthday] = useState("")
  const [user, setUser] = useState({})


  useEffect(() => {
    let url = "/users/profile"
    if (profileUser) {
      url += `/${profileUser}`
    }
    return get(url, {headers: {"token": tokens.token}})
    .then(res => res.json())
    .then(data => {
      setDisplayName(data.Profile['displayName'] || "")
      setBirthday(data.Profile['dateOfBirth'] || "")
      setGender(genders[data.Profile['gender']] || gender.NOT_SPECIFIED)
      setUser(data.user)
    })
  }, [])

  function handleSubmit(e) {
    e.preventDefault()
    return post("/users/profile", {body: {profile: {userId: user.id, displayName, gender: parseInt(gender), dateOfBirth: birthday}}, headers: {"token": tokens.token}})
      .then(res => res.json())
      .then(() => {
        setMessage("Success: updated User Profile")
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  return (
    <>
      {
        (self && !loggedIn) && <Redirect to="/" />
      }
      <div>
        {displayName || "Name not set"} ({getGenderName(gender) || "Gender not specified"})
        <br />
        @{user.username}
        <br />
        joined {getDate(user.created)}
        <br />
        <br />
        Birthday {birthday || "Birthday not set"}
      </div>
      {
        self &&
        <>
          <h2>Edit Profile</h2>
          <form onSubmit={handleSubmit}>
            <label>
              DisplayName:
              <input type="text" value={displayName} onChange={(e) => { setDisplayName(e.target.value)}} name="displayName" />
            </label>
            <label htmlFor="gender">
              Gender:
            </label>
              <select id="gender" value={gender} onChange={(e) => { setGender(e.target.value)}} name="gender">
                <option value="0" >Not Specified</option>
                <option value="1" >Male</option>
                <option value="2" >Female</option>
                <option value="3" >Other</option>
              </select>
            <label>
              Birthday:
              <input type="date" value={birthday} onChange={(e) => {
                setBirthday(e.target.value)}} name="displayName" />
            </label>
            <input type="submit" value="Submit" />
          </form>
        </>
      }
  </>
  )
}
