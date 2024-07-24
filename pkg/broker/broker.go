package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	defaultWriteTopic  = "account-balance"
	defaultConnTimeout = time.Second * 10
)

type Producer interface {
	WriteMessages(msgs ...kafka.Message) (int, error)
	Close()
}

type producer struct {
	*kafka.Conn
}

func NewProducer(url string) (Producer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	conn, err := kafka.DialLeader(ctx, "tcp", url, defaultWriteTopic, 0)
	if err != nil {
		return nil, err
	}
	return &producer{conn}, nil
}

func (p *producer) Close() {
	_ = p.Conn.Close()
}
