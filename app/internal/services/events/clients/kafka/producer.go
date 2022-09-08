package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/config"
)

type kafkaProducer struct {
	conn *kafka.Conn
}

func NewKafkaProducer(conf *config.Kafka, log *zap.Logger) (*kafkaProducer, error) {
	log.Debug("creating kafka writer", zap.String("host", conf.Host))
	conn, err := kafka.DialLeader(context.Background(), "tcp", conf.Host, conf.Topic, 0)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{
		conn: conn,
	}, nil
}

func (k *kafkaProducer) Close() error {
	return k.conn.Close()
}

func (k *kafkaProducer) WriteMessage(ctx context.Context, eventName string, data []byte) error {
	_, err := k.conn.WriteMessages(kafka.Message{
		Key:   []byte(eventName),
		Value: data,
	})
	return err
}
