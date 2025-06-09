package cmd

import (
	kafkain "delivery/internal/adapters/in/kafka"
	grpcout "delivery/internal/adapters/out/grpc/geo"
	kafkaout "delivery/internal/adapters/out/kafka"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/eventhandlers"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/sevices"
	"delivery/internal/core/ports"
	"delivery/internal/jobs"
	"delivery/internal/pkg/ddd"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"sync"
)

type CompositionRoot struct {
	configs   Config
	gormDb    *gorm.DB
	geoClient ports.GeoClient

	closers []Closer
	onceGeo sync.Once
}

func NewCompositionRoot(configs Config, gormDb *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDb:  gormDb,
	}
}

func (cr *CompositionRoot) NewDispatchService() services.DispatchService {
	return services.NewDispatchService()
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb, cr.NewMediatrWithSubscriptions())
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	createOrderCommandHandler, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWork(), cr.NewGeoClient())
	if err != nil {
		log.Fatalf("cannot create CreateOrderCommandHandler: %v", err)
	}
	return createOrderCommandHandler
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	createCourierCommandHandler, err := commands.NewCreateCourierCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create CreateCourierCommandHandler: %v", err)
	}
	return createCourierCommandHandler
}

func (cr *CompositionRoot) NewAssignOrdersCommandHandler() commands.AssignOrdersCommandHandler {
	assignOrdersCommandHandler, err := commands.NewAssignOrdersCommandHandler(
		cr.NewUnitOfWork(), cr.NewDispatchService())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersCommandHandler: %v", err)
	}
	return assignOrdersCommandHandler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	moveCouriersCommandHandler, err := commands.NewMoveCouriersCommandHandler(
		cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersCommandHandler: %v", err)
	}
	return moveCouriersCommandHandler
}

func (cr *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	getAllCouriersQueryHandler, err := queries.NewGetAllCouriersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return getAllCouriersQueryHandler
}

func (cr *CompositionRoot) NewGetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	getNotCompletedOrdersQueryHandler, err := queries.NewGetNotCompletedOrdersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetNotCompletedOrdersQueryHandler: %v", err)
	}
	return getNotCompletedOrdersQueryHandler
}

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := grpcout.NewClient(cr.configs.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}

func (cr *CompositionRoot) NewBasketConfirmedConsumer() kafkain.BasketConfirmedConsumer {
	consumer, err := kafkain.NewBasketConfirmedConsumer(
		[]string{cr.configs.KafkaHost},
		cr.configs.KafkaConsumerGroup,
		cr.configs.KafkaBasketConfirmedTopic,
		cr.NewCreateOrderCommandHandler(),
	)
	if err != nil {
		log.Fatalf("cannot create BasketConfirmedConsumer: %v", err)
	}
	cr.RegisterCloser(consumer)
	return consumer
}

func (cr *CompositionRoot) NewOrderCompletedDomainEventHandler() ddd.EventHandler {
	producer := cr.NewOrderProducer()
	handler, err := eventhandlers.NewOrderCompletedDomainEventHandler(producer)
	if err != nil {
		log.Fatalf("cannot create OrderCompletedDomainEventHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewMediatrWithSubscriptions() ddd.Mediatr {
	mediatr := ddd.NewMediatr()
	mediatr.Subscribe(cr.NewOrderCompletedDomainEventHandler(), order.NewEmptyCompletedDomainEvent())
	return mediatr
}

func (cr *CompositionRoot) NewOrderProducer() ports.OrderProducer {
	producer, err := kafkaout.NewOrderProducer([]string{cr.configs.KafkaHost}, cr.configs.KafkaOrderChangedTopic)
	if err != nil {
		log.Fatalf("cannot create OrderProducer: %v", err)
	}
	cr.RegisterCloser(producer)
	return producer
}
