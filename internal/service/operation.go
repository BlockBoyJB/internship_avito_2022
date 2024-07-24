package service

import (
	"avito_intership/internal/repo"
	"avito_intership/internal/repo/pgerrs"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const operationPrefixLog = "/service/operation"

const (
	defaultLimit = 20
)

type operationService struct {
	operation repo.Operation
}

func newOperationService(operation repo.Operation) *operationService {
	return &operationService{operation: operation}
}

func (s *operationService) GetHistory(ctx context.Context, input HistoryInput) ([]HistoryOutput, error) {
	if input.Limit <= 0 || input.Limit > defaultLimit {
		input.Limit = defaultLimit
	}
	switch input.Sort { // TODO более умную сортировку
	case "amount":
		input.Sort = "amount DESC"
	case "type":
		input.Sort = "type DESC"
	default:
		input.Sort = "created_at DESC"
	}
	history, err := s.operation.GetHistory(ctx, input.UserId, input.Sort, input.Offset, input.Limit)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return nil, ErrAccountNotFound
		}
		log.Errorf("%s/GetHistory error get account operation history: %s", operationPrefixLog, err)
		return nil, err
	}
	var result []HistoryOutput
	for _, o := range history {
		result = append(result, HistoryOutput{
			OperationId: o.Id,
			ProductId:   o.ProductId,
			OrderId:     o.OrderId,
			Amount:      o.Amount,
			Type:        o.Type,
			CreatedAt:   o.CreatedAt,
		})
	}
	return result, nil
}

func (s *operationService) CreateReport(ctx context.Context, year, month int) ([]byte, error) {
	const reportLine = "%d;%f"
	group, err := s.operation.GroupProductRevenue(ctx, year, month)
	if err != nil {
		return nil, err
	}

	result := &bytes.Buffer{}
	w := csv.NewWriter(result)

	for productId, amount := range group {
		// Не понятно, можно ли продолжать или стоит сразу ошибку и выход. Решил делать возврат сразу после ошибки,
		// потому что тут собирается отчет для налоговой, следовательно, ошибки или пропуски тут недопустимы
		if err = w.Write([]string{fmt.Sprintf(reportLine, productId, amount)}); err != nil {
			log.Errorf("%s/CreateReport error write line: %s", operationPrefixLog, err)
			return nil, err
		}
	}
	w.Flush()
	if err = w.Error(); err != nil {
		log.Errorf("%s/CreateReport error write buffer: %s", operationPrefixLog, err)
		return nil, err
	}

	return result.Bytes(), nil
}
