package dbmodel

import "time"

type Account struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}
