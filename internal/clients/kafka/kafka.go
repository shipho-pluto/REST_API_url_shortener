package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"

	"github.com/segmentio/kafka-go"
)

type Broker struct {
	log    *slog.Logger
	conn   *kafka.Conn
	Reader *kafka.Reader
	Writer *kafka.Writer
}

func New(log *slog.Logger, cfg config.Broker) (*Broker, error) {
	const op = "kafka.New"

	conn, err := kafka.Dial(cfg.Network, cfg.Address)
	if err != nil {
		return &Broker{}, fmt.Errorf("%s: %w", op, err)
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             cfg.TopicName,
			NumPartitions:     cfg.Partitions,
			ReplicationFactor: cfg.Replications,
		},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Address},
		Topic:   cfg.TopicName,
		GroupID: cfg.GroupID,
	})

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.Address},
		Topic:   cfg.TopicName,
	})

	if err := conn.CreateTopics(topicConfigs...); err != nil {
		return &Broker{}, fmt.Errorf("%s: %w", op, err)
	}

	return &Broker{
		conn:   conn,
		log:    log,
		Reader: r,
		Writer: w,
	}, nil
}

func (b *Broker) Stop() error {
	return b.conn.Close()
}

func (b *Broker) ConsumeMessage(ctx context.Context) {
	const op = "kafka.ConsumeMessage"
	for {
		msg, err := b.Reader.ReadMessage(ctx)
		if err != nil {
			b.log.Error("error with read message", sl.Err(err))
		}

		b.log.Info("[CATCHED MESSAGE]",
			slog.String("key", string(msg.Key)),
			slog.String("value", string(msg.Value)),
			slog.Int64("offset", msg.Offset),
		)
	}
}

type MessageToBroker struct {
	Key   string
	Value string
}

func (b *Broker) ProduceMessage(ctx context.Context, msg MessageToBroker) error {
	const op = "kafka.ProduceMessage"

	kafkaMsg := kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%s", msg.Key)),
		Value: []byte(fmt.Sprintf("Mwssage-%s", msg.Value)),
	}

	err := b.Writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		b.log.Error("error with sending message", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	b.log.Info("sent message successfully")
	return nil
}
