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
	result := database.Where("username = ?", username).First(&user)

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
	if len(usersId) <= 0 {
		return users, nil
	}

	result := database.Find(&users, usersId)
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

	var count int64
	err := database.Model(&Follow{}).Where("source_user_id = ? and dest_user_id = ?", srcUserId, dstUserId).Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
		// already follows
		return nil, nil
	}

	// not followed yet
	// create follow
	follow = Follow {
		SourceUserID: int(srcUserId),
		DestUserID: int(dstUserId),
	}
	result := database.Create(&follow)
	return &follow, result.Error
}

func DoesFollow(srcUserId, dstUserId uint) (bool, error) {
	follow := Follow{}

	result := database.Where("source_user_id = ? and dest_user_id = ?", srcUserId, dstUserId).Take(&follow)

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
	followed := []Follow{}
	result := database.Where("dest_user_id = ?", target).Find(&followed)
	return followed, result.Error
}

func CreatePost(userId uint, content string) (Post, error) {
	post := Post{}
	post.Content = content
	post.UserID = int(userId)
	err := database.Create(&post).Error
	return post, err
}

func GetPost(id uint) (Post, error) {
	post := Post{}
	post.ID = id
	err := database.Take(&post).Error
	return post, err
}

func UpdatePost(data Post) (Post, error) {
	// can only update content
	data_ := Post{}
	data_.ID = data.ID
	data_.Content = data.Content

	outPost := Post{}
	result := database.Model(&outPost).Updates(&data_)
	return outPost, result.Error
}

func GetRecentPosts(count int) ([]Post, error) {
	posts := make([]Post, count)
	err := database.Limit(count).Order("created_at").Find(&posts).Error
	return posts, err
}