package campaign

import (
	"context"
	"encoding/json"
	"html/template"
	"strings"
	"time"

	"github.com/tusharsoni/copper/cerror"

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
			go RunMailer(context.Background(), p)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// todo: send signal to mailer to stop processing new tasks
			return nil
		},
	})
}

func RunMailer(ctx context.Context, p MailerParams) {
	const (
		runTime   = 5 * time.Minute
		rateLimit = 50 * time.Millisecond
	)

	var limiter = time.Tick(rateLimit)

	for {
		select {
		case <-limiter:
			task, err := p.Queue.NextSendTask(ctx)
			if err != nil {
				p.Logger.Error("Failed to get next send task, exiting..", err)
				return
			}

			ctx, cancel := context.WithTimeout(ctx, runTime)

			err = runSendTask(ctx, p.Mailer, task)
			if err != nil {
				p.Logger.Error("Failed to run send task", err)
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
			}

			cancel()
		}
	}
}

func runSendTask(ctx context.Context, mailer cmailer.Mailer, task *SendTask) error {
	var (
		params    = make(map[string]interface{})
		emailBody strings.Builder
	)

	tmpl, err := template.New(task.UUID).Parse(task.HTMLBody)
	if err != nil {
		return cerror.New(err, "failed to parse html body", nil)
	}

	err = json.Unmarshal([]byte(task.Params), &params)
	if err != nil {
		return cerror.New(err, "failed to parse params", nil)
	}

	err = tmpl.Execute(&emailBody, params)
	if err != nil {
		return cerror.New(err, "failed to execute email template", nil)
	}

	_, err = mailer.SendHTML(ctx,
		task.FromEmail,
		task.ToEmail,
		task.Subject,
		emailBody.String(),
	)
	if err != nil {
		return cerror.New(err, "failed to send email", nil)
	}

	return nil
}
