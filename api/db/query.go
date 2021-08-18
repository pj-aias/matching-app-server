package db

import (
	"errors"

	"gorm.io/gorm"
)

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

func LookupUser(username string) (User, error) {
	user := User{}
	user.Username = username
	result := database.Take(&user)

	return user, result.Error
}

func UpdateUser(userId uint, user User) (User, error) {
	outUser := User{}
	outUser.ID = userId
	result := database.Model(&outUser).Updates(&user)
	return outUser, result.Error
}

func GetUsers(usersId []uint) ([]User, error) {
	users := make([]User, len(usersId))
	for i := range users {
		users[i].ID = usersId[i]
	}

	result := database.Model(&User{}).Find(&users)
	return users, result.Error
}

func GetPasswordHash(userId uint64) (PasswordHash, error) {
	hash := PasswordHash{UserID: int(userId)}
	result := database.Take(&hash)

	return hash, result.Error
}

func AddPasswordHash(userId uint64, hash []byte) (PasswordHash, error) {
	hashData := PasswordHash{
		UserID: int(userId),
		Hash:   hash,
	}

	result := database.Create(&hashData)
	return hashData, result.Error
}

func CreateFollow(srcUserId, dstUserId uint) (*Follow, error) {
	follow := Follow {}

	err := database.Where("source_user_id = ? and dest_user_id = ?", srcUserId, dstUserId).Find(&follow).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// not followed yet
		// create follow
		follow = Follow {
			SourceUserID: int(srcUserId),
			DestUserID: int(dstUserId),
		}
		result := database.Create(&follow)
		return &follow, result.Error
	} else if err == nil {
		// already follows
		return nil, nil
	} else {
		// some error occured
		return nil, err
	}
}

func DoesFollow(srcUserId, dstUserId uint) (bool, error) {
	follow := Follow {
		SourceUserID: int(srcUserId),
		DestUserID: int(dstUserId),
	}

	result := database.Take(&follow)

	if err := result.Error; err == nil {
		// following
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// not following
		return false, nil
	} else {
		// an error occured
		return false, err
	}
}

func DestroyFollow(srcUserId, dstUserId uint) (error) {
	result := database.Where("source_user_id = ? and dest_user_id = ?", srcUserId, dstUserId).Delete(&Follow{})
	if err := result.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// not following
		return nil
	} else {
		return err
	}
}

func GetFollowing(source uint) ([]Follow, error) {
	following := make([]Follow, 16)
	result := database.Where("source_user_id = ?", source).Find(&following)
	return following, result.Error
}

func GetFollowed(target uint) ([]Follow, error) {
	followed := make([]Follow, 16)
	result := database.Where("dest_user_id = ?", target).Find(&followed)
	return followed, result.Error
}