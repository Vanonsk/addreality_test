package withdb

import (
	"addreality_t/internal/resources"
)

type WithDB struct {
	db *resources.DB
}

func New(db *resources.DB) *WithDB {
	return &WithDB{
		db: db,
	}
}

// SaveRateUser, save new user information in storage
func (wdb *WithDB) SaveRateUser(user *User) {
	wdb.db.Log.Debugf("Save user %s with counter %d", user.ID, user.Counter.Rate())
	wdb.db.Store(user.ID, user.Counter)
}

// GetUser, get user from storage by userId
func (wdb *WithDB) GetUser(userId string) (*User, error) {
	counter, ok := wdb.db.Load(userId)
	if ok {
		return &User{userId, counter}, nil
	}
	return nil, nil
}

// GetAllUsers, get all users from storage
func (wdb *WithDB) GetAllUsers() ([]*User, error) {
	robotUsers := make([]*User, 0)
	for _, k := range wdb.db.GetAllKeys() {
		counter, ok := wdb.db.Load(k)
		if ok {
			robotUsers = append(robotUsers, &User{ID: k, Counter: counter})
		}
	}
	wdb.db.Log.Debug(robotUsers)
	return robotUsers, nil
}
