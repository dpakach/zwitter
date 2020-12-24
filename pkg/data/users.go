package data

var UserStore = new(Users)
var ProfileStore = new(Profiles)

type Profiles struct {
	Profiles []*Profile `json:"profiles"`
}

// Users type for Users List
type Users struct {
	Users []User `json:"users"`
}

// AddDbList adds a User type to Users.Users
func (p *Users) AddDbList(user *User) int64 {
	userId := p.NewID()
	user.ID = userId
	p.Users = append(p.Users, *user)
	return userId
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

// AddDbList adds a User type to Users.Users
func (p *Profiles) AddDbList(profile *Profile) {
	p.Profiles = append(p.Profiles, profile)
}

// ReadFromDb Updates the database Instance from the file
func (p *Profiles) ReadFromDb() *Profiles {
	return p
}

// CommitDb writes the current database instance to the file
func (p *Profiles) CommitDb() {
	return
}

// GetByID returns a user with given ID from the database instance
func (p *Profiles) GetByID(id int64) *Profile {
	for _, item := range p.Profiles {
		if item.UserId == id {
			return item
		}
	}
	return nil
}

// GetByID returns a user with given ID from the database instance
func (p *Profiles) UpdateProfiles(id int64, profile *Profile) *Profile {
	var uId int64 = 0
	for _, item := range p.Profiles {
		uId = item.UserId
		if item.UserId == id {
			item = profile
			break
		}
	}
	return p.GetByID(uId)
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

type Profile struct {
	UserId      int64  `json:"userId"`
	DisplayName string `json:"displayName"`
	Birthday    string `json:"birthday"`
	Gender      int    `json:"gender"`
}

// GetID returns ID of the user
func (p *Profile) GetID() int64 { return p.UserId }

// SaveToStore saves the user to the store
func (p *Profile) SaveToStore() {
	ProfileStore.AddDbList(p)
}
