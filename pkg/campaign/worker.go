package campaign

import (
	"context"
	"droplist/pkg/profile"
	"encoding/json"
	"html/template"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type MailerParams struct {
	fx.In

	Svc       Svc
	Queue     Queue
	Profile   profile.Svc
	Secrets   profile.Secrets
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

	campaign, err := p.Svc.GetCampaign(ctx, task.CampaignUUID)
	if err != nil {
		return cerror.New(err, "failed to get campaign", map[string]interface{}{
			"campaignUUID": task.CampaignUUID,
		})
	}

	userProfile, err := p.Profile.GetProfile(ctx, campaign.CreatedBy)
	if err != nil {
		return cerror.New(err, "failed to get profile", map[string]interface{}{
			"userUUID": campaign.CreatedBy,
		})
	}

	awsKey, err := p.Secrets.Decrypt(ctx, userProfile.AWSSecretAccessKey)
	if err != nil {
		return cerror.New(err, "failed to decrypt aws secret access key", map[string]interface{}{
			"userUUID": campaign.CreatedBy,
		})
	}

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

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(userProfile.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			userProfile.AWSAccessKeyID,
			awsKey,
			"",
		),
	})
	if err != nil {
		return cerror.New(err, "failed to create new aws session", nil)
	}

	_, err = ses.New(awsSession).SendEmailWithContext(ctx, &ses.SendEmailInput{
		Source: aws.String(task.FromName + " <" + task.FromEmail + ">"),
		Destination: &ses.Destination{
			ToAddresses: []*string{&task.ToEmail},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(task.Subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailBody.String()),
				},
			},
		},
	})
	if err != nil {
		return cerror.New(err, "failed to send email", map[string]interface{}{
			"taskUUID": task.UUID,
		})
	}

	return nil
}
