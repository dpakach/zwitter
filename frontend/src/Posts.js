import React, {useState, useEffect} from "react"
import {sendRequest} from "./helpers/request"
import Post from "./Post"

function Posts(props) {
  const [postText, setPostText] = useState("")
  const [posts, setPosts] = useState([])
  const [message, setMessage] = useState("")

  useEffect(() => {
    sendRequest("/posts/get")
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
      {posts.map(post => <Post post={post} key={post.id} tokens={props.tokens} level={0} {...props} />)}
    </>
  )
}

export default Posts
