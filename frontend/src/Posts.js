import React, {useState, useEffect} from "react"
import {baseUrl} from "./const.js"
import {sendRequest} from "./helpers/request"

function Posts(props) {
  const [postText, setPostText] = useState("")
  const [posts, setPosts] = useState([])
  const [message, setMessage] = useState("")

  const dateOptions = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' };
  const timeOptions = {  };

  useEffect(() => {
    fetch(`${baseUrl}/posts/get`, {method: "POST"})
    .then(res => res.json())
    .then(json => {
      setPosts(json.posts.reverse() || [])
    }, (error) => {
      setMessage("Could not Get Posts: " + error.message)
    })
  }, [])

  function handleSubmit(e) {
    e.preventDefault()
    return sendRequest("/posts/create", {text: postText}, {"token": props.tokens.token})
      .then(res => res.json())
      .then((json) => {
        setPosts([json.post, ...posts])
        setMessage("Successfully created a post")
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  return (
    <>
      <h2>Posts</h2>
      {!message || <p>{message}</p>}
      { !props.loggedIn ||
        <form onSubmit={handleSubmit}>
          <textarea placeholder="Create a new post" type="text" value={postText} onChange={(e) => { setPostText(e.target.value)}} name="text" />
          <input type="submit" value="Submit" />
        </form>
      }
      {posts.map(post => {
        let created = new Date(post.created * 1000)
        const format = created.toLocaleTimeString("en-US", dateOptions)
        return (<div key={post.id}>
          {post.text}
          <br/>
          <b>@{post.author.username}</b>
          <br/>
          {format}
          <br/>
          <br/>
        </div>)
      })}
    </>
  )
}

export default Posts
