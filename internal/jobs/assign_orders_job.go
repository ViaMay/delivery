package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrdersCommandHandler commands.AssignOrdersCommandHandler
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

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	log.Infof("[AssignOrdersJob] triggered")
	command := &commands.AssignOrdersCommand{}
	if err := j.assignOrdersCommandHandler.Handle(ctx, command); err != nil {
		log.Error(err)
	}
}
