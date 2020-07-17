package campaign

import (
	"context"
	"shoot/pkg/audience"
	"time"

	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"go.uber.org/fx"
)

type MailerParams struct {
	fx.In

	Svc       Svc
	Queue     Queue
	Audience  audience.Svc
	Mailer    cmailer.Mailer
	Logger    clogger.Logger
	Lifecycle fx.Lifecycle
}

func RegisterMailer(p MailerParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go RunMailer(p)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// todo: send signal to mailer to stop processing new tasks
			return nil
		},
	})
}

func RunMailer(p MailerParams) {
	const (
		runTime  = 5 * time.Minute
		waitTime = 5 * time.Second
	)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), runTime)

		task, err := p.Queue.NextSendTask(ctx)
		if err != nil {
			p.Logger.Error("Failed to get next send task, exiting..", err)
			cancel()
			return
		}

		if task == nil {
			time.Sleep(waitTime)
			cancel()
			continue
		}

		campaign, err := p.Svc.GetCampaign(ctx, task.CampaignUUID)
		if err != nil {
			p.Logger.Error("Failed to get campaign", err)
			err = p.Svc.CompleteSendTask(ctx, task.UUID, SendTaskStatusFailed)
			if err != nil {
				p.Logger.Error("Failed to mark task as failed", err)
			}
			cancel()
			continue
		}

		contact, err := p.Audience.GetContact(ctx, task.ContactUUID)
		if err != nil {
			p.Logger.Error("Failed to get contact", err)
			err = p.Svc.CompleteSendTask(ctx, task.UUID, SendTaskStatusFailed)
			if err != nil {
				p.Logger.Error("Failed to mark task as failed", err)
			}
			cancel()
			continue
		}

		_, err = p.Mailer.SendPlain(
			campaign.FromEmail,
			contact.Email,
			"Test Subject",
			"Test Body",
		)
		if err != nil {
			p.Logger.Error("Failed to send plain email", err)
			err = p.Svc.CompleteSendTask(ctx, task.UUID, SendTaskStatusFailed)
			if err != nil {
				p.Logger.Error("Failed to mark task as failed", err)
			}
			cancel()
			continue
		}

		err = p.Svc.CompleteSendTask(ctx, task.UUID, SendTaskStatusSent)
		if err != nil {
			p.Logger.Error("Failed to mark task as sent", err)
			cancel()
			continue
		}

		cancel()
	}
}
