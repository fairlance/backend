package messaging

import (
	"sync"
)

type ProjectRoom struct {
	id           uint
	usersMU      sync.RWMutex
	users        map[*User]bool
	allowedUsers map[*AllowedUser]bool
}

func newProjectRoom(projectID uint, allowedUsers map[*AllowedUser]bool) *ProjectRoom {
	return &ProjectRoom{
		id:           projectID,
		allowedUsers: allowedUsers,
		users:        make(map[*User]bool),
	}
}

func (r *ProjectRoom) getAllowedUser(userType string, userID uint) *AllowedUser {
	for u := range r.allowedUsers {
		if userType == u.fairlanceType && userID == u.fairlanceID {
			return u
		}
	}
	return nil
}

func (r *ProjectRoom) addUser(user *User) {
	r.usersMU.Lock()
	r.users[user] = true
	r.usersMU.Unlock()
}

func (r *ProjectRoom) removeUser(user *User) {
	r.usersMU.Lock()
	delete(r.users, user)
	r.usersMU.Unlock()
}

func (r *ProjectRoom) hasReasonToExist() bool {
	r.usersMU.RLock()
	defer r.usersMU.RUnlock()
	if len(r.users) > 0 {
		return true
	}
	return false
}

func (r *ProjectRoom) isUserAllowed(userType string, userID uint) bool {
	for u := range r.allowedUsers {
		if userType == u.fairlanceType && userID == u.fairlanceID {
			return true
		}
	}
	return false
}

func (r *ProjectRoom) getAbsentUsers() []*AllowedUser {
	var absent []*AllowedUser
	for allowedUser := range r.allowedUsers {
		if !r.isUserConnected(allowedUser) {
			absent = append(absent, allowedUser)
		}
	}
	return absent
}

func (r *ProjectRoom) isUserConnected(allowedUser *AllowedUser) bool {
	for u := range r.users {
		if allowedUser.UniqueID() == u.UniqueID() {
			return true
		}
	}
	return false
}
