package service

import (
	"avito_intership/internal/repo"
	"avito_intership/pkg/broker"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	DepositInput struct {
		UserId int
		Amount float64
	}
	WithdrawInput struct {
		UserId int
		Amount float64
	}
	TransferInput struct {
		From   int
		To     int
		Amount float64
	}
)

type (
	ReservationInput struct {
		UserId    int
		ProductId int
		OrderId   int
		Amount    float64
	}
)

type (
	HistoryInput struct {
		UserId int
		Sort   string
		Offset int
		Limit  int
	}
	HistoryOutput struct {
		OperationId int       `json:"operation_id"`
		ProductId   *int      `json:"product_id"`
		OrderId     *int      `json:"order_id"`
		Amount      float64   `json:"amount"`
		Type        string    `json:"type"`
		CreatedAt   time.Time `json:"created_at"`
	}
)

type Auth interface {
	ValidateToken(token string) bool
	CreateToken() (string, error)
}

type Account interface {
	CreateAccount(ctx context.Context, userId int) error
	GetBalance(ctx context.Context, userId int) (float64, error)

	Deposit(ctx context.Context, input DepositInput) error
	Withdraw(ctx context.Context, input WithdrawInput) error
	Transfer(ctx context.Context, input TransferInput) error
}

type Reservation interface {
	CreateReservation(ctx context.Context, input ReservationInput) (int, error)
	CancelReservation(ctx context.Context, reservationId int) error
	RevenueReservation(ctx context.Context, reservationId int) error
}

type Operation interface {
	GetHistory(ctx context.Context, input HistoryInput) ([]HistoryOutput, error)
	CreateReport(ctx context.Context, year, month int) ([]byte, error)
}

type (
	Services struct {
		Auth        Auth
		Account     Account
		Reservation Reservation
		Operation   Operation
	}
	ServicesDependencies struct {
		Repos      *repo.Repositories
		Producer   broker.Producer
		PrivateKey string
		PublicKey  string
	}
)

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Auth:        newAuthService(d.PrivateKey, d.PublicKey),
		Account:     newAccountService(d.Repos.Account, d.Producer),
		Reservation: newReservationService(d.Repos.Reservation, d.Producer),
		Operation:   newOperationService(d.Repos.Operation),
	}
}

type message struct {
	UserId int
	Amount float64
}

// Функция, которая пушит сообщения в брокер.
// Представим, что у нас есть микросервис нотификаций,
// который отправляет сообщение пользователю о новой операции на аккаунте.
// Тк микросервис вымышленный, то выбрал условный формат сообщения userId + amount.
// Ключ для consumer`а - тип операции
func pushMessage(producer broker.Producer, key string, input any) error {
	body, err := json.Marshal(input)
	if err != nil {
		log.Errorf("/service/service/pushMessage error marshal input: %s", err)
		return err
	}
	_, err = producer.WriteMessages(kafka.Message{
		Key:   []byte(key),
		Value: body,
	})
	if err != nil {
		log.Errorf("/service/service/pushMessage error push message to broker: %s", err)
		return err
	}
	return nil
}
