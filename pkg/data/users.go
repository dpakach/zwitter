package data

var UserStore = new(Users)

// Users type for Users List
type Users struct {
	Users []User `json:"users"`
}

// AddDbList adds a User type to Users.Users
func (p *Users) AddDbList(user *User) {
	user.ID = p.NewID()
	p.Users = append(p.Users, *user)
}

// ReadFromDb Updates the database Instance from the file
func (p *Users) ReadFromDb() *Users {
	return p
}

// CommitDb writes the current database instance to the file
func (p *Users) CommitDb() {
	return
}

// GetByID returns a user with given ID from the database instance
func (p *Users) GetByID(id int64) *User {
	for _, item := range p.Users {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

// NewID returns a new ID for creating new database object
func (p *Users) NewID() int64 {
	id := int64(0)
	for _, user := range p.Users {
		if user.ID > id {
			id = user.ID
		}
	}
	return id + 1
}

// GetByID returns a user with given ID from the database instance
func (p *Users) GetByUsername(username string) *User {
	p.ReadFromDb()
	for _, item := range p.Users {
		if item.Username == username {
			return &item
		}
	}
	return nil
}

// User type for a single user instance
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Created  int64  `json:"created"`
}

// GetID returns ID of the user
func (p *User) GetID() int64 { return p.ID }

// SetID sets ID of the user
func (p *User) SetID(id int64) { p.ID = id }

// SaveToStore saves the user to the store
func (p *User) SaveToStore() {
	UserStore.AddDbList(p)
}
