package campaign

import (
	"context"
	"encoding/json"
	"html/template"
	"strings"
	"time"

	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"go.uber.org/fx"
)

type MailerParams struct {
	fx.In

	Svc       Svc
	Queue     Queue
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
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), runTime)
			contactParams = make(map[string]interface{})
			emailBody     strings.Builder
		)

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

		handleTaskErr := func(log string, err error) {
			p.Logger.Error(log, err)
			err = p.Svc.CompleteSendTask(ctx, task.UUID, SendTaskStatusFailed)
			if err != nil {
				p.Logger.Error("Failed to mark task as failed", err)
			}
			cancel()
		}

		tmpl, err := template.New(task.UUID).Parse(task.HTMLBody)
		if err != nil {
			handleTaskErr("Failed to parse HTML body", err)
			continue
		}

		err = json.Unmarshal([]byte(task.ContactParams), &contactParams)
		if err != nil {
			handleTaskErr("Failed to parse contact params", err)
			continue
		}

		params := map[string]interface{}{
			"Contact": contactParams,
		}

		err = tmpl.Execute(&emailBody, params)
		if err != nil {
			handleTaskErr("Failed to execute email template", err)
			continue
		}

		_, err = p.Mailer.SendHTML(
			task.FromEmail,
			task.ToEmail,
			task.Subject,
			emailBody.String(),
		)
		if err != nil {
			handleTaskErr("Failed to send email", err)
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
