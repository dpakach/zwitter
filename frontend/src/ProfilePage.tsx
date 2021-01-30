import * as React from 'react';
import { useState, useEffect } from 'react';
import { Redirect, useParams } from 'react-router-dom';
import { get, post } from './helpers/request';
import { getDate } from './helpers/date';
import { User, Tokens, Genders } from './types/types';
import Posts from './Posts';

type ProfileProps = {
  tokens: Tokens;
  loggedIn: boolean;
  self?: boolean;
};

function ProfilePage(props: ProfileProps) {
  const { tokens, loggedIn, self } = props;

  function getGenderName(id: number) {
    for (const enumMember in Genders) {
      const isValueProperty = parseInt(enumMember, 10) >= 0;
      if (isValueProperty) {
        return Genders[enumMember].toLowerCase();
      }
    }
    return 'not specified';
  }

  const params = useParams<Record<string, string | undefined>>();

  const [profileUser, setProfileUser] = useState(params.user);
  const [message, setMessage]: [string, (prop: string) => void] = useState('');
  const [displayName, setDisplayName]: [string, (prop: string) => void] = useState('');
  const [gender, setGender] = useState(Genders.NOT_SPECIFIED);
  const [birthday, setBirthday]: [string, (prop: string) => void] = useState('');
  const [user, setUser] = useState({} as User);
  const [following, setFollowing] = useState(false);

  useEffect(() => {
    let url;
    if (profileUser) {
      url = `/users/profile/${profileUser}`;
    } else {
      url = '/users/self/profile';
    }
    const headers = loggedIn ? { token: tokens.token } : {};
    get(url, { headers })
      .then((res) => res.json())
      .then((data) => {
        setDisplayName(data.Profile['displayName'] || '');
        setBirthday(data.Profile['dateOfBirth'] || '');
        setGender(data.Profile['gender'] || Genders.NOT_SPECIFIED);
        setUser(data.user);
        setFollowing(data.following);
      });
  }, []);

  function handleSubmit(e) {
    if (!loggedIn) {
      return;
    }
    e.preventDefault();
    return post('/users/profile', {
      body: { profile: { userId: user.id, displayName, gender: gender, dateOfBirth: birthday } },
      headers: { token: tokens.token },
    })
      .then((res) => res.json())
      .then(
        () => {
          setMessage('Success: updated User Profile');
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }

  function toggleFollow(e) {
    if (!loggedIn) {
      return;
    }
    e.preventDefault();
    const urlPrefix = following ? 'unfollow' : 'follow';
    return post(`/users/${urlPrefix}/${user.username}`, {
      headers: { token: tokens.token },
      body: { username: user.username },
    })
      .then((res) => res.json())
      .then(
        () => {
          setMessage(`Success: ${urlPrefix}ed the user`);
          setFollowing(!following);
        },
        (error) => {
          setMessage('Error: ' + error.message);
        },
      );
  }

  return (
    <>
      {self && !loggedIn && <Redirect to="/" />}
      <div>
        {displayName || 'Name not set'} ({getGenderName(gender) || 'Gender not specified'})
        <br />@{user.username}{' '}
        {loggedIn && profileUser && <a onClick={toggleFollow}>{following ? 'Unfollow' : 'follow'}</a>}
        <br />
        joined {getDate(user.created)}
        <br />
        <br />
        Birthday {birthday || 'Birthday not set'}
      </div>
      {self && (
        <>
          <h2>Edit Profile</h2>
          <form onSubmit={handleSubmit}>
            <label>
              DisplayName:
              <input
                type="text"
                value={displayName}
                onChange={(e) => {
                  setDisplayName(e.target.value);
                }}
                name="displayName"
              />
            </label>
            <label htmlFor="gender">Gender:</label>
            <select
              id="gender"
              value={gender}
              onChange={(e) => {
                setGender(parseInt(e.target.value));
              }}
              name="gender"
            >
              <option value={Genders.NOT_SPECIFIED}>Not Specified</option>
              <option value={Genders.MALE}>Male</option>
              <option value={Genders.FEMALE}>Female</option>
              <option value={Genders.OTHER}>Other</option>
            </select>
            <label>
              Birthday:
              <input
                type="date"
                value={birthday}
                onChange={(e) => {
                  setBirthday(e.target.value);
                }}
                name="displayName"
              />
            </label>
            <input type="submit" value="Submit" />
          </form>
        </>
      )}
      <h3>Posts by @{user.username}</h3>
      <Posts loggedIn={loggedIn} tokens={tokens} showInput={false} />
    </>
  );
}

export default ProfilePage;
