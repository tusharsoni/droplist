package main

import (
	"droplist/pkg/audience"
	"droplist/pkg/campaign"
	"droplist/pkg/profile"
	"net/url"
	"os"

	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	HTTP      chttp.Config
	SQL       csql.Config
	Audience  audience.Config
	Campaign  campaign.Config
	AWSMailer cmailer.AWSConfig
	Profile   profile.Config
}

func NewConfig() (Config, error) {
	baseURL, err := url.Parse(os.Getenv("SHOOT_BASE_URL"))
	if err != nil {
		return Config{}, cerror.New(err, "failed to parse base url", map[string]interface{}{
			"url": os.Getenv("SHOOT_BASE_URL"),
		})
	}

	return Config{
		HTTP: chttp.Config{
			Port:       9712,
			HealthPath: "/api/health",
		},
		SQL: csql.Config{
			Host: "localhost",
			Port: 5432,
			Name: "shoot",
			User: "postgres",
		},
		Audience: audience.Config{
			BaseURL: baseURL,
		},
		Campaign: campaign.Config{
			BaseURL: baseURL,
		},
		AWSMailer: cmailer.AWSConfig{
			Region:          "us-east-1",
			AccessKeyId:     "AKIAIKIZY7B54QZ5M4UA",
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		Profile: profile.Config{
			Passphrase: os.Getenv("DROPLIST_PASSPHRASE"),
		},
	}, nil
}
