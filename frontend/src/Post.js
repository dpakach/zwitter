import React, {useState, useEffect} from "react"
import {post as httpPost, get} from "./helpers/request"
import {Link, useHistory} from "react-router-dom"
import {baseUrl} from "./const"

function Post({post: p, tokens, level, loggedIn, clickable}) {
  const reactTypes = Object.freeze({
    REPLY: 'reply',
    REZWEET: 'rezweet',
    LIKE: 'LIKE',
  })
  const [replyShown, setReplyShown] = useState(false)
  const [rezweetShown, setRezweetShown] = useState(false)
  const [replyText, setReplyText] = useState("")
  const [post, setPost] = useState(p)
  const [message, setMessage] = useState("")
  const history = useHistory()

  const [updateKey, updatePage] = useState(0)

  useEffect(() => {
    setPost(p)
  }, [])

  const dateOptions = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' };
  const created = new Date(post.created * 1000)
  const formattedDate = created.toLocaleTimeString("en-US", dateOptions)

  function getRezweet(rezweet) {
    if (Object.keys(rezweet).length === 0) {
      return
    }
    const created = new Date(rezweet.created * 1000)
    const formattedDate = created.toLocaleTimeString("en-US", dateOptions)
    return (
       <div style={{
         margin: "10px",
         border: "2px solid #333",
         padding: "10px",
         maxWidth: "400px",
       }}>
          <Link onClick={() => updatePage(updateKey + 1)} to={`/post/${rezweet.id}`}>
            <p>{rezweet.text}</p>
          </Link>
          {rezweet.media && (
            <img src={`${baseUrl}/media/${rezweet.media}`} alt={post.text} style={{width: '400px'}} />
          )}
          <b>@{rezweet.author.username}</b>
          <br/>
          {formattedDate}
          <br/>
        </div>
    )
  }

  function handleSubmit(e) {
    e.preventDefault()
    let body = {text: replyText}
    let url = "/posts"

    if (rezweetShown) {
      url = `/posts/${post.id}/rezweet`
    } else {
      body = {...body, parentid: post.id}
    }

    return httpPost(url, {body, headers: {"token": tokens.token}})
      .then(res => res.json())
      .then((json) => {
        setReplyText("")

        if (url === "/posts") {
          setPost({...post, children: [json.post, ...(post.children || [])]})
        } else {
          history.push("/post/" + json.post.id)
        }
        setReplyShown(false)
      }, (error) => {
        setMessage("Error: " + error.message)
      })
  }

  function likePost() {
    return httpPost(`/posts/${post.id}/like`, {headers: {"token": tokens.token}})
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

  function toggleReplyRezweet(type) {
    if (type == reactTypes.REZWEET) {
      setReplyShown(false)
      setRezweetShown(!rezweetShown)
    } else if (type == reactTypes.REPLY) {
      setRezweetShown(false)
      setReplyShown(!replyShown)
    } else {
      return
    }
  }

  return(
   <div style={{
     marginLeft: `${level === 0 ? 0 : 20}px`,
     borderLeft: level === 0 ? "" : "2px solid #333",
     paddingLeft: "10px",
     marginBottom: "5rem",
   }}>
      {Object.keys(post.rezweet).length === 0 ?
        <></> :
        <>
          <small>@{post.author.username} rezweeted</small>
          <br/>
        </>
      }
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
      {(post.rezweet === {}) ?
        <></> :
        getRezweet(post.rezweet)
      }
      <b>@{post.author.username}</b>
      <br/>
      {formattedDate}
      <br/>
      {post.media && (
        <img src={`${baseUrl}/media/${post.media}`} alt={post.text} style={{width: '400px'}} />
      )}
      <br/>
      {!loggedIn || (
        <>
          <button onClick={() => toggleReplyRezweet(reactTypes.REPLY)}>reply</button>
          <button onClick={() => toggleReplyRezweet(reactTypes.REZWEET)}>rezweet</button>
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
      {!(replyShown || rezweetShown) || (
        <>
          <p>{message}</p>
          <form onSubmit={handleSubmit}>
            <textarea placeholder={
              replyShown ? "Reply to the post" : "Retweet this post"
            } type="text" value={replyText} onChange={(e) => { setReplyText(e.target.value)}} name="text" />
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
