package dbmodel

import "time"

type Reservation struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	ProductId int       `db:"product_id"`
	OrderId   int       `db:"order_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
