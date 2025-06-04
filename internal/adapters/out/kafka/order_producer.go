package kafka

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/generated/queues/orderstatuschangedpb"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type orderProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewOrderProducer(brokers []string, topic string) (ports.OrderProducer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &orderProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (p *orderProducer) Publish(ctx context.Context, domainEvent ddd.DomainEvent) error {
	completedDomainEvent, ok := domainEvent.(*order.CompletedDomainEvent)
	if !ok {
		return fmt.Errorf("unexpected domain event type: %T", domainEvent)
	}

	integrationEvent, err := p.mapDomainEventToIntegrationEvent(completedDomainEvent)
	if err != nil {
		return fmt.Errorf("map event: %w", err)
	}

	bytes, err := json.Marshal(integrationEvent)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(completedDomainEvent.OrderID.String()),
		Value: sarama.ByteEncoder(bytes),
	}

	resultCh := make(chan error, 1)

	go func() {
		_, _, err := p.producer.SendMessage(msg)
		resultCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}

func (p *orderProducer) Close() error {
	return p.producer.Close()
}

func (p *orderProducer) mapDomainEventToIntegrationEvent(completedDomainEvent *order.CompletedDomainEvent) (*orderstatuschangedpb.OrderStatusChangedIntegrationEvent, error) {
	status, ok := orderstatuschangedpb.OrderStatus_value[completedDomainEvent.OrderStatus]
	if !ok {
		return nil, errs.NewValueIsInvalidError("OrderStatus")
	}

	integrationEvent := orderstatuschangedpb.OrderStatusChangedIntegrationEvent{
		OrderId:     completedDomainEvent.OrderID.String(),
		OrderStatus: orderstatuschangedpb.OrderStatus(status),
	}
	return &integrationEvent, nil
}
