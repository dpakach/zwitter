import * as React from "react"
import {useState, useEffect} from "react"
import { tokenToString } from "typescript"
import {get, post, sendFileUploadRequest} from "./helpers/request"
import Post from "./Post"
import { CreatePostRequest, Tokens } from "./types/types"

type PostsProps = {
  tokens: Tokens, 
  loggedIn: boolean,
}
function Posts(props: PostsProps) {
  const {loggedIn, tokens} = props
  const [postText, setPostText]: [string, (prop: string) => void] = useState("")
  const [posts, setPosts] = useState([])
  const [message, setMessage]: [string, (prop: string) => void] = useState("")

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

    let file: File = null
    let fileid: string
    try {
      const input = document.getElementById("create-post-media-input") as HTMLInputElement
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

    const body: CreatePostRequest = {text: postText, media: fileid}

    return post("/posts", {body , headers: {"token": props.tokens.token}})
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
          <textarea placeholder="Create a new post" value={postText} onChange={(e) => { setPostText(e.target.value)}} name="text" />
          <input type="submit" value="Submit" />
          <input type="file" id="create-post-media-input" />
        </form>
      }
      {posts.map(post => <Post post={post} key={post.id} tokens={props.tokens} level={0} clickable={true} {...props} />)}
    </>
  )
}

export default Posts
