package db

func (room Chatroom) ContainsUser(userId uint) bool {
	for _, u := range room.Users {
		if u.ID == userId {
			return true
		}
	}

	return false
}
