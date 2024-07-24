package service

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo"
	"avito_intership/internal/repo/pgerrs"
	"avito_intership/pkg/broker"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
)

const (
	accountServicePrefixLog = "/service/account"
)

type accountService struct {
	account  repo.Account
	producer broker.Producer
}

func newAccountService(account repo.Account, producer broker.Producer) *accountService {
	return &accountService{
		account:  account,
		producer: producer,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, userId int) error {
	if err := s.account.CreateAccount(ctx, userId); err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return ErrAccountAlreadyExists
		}
		log.Errorf("%s/CreateAccount error create account: %s", accountServicePrefixLog, err)
		return err
	}
	return nil
}

func (s *accountService) GetBalance(ctx context.Context, userId int) (float64, error) {
	balance, err := s.account.GetBalance(ctx, userId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return 0, ErrAccountNotFound
		}
		log.Errorf("%s/GetBalance error get balance: %s", accountServicePrefixLog, err)
		return 0, err
	}
	return balance, nil
}

func (s *accountService) Deposit(ctx context.Context, input DepositInput) error {
	if err := s.account.Deposit(ctx, input.UserId, input.Amount); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrAccountNotFound
		}
		log.Errorf("%s/Deposit error update account balance: %s", accountServicePrefixLog, err)
		return err
	}

	if err := pushMessage(s.producer, dbmodel.OperationDeposit, input); err != nil {
		return err
	}
	return nil
}

func (s *accountService) Withdraw(ctx context.Context, input WithdrawInput) error {
	if err := s.account.Withdraw(ctx, input.UserId, input.Amount); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrAccountNotFound
		}
		if errors.Is(err, pgerrs.ErrNotEnoughBalance) {
			return ErrNotEnoughBalance
		}
		log.Errorf("%s/Withdraw error update account balance: %s", accountServicePrefixLog, err)
		return ErrCannotUpdateBalance
	}

	if err := pushMessage(s.producer, dbmodel.OperationWithdraw, input); err != nil {
		return err
	}
	return nil
}

func (s *accountService) Transfer(ctx context.Context, input TransferInput) error {
	if err := s.account.Transfer(ctx, input.From, input.To, input.Amount); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrAccountNotFound
		}
		if errors.Is(err, pgerrs.ErrNotEnoughBalance) {
			return ErrNotEnoughBalance
		}
		log.Errorf("%s/Transfer error transfer: %s", accountServicePrefixLog, err)
		return ErrCannotUpdateBalance
	}

	// очевидно, что для того кто отправил перевод уведомление не нужно
	err := pushMessage(s.producer, dbmodel.OperationIncomingTransfer, message{
		UserId: input.To,
		Amount: input.Amount,
	})
	if err != nil {
		return err
	}
	return nil
}
