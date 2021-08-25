package db_test

import (
	"testing"

	"github.com/pj-aias/matching-app-server/db"
)

func TestRoomContainsUser(t *testing.T) {
	user1 := sampleUser(1, "hoge")
	user2 := sampleUser(2, "fuga")

	room := db.Chatroom{
		Users: []db.User{user1, user2},
	}

	tests := []uint{1, 3}
	wants := []bool{true, false}

	for i, tt := range tests {
		got, want := room.ContainsUser(tt), wants[i]
		if got != want {
			t.Errorf("Room.ContainsUSer() = %v, want %v", got, want)
		}
	}
}
