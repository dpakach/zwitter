import * as React from 'react';
import { useState, useEffect } from 'react';
import { get } from './helpers/request';
import Post from './Post';
import { RouteComponentProps } from 'react-router-dom';
import { Tokens } from './types/types';

type TParams = {
  id: string;
};

interface SinglePostParams extends RouteComponentProps<TParams> {
  tokens: Tokens;
  loggedIn: boolean;
}

function SinglePost(props: SinglePostParams) {
  const [message, setMessage]: [string, (message: string) => void] = useState('');
  const [post, setPost] = useState(null);
  const [postId, setPostId] = useState(0);

  useEffect(() => {
    if (postId) {
      get(`/posts/${postId}`, {})
        .then((res) => {
          return res.json();
        })
        .then(
          (json) => {
            setPost(json.post);
          },
          (error) => {
            setMessage(error);
          },
        );
    }
  }, [postId]);

  useEffect(() => {
    setPostId(parseInt(props.match.params.id));
  }, [props.match.params.id]);

  return (
    <>
      <h2>Post</h2>
      {message && <p>{message}</p>}
      {post && <Post post={post} tokens={props.tokens} level={0} {...props} clickable={false} />}
    </>
  );
}

export default SinglePost;
