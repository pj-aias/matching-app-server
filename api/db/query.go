package db

func GetUser(id uint64) (User, error) {
	user := User{}
	user.ID = uint(id)

	result := database.Take(&user)

	return user, result.Error
}

func AddUser(user User) (User, error) {
	result := database.Create(&user)
	return user, result.Error
}
