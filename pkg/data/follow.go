package data

import (
	"errors"
)

var FollowStore = new(Follows)

// follows type for Users List
type Follows struct {
	Follows []Follow `json:"follow"`
}

// AddDbList adds a new follow to follows.follows
func (f *Follows) AddDbList(follow *Follow) {
	f.Follows = append(f.Follows, *follow)
}

// ReadFromDb updates follows.follows from the file
func (f *Follows) ReadFromDb() *Follows {
	return f
}

// CommitDb writes follows.follows to the file
func (f *Follows) CommitDb() {
	return
}

// followers of id
func (f *Follows) GetFollowers(id int64) []Follow {
	follows := []Follow{}
	for _, item := range f.Follows {
		if item.Follows == id {
			follows = append(follows, item)
		}
	}
	return follows
}

// id is following
func (f *Follows) GetfollowByUser(id int64) []Follow {
	follows := []Follow{}
	for _, item := range f.Follows {
		if item.User == id {
			follows = append(follows, item)
		}
	}
	return follows
}

// get follow relation between 2 users
func (f *Follows) GetByUserAndFollows(userId, follows int64) *Follow {
	for _, item := range f.Follows {
		if item.User == userId && item.Follows == follows {
			return &item
		}
	}
	return nil
}

func (f *Follows) DeleteFollow(userId, follows int64) error {
	found := false
	for idx, item := range f.Follows {
		if item.User == userId && item.Follows == follows {
			f.Follows = append(f.Follows[:idx], f.Follows[idx+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("Could not find the item")
	}
	return nil
}

func (f *Follows) GetFollowsCount(userId int64) int64 {
	return int64(len(f.GetFollowers(userId)))
}

func (f *Follows) GetFollowingCount(userId int64) int64 {
	return int64(len(f.GetfollowByUser(userId)))
}

// follow type for a follow Instance
type Follow struct {
	User    int64 `json:"user"`
	Follows int64 `json:"follows"`
}

// SaveToStore saves the follow to the given dB store
func (f *Follow) SaveToStore() {
	FollowStore.AddDbList(f)
}
