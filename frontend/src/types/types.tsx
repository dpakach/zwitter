type User = {
  id: number;
  username: string;
  created: number;
};

type Tokens = {
  token: string;
  refreshToken: string;
  user: User;
};

type PostType = {
  id: string;
  created: string;
  text: string;
  children: Array<PostType>;
  liked: boolean;
  likes: string;
  rezweet?: PostType;
  author: User;
  media: string;
};

type UserRequest = {
  username: string;
  password: string;
};

type CreatePostRequest = {
  text: string;
  parentid?: string;
  media?: string;
};

enum Genders {
  NOT_SPECIFIED = 0,
  MALE = 1,
  FEMALE = 2,
  OTHER = 3,
}

enum PostReactTypes {
  REPLY = 0,
  REZWEET = 1,
  LIKE = 2,
}

export { User, Tokens, Genders, PostType, CreatePostRequest, UserRequest, PostReactTypes };
