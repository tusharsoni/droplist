package main

import (
	"os"
	"shoot/pkg/audience"
	"shoot/pkg/campaign"

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
}

func NewConfig() Config {
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
			BaseURL: os.Getenv("SHOOT_BASE_URL"),
		},
		Campaign: campaign.Config{
			BaseURL: os.Getenv("SHOOT_BASE_URL"),
		},
		AWSMailer: cmailer.AWSConfig{
			Region:          "us-east-1",
			AccessKeyId:     "AKIAIKIZY7B54QZ5M4UA",
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}
