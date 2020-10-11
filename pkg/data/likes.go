package data

import (
	"errors"
)

var LikeStore = new(Likes)

// likes type for Users List
type Likes struct {
	Likes []Like `json:"like"`
}

// AddDbList adds a new like to likes.likes
func (l *Likes) AddDbList(like *Like) {
	l.Likes = append(l.Likes, *like)
}

// ReadFromDb updates likes.likes from the file
func (l *Likes) ReadFromDb() *Likes {
	return l
}

// CommitDb writes likes.likes to the file
func (l *Likes) CommitDb() {
	return
}

func (l *Likes) GetLikesForPost(postid int64) []Like {
	likes := []Like{}
	for _, item := range l.Likes {
		if item.Post == postid {
			likes = append(likes, item)
		}
	}
	return likes
}

func (l *Likes) GetlikeByAuthor(id int64) []Like {
	likes := []Like{}
	for _, item := range l.Likes {
		if item.Author == id {
			likes = append(likes, item)
		}
	}
	return likes
}

func (l *Likes) GetByPostAndAuthor(postId, authorId int64) *Like {
	for _, item := range l.Likes {
		if item.Author == authorId && item.Post == postId {
			return &item
		}
	}
	return nil
}

func (l *Likes) DeleteLike(postId, authorId int64) error {
	found := false
	for idx, item := range l.Likes {
		if item.Author == authorId && item.Post == postId {
			l.Likes = append(l.Likes[:idx], l.Likes[idx+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("Could not find the item")
	}
	return nil
}

func (l *Likes) GetLikesCount(postId int64) int64 {
	return int64(len(l.GetLikesForPost(postId)))
}

// like type for a like Instance
type Like struct {
	Post   int64 `json:"post"`
	Author int64 `json:"author"`
}

// SaveToStore saves the like to the given dB store
func (l *Like) SaveToStore() {
	LikeStore.AddDbList(l)
}
