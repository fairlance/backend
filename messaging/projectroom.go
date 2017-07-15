package messaging

type ProjectRoom struct {
	id      uint
	users   map[User]bool
	project *Project
}

// func (r *ProjectRoom) getUser(userType string, userID uint) *User {
// 	for user := range r.Users {
// 		if user.userType == userType && user.ID == userID {
// 			return &user
// 		}
// 	}
// 	return nil
// }

// func (r *ProjectRoom) HasReasonToExist() bool {
// 	for _, user := range r.Users {
// 		if user.online {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (r *ProjectRoom) StartReadWrite(conn *userConn) (*User, error) {
// 	user, ok := r.Users[conn.id]
// 	if ok {
// 		user.Activate(conn)
// 		go user.startWriting()
// 		go user.startReading()
// 		return user, nil
// 	}

// 	return nil, fmt.Errorf("user %s not found", conn.id)
// }

// func (r *ProjectRoom) Close() {
// 	for _, user := range r.Users {
// 		user.Close()
// 	}
// }

// func (r *ProjectRoom) HasUser(user *models.User) bool {
// 	for _, u := range r.Users {
// 		if user.ID == u.id && user.Type == u.userType {
// 			return true
// 		}
// 	}
// 	return false
// }
