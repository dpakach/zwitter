import * as React from "react"
import {useState, useEffect} from "react"
import { Redirect, useParams } from 'react-router-dom'
import {get, post} from "./helpers/request"
import {getDate} from "./helpers/date"
import {User, Tokens, Genders } from "./types/types"

type ProfileProps = {
  tokens: Tokens,
  loggedIn: boolean,
  self?: boolean,
}

function ProfilePage(props: ProfileProps) {
  const {tokens, loggedIn, self} = props

  function getGenderName(id: number) {
    for (var enumMember in Genders) {
      var isValueProperty = parseInt(enumMember, 10) >= 0
      if (isValueProperty) {
         return Genders[enumMember].toLowerCase()
      }
    }
    return "not specified"
  }
  
  const params = useParams<Record<string, string | undefined>>();

  const [profileUser, setProfileUser] = useState(params.user)
  const [message, setMessage]: [string, (prop: string) => void] = useState("")
  const [displayName, setDisplayName]: [string, (prop: string) => void] = useState("")
  const [gender, setGender] = useState(Genders.NOT_SPECIFIED)
  const [birthday, setBirthday]: [string, (prop: string) => void] = useState("")
  const [user, setUser] = useState({} as User)


  useEffect(() => {
    let url = "/users/profile"
    if (profileUser) {
      url += `/${profileUser}`
    }
    get(url, { headers: { "token": tokens.token } })
    .then(res => res.json())
    .then(data => {
      setDisplayName(data.Profile['displayName'] || "")
      setBirthday(data.Profile['dateOfBirth'] || "")
      setGender(data.Profile['gender'] || Genders.NOT_SPECIFIED)
      setUser(data.user)
    })   
  }, [])

  function handleSubmit(e) {
    e.preventDefault()
    return post("/users/profile", {body: {profile: {userId: user.id, displayName, gender: gender, dateOfBirth: birthday}}, headers: {"token": tokens.token}})
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
              <select id="gender" value={gender} onChange={(e) => { setGender(parseInt(e.target.value))}} name="gender">
                <option value={Genders.NOT_SPECIFIED} >Not Specified</option>
                <option value={Genders.MALE} >Male</option>
                <option value={Genders.FEMALE} >Female</option>
                <option value={Genders.OTHER} >Other</option>
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

export default ProfilePage