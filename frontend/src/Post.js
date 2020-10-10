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

  return(
   <div style={{marginLeft: `${level === 0 ? 0 : 20}px`}}>
      <br/>
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
        <button onClick={() => setReplyShown(!replyShown)}>reply</button>
      )}

      {!replyShown || (
        <>
          <p>{message}</p>
          <form onSubmit={handleSubmit}>
            <textarea placeholder="Reply to the post" type="text" value={replyText} onChange={(e) => { setReplyText(e.target.value)}} name="text" />
            <input type="submit" value="Submit" />
          </form>
        </>
      )}

      {post.children && post.children.map(child => {
          return <Post tokens={tokens} post={child} key={child.id} level={level + 1} loggedIn={loggedIn} clickable={true} />
      })}
      <br/>
    </div>
  )
}

export default Post
