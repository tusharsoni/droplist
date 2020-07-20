package campaign

import (
	"context"
	"encoding/json"
	"html/template"
	"net/url"
	"path"
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
	Config    Config
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
		rateLimit = 100 * time.Millisecond
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

			err = runSendTask(ctx, p, task)
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

func runSendTask(ctx context.Context, p MailerParams, task *SendTask) error {
	var (
		params    = make(map[string]interface{})
		emailBody strings.Builder
	)

	tmpl, err := template.New(task.UUID).Funcs(template.FuncMap{
		"trackURL": func(redirectTo string) string {
			imgURL, _ := url.Parse(path.Join("/api/campaigns/", task.CampaignUUID, "/events/", task.ContactUUID, "/click"))

			imgURLQuery := imgURL.Query()
			imgURLQuery.Add("url", redirectTo)

			imgURL.RawQuery = imgURLQuery.Encode()

			return p.Config.BaseURL.ResolveReference(imgURL).String()
		},
	}).Parse(task.HTMLBody)
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

	_, err = p.Mailer.SendHTML(ctx,
		task.FromName+" <"+task.FromEmail+">",
		task.ToEmail,
		task.Subject,
		emailBody.String(),
	)
	if err != nil {
		return cerror.New(err, "failed to send email", nil)
	}

	return nil
}
