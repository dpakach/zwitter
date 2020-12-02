import React, {useState, useEffect} from "react"
import {get, post, sendFileUploadRequest} from "./helpers/request"
import Post from "./Post"

function Posts(props) {
  const {loggedIn, tokens} = props
  const [postText, setPostText] = useState("")
  const [posts, setPosts] = useState([])
  const [message, setMessage] = useState("")

  useEffect(() => {
    get("/posts", {headers: loggedIn ? {"token": tokens.token} : {}})
    .then(res => res.json())
    .then(json => {
      setPosts(json.posts ? json.posts.reverse() : [])
    }, (error) => {
      setMessage("Could not Get Posts: " + error.message)
    })
  }, [])

  async function handleSubmit(e) {
    e.preventDefault()

    let file = null
    let fileid
    try {
      const input = document.getElementById("create-post-media-input")
      file = input.files[0]
    } catch (e) {
    }
    if (file) {
      await sendFileUploadRequest(`/media/${file.name}`, file, {
        token: props.tokens.token
      })
      .then(body => body.json())
      .then(data => {
        fileid = data.id
      })
      if (!fileid) {
        setMessage("Failed while trying to upload media")
        return
      }
    }

    return post("/posts", {body: {text: postText, media: fileid}, headers: {"token": props.tokens.token}})
      .then(res => res.json())
      .then((json) => {
        setPosts([json.post, ...posts])
        setPostText("")
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
          <input type="file" id="create-post-media-input" />
        </form>
      }
      {posts.map(post => <Post post={post} key={post.id} tokens={props.tokens} level={0} clickable={true} {...props} />)}
    </>
  )
}

export default Posts
