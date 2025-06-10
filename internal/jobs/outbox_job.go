package jobs

import (
	"context"
	"delivery/internal/adapters/out/postgres/outboxrepo"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
	"time"
)

var _ cron.Job = &OutboxJob{}

type OutboxJob struct {
	outboxRepository outboxrepo.OutboxRepository
	registry         outbox.EventRegistry
	mediatr          ddd.Mediatr
}

func NewOutboxJob(outboxRepository outboxrepo.OutboxRepository,
	registry outbox.EventRegistry,
	mediatr ddd.Mediatr) (*OutboxJob, error) {
	if outboxRepository == nil {
		return nil, errs.NewValueIsRequiredError("outboxRepository")
	}
	if registry == nil {
		return nil, errs.NewValueIsRequiredError("registry")
	}
	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	return &OutboxJob{
		outboxRepository: outboxRepository,
		registry:         registry,
		mediatr:          mediatr}, nil
}

func (j *OutboxJob) Run() {
	ctx := context.Background()

	// Получаем не отправленные Outbox Events
	outboxMessages, err := j.outboxRepository.GetNotPublishedMessages()
	if err != nil {
		log.Error(err)
	}

	// Перебираем в цикле
	for _, outboxMessage := range outboxMessages {
		// Приводим Outbox Message -> Domain Event
		domainEvent, err := j.registry.DecodeDomainEvent(outboxMessage)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(domainEvent)

		err = j.mediatr.Publish(ctx, domainEvent)
		if err != nil {
			log.Error(err)
			continue
		}

		// Если ошибок нет, помечаем Outbox Message как отправленное и сохраняем в БД
		now := time.Now().UTC()
		outboxMessage.ProcessedAtUtc = &now
		err = j.outboxRepository.Update(ctx, outboxMessage)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}
