package cmd

import (
	"delivery/internal/adapters/in/jobs"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	services "delivery/internal/core/domain/sevices"
	"delivery/internal/core/ports"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
)

type CompositionRoot struct {
	configs         Config
	gormDb          *gorm.DB
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
	Jobs            Jobs
}

func NewCompositionRoot(_ Config) CompositionRoot {
	app := CompositionRoot{}
	return app
}

func (cr *CompositionRoot) NewDispatchService() services.DispatchService {
	return services.NewDispatchService()
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	createOrderCommandHandler, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWork())
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

type DomainServices struct {
	OrderDispatcher services.DispatchService
}

type Repositories struct {
	UnitOfWork        ports.UnitOfWork
	OrderRepository   ports.OrderRepository
	CourierRepository ports.CourierRepository
}

type CommandHandlers struct {
	AssignOrdersCommandHandler  commands.AssignOrdersCommandHandler
	CreateOrderCommandHandler   commands.CreateOrderCommandHandler
	CreateCourierCommandHandler commands.CreateCourierCommandHandler
	MoveCouriersCommandHandler  commands.MoveCouriersCommandHandler
}

type QueryHandlers struct {
	GetAllCouriersQueryHandler        queries.GetAllCouriersQueryHandler
	GetNotCompletedOrdersQueryHandler queries.GetNotCompletedOrdersQueryHandler
}

type Jobs struct {
	AssignOrdersJob cron.Job
	MoveCouriersJob cron.Job
}

func (cr *CompositionRoot) NewAssignOrdersJob() *jobs.AssignOrdersJob {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("failed to create AssignOrdersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() *jobs.MoveCouriersJob {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("failed to create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) StartCronJobs() {
	c := cron.New(cron.WithSeconds())

	assignJob := cr.NewAssignOrdersJob()
	moveJob := cr.NewMoveCouriersJob()

	if _, err := c.AddJob("@every 1s", assignJob); err != nil {
		log.Fatalf("failed to schedule AssignOrdersJob: %v", err)
	}

	if _, err := c.AddJob("@every 2s", moveJob); err != nil {
		log.Fatalf("failed to schedule MoveCouriersJob: %v", err)
	}

	c.Start()
	log.Println("Cron jobs started")
}

func (cr *CompositionRoot) NewJobs() Jobs {
	assignOrdersJob, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("failed to create assignOrdersJob: %v", err)
	}

	moveCouriersJob, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("failed to create moveCouriersJob: %v", err)
	}

	return Jobs{
		AssignOrdersJob: assignOrdersJob,
		MoveCouriersJob: moveCouriersJob,
	}
}
