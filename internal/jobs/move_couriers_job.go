package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriersJob{}

type MoveCouriersJob struct {
	moveCouriersCommandHandler commands.MoveCouriersCommandHandler
}

func NewMoveCouriersJob(
	moveCouriersCommandHandler commands.MoveCouriersCommandHandler,
) (*MoveCouriersJob, error) {
	if moveCouriersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("MoveCouriersCommandHandler")
	}
	return &MoveCouriersJob{
		moveCouriersCommandHandler: moveCouriersCommandHandler,
	}, nil
}

func (j *MoveCouriersJob) Run() {
	ctx := context.Background()
	log.Infof("[MoveCouriersJob] triggered")
	command := &commands.MoveCouriersCommand{}
	if err := j.moveCouriersCommandHandler.Handle(ctx, command); err != nil {
		log.Error(err)
	}
}
