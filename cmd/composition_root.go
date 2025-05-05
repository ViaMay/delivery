package cmd

import services "delivery/internal/core/domain/sevices"

type CompositionRoot struct {
}

func NewCompositionRoot(_ Config) CompositionRoot {
	app := CompositionRoot{}
	return app
}

func (cr *CompositionRoot) NewDispatchService() services.DispatchService {
	return services.NewDispatchService()
}
