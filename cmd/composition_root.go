package cmd

import (
	"delivery/internal/adapters/out/postgres"
	services "delivery/internal/core/domain/sevices"
	"delivery/internal/core/ports"
	"gorm.io/gorm"
	"log"
)

type CompositionRoot struct {
	configs Config
	gormDb  *gorm.DB
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
