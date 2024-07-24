package dbmodel

import "time"

// Operations
const (
	OperationDeposit  = "deposit"  // Пополнение
	OperationWithdraw = "withdraw" // Снятие

	OperationOutgoingTransfer = "outgoing-transfer" // Исходящий перевод
	OperationIncomingTransfer = "incoming-transfer" // Входящий перевод

	OperationReservation   = "reservation"    // Резервация денег (удержание)
	OperationDereservation = "de-reservation" // Дерезервация денег (возврат)
	OperationRevenue       = "revenue"        // Признание выручки
)

type Operation struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	ProductId *int      `db:"product_id"` // pointer because value in db can be null
	OrderId   *int      `db:"order_id"`   // pointer because value in db can be null
	Amount    float64   `db:"amount"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
}
