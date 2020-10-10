import React, {useState, useEffect} from "react"
import {sendRequest} from "./helpers/request"
import Post from "./Post"

function SinglePost(props) {
  const [message, setMessage] = useState("")
  const [post, setPost] = useState(null)
  const [postId, setPostId] = useState(0)


  useEffect(() => {
    if (postId) {
      sendRequest("/posts/get/id", {id: postId})
      .then(res => res.json())
      .then(json => {
        setPost(json.post)
        setMessage("success")
      }, (error) => {
        setMessage(error)
      })
    }
  }, [postId])
  

  useEffect(() => {
    setPostId(props.match.params.id)
  }, [props.match.params.id])

  return (
    <>
      <h2>Post</h2>
      {message && <p>{message}</p>}
      {post && (
        <Post post={post} tokens={props.tokens} level={0} {...props} clickable={false} setPostId={setPostId} />
      )}
    </>
  )
}

export default SinglePost
