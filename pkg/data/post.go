package data

var PostStore = new(Posts)

// Posts type for Users List
type Posts struct {
	Posts []Post `json:"posts"`
}

// AddDbList adds a new Post to Posts.Posts
func (p *Posts) AddDbList(post *Post) {
	post.ID = p.NewID()
	p.Posts = append(p.Posts, *post)
}

// ReadFromDb updates Posts.Posts from the file
func (p *Posts) ReadFromDb() *Posts {
	return p
}

// CommitDb writes Posts.Posts to the file
func (p *Posts) CommitDb() {
	return
}

// GetByID returns Post with given ID
func (p *Posts) GetByID(id int64) *Post {
	for _, item := range p.Posts {
		if item.ID == int64(id) {
			return &item
		}
	}
	return nil
}

// NewID returns new ID for the new Post
func (p *Posts) NewID() int64 {
	id := int64(0)
	for _, post := range p.Posts {
		if post.ID > id {
			id = post.ID
		}
	}
	return id + 1
}

// Post type for a Post Instance
type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Created int64  `json:"created"`
	Author  int64  `json:"author"`
}

// GetID returns ID of the Post
func (p *Post) GetID() int64 { return p.ID }

// SetID sets ID of the Post
func (p *Post) SetID(id int64) { p.ID = id }

// SaveToStore saves the Post to the given dB store
func (p *Post) SaveToStore() {
	PostStore.AddDbList(p)
}
