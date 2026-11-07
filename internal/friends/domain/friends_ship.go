package domain

import "time"

type Friendship struct {
	ID        int       `db:"id"`
	UserIDA   string    `db:"user_id_a"`
	UserIDB   string    `db:"user_id_b"`
	CreatedAt time.Time `db:"created_at"`
}

// NewFriendship упорядочивает пару так, чтобы user_id_a < user_id_b (CHECK в миграции).
// Сравнение строк совпадает с сравнением UUID в Postgres, пока id в каноничном
// lowercase-формате — именно так их отдаёт Postgres.
func NewFriendship(user1, user2 string) *Friendship {
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	return &Friendship{UserIDA: user1, UserIDB: user2}
}

// FriendOf возвращает id второго участника дружбы.
func (f *Friendship) FriendOf(userID string) string {
	if f.UserIDA == userID {
		return f.UserIDB
	}
	return f.UserIDA
}
