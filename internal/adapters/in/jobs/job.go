package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

// Интерфейсная проверка
var _ cron.Job = &AssignOrdersJob{}
var _ cron.Job = &MoveCouriersJob{}

type AssignOrdersJob struct {
	assignOrdersCommandHandler commands.AssignOrdersCommandHandler
}

type MoveCouriersJob struct {
	moveCouriersCommandHandler commands.MoveCouriersCommandHandler
}

func NewAssignOrdersJob(
	assignOrdersCommandHandler commands.AssignOrdersCommandHandler,
) (*AssignOrdersJob, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("AssignOrdersCommandHandler")
	}
	return &AssignOrdersJob{
		assignOrdersCommandHandler: assignOrdersCommandHandler,
	}, nil
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

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	log.Infof("[AssignOrdersJob] triggered")
	command := &commands.AssignOrdersCommand{}
	if err := j.assignOrdersCommandHandler.Handle(ctx, command); err != nil {
		log.Error(err)
	}
}

func (j *MoveCouriersJob) Run() {
	ctx := context.Background()
	log.Infof("[MoveCouriersJob] triggered")
	command := &commands.MoveCouriersCommand{}
	if err := j.moveCouriersCommandHandler.Handle(ctx, command); err != nil {
		log.Error(err)
	}
}
