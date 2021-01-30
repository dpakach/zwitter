import * as React from 'react';
import { post as httpPost } from './helpers/request';
import { Link, useHistory } from 'react-router-dom';
import { baseUrl } from './const';
import { Tokens, PostType, CreatePostRequest, PostReactTypes } from './types/types';
import { timeSince } from "./helpers/date";

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart, faCommentDots, faRetweet, faReplyAll, faRepublican, faReply } from '@fortawesome/free-solid-svg-icons';

type PostProps = {
  loggedIn: boolean;
  post: PostType;
  tokens: Tokens;
  level: number;
  clickable: boolean;
};

function Post({ post: p, tokens, level, loggedIn, clickable }: PostProps) {
  const [replyShown, setReplyShown]: [boolean, (replyShown: boolean) => void] = React.useState<boolean>(false);
  const [rezweetShown, setRezweetShown]: [boolean, (rezweetShown: boolean) => void] = React.useState<boolean>(false);
  const [replyText, setReplyText]: [string, (replyText: string) => void] = React.useState('');
  const [post, setPost]: [PostType, (post: PostType) => void] = React.useState<PostType>(p);
  const [message, setMessage]: [string, (message: string) => void] = React.useState('');
  const history = useHistory();
  const [timeString, setTimeString]: [string, (timeString: string) => void] = React.useState('');
  const [updateKey, updatePage]: [number, (updateKey: number) => void] = React.useState<number>(0);

  React.useEffect(() => {
    setPost(p);
    let formattedDate: string = timeSince(parseInt(post.created)) + " ago";
    setTimeString(formattedDate); 
  }, []);

  React.useEffect(() => {
    console.log('rerender')
    const timer=setTimeout(() => {
      let formattedDate = timeSince(parseInt(post.created)) + " ago";
 
      setTimeString(formattedDate);
    }, 1000);

    // Clear timeout if the component is unmounted
    return () => clearTimeout(timer);
  })

  const updateTimeAgo = function() {
    const created: Date = new Date(parseInt(post.created) * 1000);
    const formattedDate: string = timeSince(parseInt(post.created)) + " ago"
    setTimeString(formattedDate)
  }

  const dateOptions: Intl.DateTimeFormatOptions = {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  };
  
  function getRezweet(rezweet) {
    if (Object.keys(rezweet).length === 0) {
      return;
    }
    const created: Date = new Date(rezweet.created * 1000);
    const formattedDate: string = created.toLocaleTimeString('en-US', dateOptions);
    return (
      <Link onClick={() => updatePage(updateKey + 1)} to={`/post/${rezweet.id}`}>
        <div className="post-section rezweet-section">
          <p>{rezweet.text}</p>
          {rezweet.media && (
            <img src={`${baseUrl}/media/${rezweet.media}`} alt={post.text} style={{ width: '400px' }} />
          )}
          <Link className="username" to={`/profile/${rezweet.author.username}`}>
            <b>@{rezweet.author.username}</b>
          </Link>
          <p>{formattedDate}</p>
        </div>
      </Link>
    );
  }

  function handleSubmit(e: React.FormEvent): Promise<void> {
    e.preventDefault();

    let body: CreatePostRequest = { text: replyText };
    let url = '/posts';

    if (rezweetShown) {
      url = `/posts/${post.id}/rezweet`;
    } else {
      body = { ...body, parentid: post.id };
    }

    return httpPost(url, { body, headers: { token: tokens.token } })
      .then((res) => res.json())
      .then(
        (json) => {
          setReplyText('');

          if (url === '/posts') {
            setPost({ ...post, children: [json.post, ...(post.children || [])] });
          } else {
            history.push('/post/' + json.post.id);
          }
          setReplyShown(false);
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }

  function likePost(): Promise<void> {
    if (!loggedIn) {
      history.push("/login")
      return
    }
    return httpPost(`/posts/${post.id}/like`, { headers: { token: tokens.token } })
      .then((res) => res.json())
      .then(
        () => {
          setReplyText('');
          if (!post.liked) {
            setPost({ ...post, likes: String((parseInt(post.likes) || 0) + 1), liked: true });
          } else {
            setPost({ ...post, likes: String(parseInt(post.likes) - 1), liked: false });
          }
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }

  function toggleReplyRezweet(type: PostReactTypes): void {
    if (!loggedIn) {
      return history.push("/login")
    }
    if (type == PostReactTypes.REZWEET) {
      setReplyShown(false);
      setRezweetShown(!rezweetShown);
    } else if (type == PostReactTypes.REPLY) {
      setRezweetShown(false);
      setReplyShown(!replyShown);
    } else {
      return;
    }
  }

  return (
    <>
      {post.rezweet.id === '0' ? (
        <></>
      ) : (
        <>
          <small>
            <Link className="username" to={`/profile/${post.author.username}`}>
              @{post.author.username}
            </Link>
            {' rezweeted'}
          </small>
        </>
      )}
      <div
        style={{
          marginLeft: `${level === 0 ? 0 : 20}px`,
          borderLeft: level === 0 ? '' : '2px solid #333',
        }}
        className={clickable && 'post-section'}
      >
        <div className="list-item">
          <Link onClick={() => updatePage(updateKey + 1)} to={`/post/${post.id}`}>
            <div className="list-item__content">{post.text}</div>

            {post.rezweet.id === '0' ? <></> : getRezweet(post.rezweet)}
            <Link className="username" to={`/profile/${post.author.username}`}>
              <b>@{post.author.username}</b>
            </Link>
            <p>{timeString}</p>
            {post.media && (
              <>
                <img src={`${baseUrl}/media/${post.media}`} alt={post.text} style={{ width: '400px' }} />
                <br />
              </>
            )}
          </Link>
          <div className="facts">
            <div className={"fact__icon " + (replyShown ? "fact__icon--1" : "")} onClick={() => toggleReplyRezweet(PostReactTypes.REPLY)}>
              <FontAwesomeIcon className="fact__icon" icon={faCommentDots} />
              <div className="fact__value">13k</div>
            </div>

            <div className={"fact__icon " + (rezweetShown ? "fact__icon--1" : "")} onClick={() => toggleReplyRezweet(PostReactTypes.REZWEET)}>
              <FontAwesomeIcon className="fact__icon" icon={faRetweet} />
              <div className="fact__value">20k</div>
            </div>

            <div className="fact" onClick={likePost}>
              <FontAwesomeIcon className={"fact__icon " + (post.liked ? "fact__icon--1" : "")} icon={faHeart} />
              <div className="fact__value">{post.likes || 0}</div>
            </div>
          </div>
        </div>

        {!(replyShown || rezweetShown) || (
          <>
            <p>{message}</p>
            <form onSubmit={handleSubmit}>
              <textarea
                placeholder={replyShown ? 'Reply to the post' : 'Retweet this post'}
                value={replyText}
                onChange={(e) => {
                  setReplyText(e.target.value);
                }}
                name="text"
              />
              <input type="submit" value="Submit" />
            </form>
          </>
        )}
        {post.children &&
          post.children.map((child) => (
            <Post tokens={tokens} post={child} key={child.id} level={level + 1} loggedIn={loggedIn} clickable={true} />
          ))}
      </div>
    </>
  );
}

export default Post;
