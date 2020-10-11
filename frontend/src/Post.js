import React, {useState, useEffect} from "react"
import {sendRequest} from "./helpers/request"
import {Link} from "react-router-dom"

function Post({post: p, tokens, level, loggedIn, clickable}) {
  const [replyShown, setReplyShown] = useState(false)
  const [replyText, setReplyText] = useState("")
  const [post, setPost] = useState(p)
  const [message, setMessage] = useState("")

  const [updateKey, updatePage] = useState(0)

  useEffect(() => {
    setPost(p)
  }, [])

  const dateOptions = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' };
  const created = new Date(post.created * 1000)
  const formattedDate = created.toLocaleTimeString("en-US", dateOptions)

  function handleSubmit(e) {
    e.preventDefault()
    return sendRequest("/posts/create", {text: replyText, parentid: post.id}, {"token": tokens.token})
      .then(res => res.json())
      .then((json) => {
        setReplyText("")
        setPost({...post, children: [json.post, ...(post.children || [])]})
        setReplyShown(false)
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  function likePost() {
    return sendRequest("/posts/id/like", {post: post.id}, {"token": tokens.token})
      .then(res => res.json())
      .then(() => {
        setReplyText("")
        if (!post.liked) {
          setPost({...post, likes: parseInt(post.likes || 0) + 1, liked: true})
        } else {
          setPost({...post, likes: parseInt(post.likes) - 1, liked: false})
        }
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  return(
   <div style={{
     marginLeft: `${level === 0 ? 0 : 20}px`,
     borderLeft: level === 0 ? "" : "2px solid #333",
     paddingLeft: "10px"
   }}>
      {clickable ? 
        (
          <Link onClick={() => updatePage(updateKey + 1)} to={`/post/${post.id}`}>
            {post.text}
          </Link>
        ):(
          <span>{post.text}</span>
        )
      }
      <br/>
      <b>@{post.author.username}</b>
      <br/>
      {formattedDate}
      <br/>

      {!loggedIn || (
        <>
          <button onClick={() => setReplyShown(!replyShown)}>reply</button>
          <button
            onClick={likePost}
            style={post.liked ? {
              color: "#ccc",
              backgroundColor: "#3f6ea1"
            } : {}}
          >Like</button>
        </>
      )}
      <p>Likes: {post.likes || 0}</p>
      {!replyShown || (
        <>
          <p>{message}</p>
          <form onSubmit={handleSubmit}>
            <textarea placeholder="Reply to the post" type="text" value={replyText} onChange={(e) => { setReplyText(e.target.value)}} name="text" />
            <input type="submit" value="Submit" />
          </form>
        </>
      )}

      {post.children && post.children.map(child => (
          <Post tokens={tokens} post={child} key={child.id} level={level + 1} loggedIn={loggedIn} clickable={true} />
      ))}
    </div>
  )
}

export default Post
