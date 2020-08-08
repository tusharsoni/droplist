package main

import (
	"droplist/pkg/audience"
	"droplist/pkg/campaign"
	"droplist/pkg/config"
	"droplist/pkg/credit"
	"droplist/pkg/profile"
	"flag"
	"net/url"
	"time"

	cauthemailotp "github.com/tusharsoni/copper/cauth/emailotp"

	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	HTTP         chttp.Config
	SQL          csql.Config
	Audience     audience.Config
	Campaign     campaign.Config
	AWSMailer    cmailer.AWSConfig
	Profile      profile.Config
	AuthEmailOTP cauthemailotp.Config
	Credit       credit.Config
}

func NewConfig() (Config, error) {
	configPath := flag.String("config", "./config.toml", "Path to the config file")
	flag.Parse()

	appConfig, err := config.Read(*configPath)
	if err != nil {
		return Config{}, cerror.New(err, "failed to read config", map[string]interface{}{
			"path": configPath,
		})
	}

	baseURL, err := url.Parse(appConfig.HTTP.BaseURL)
	if err != nil {
		return Config{}, cerror.New(err, "failed to parse base url", map[string]interface{}{
			"url": appConfig.HTTP.BaseURL,
		})
	}

	authEmailOTP := cauthemailotp.GetDefaultConfig()
	authEmailOTP.VerificationEmail.From = appConfig.Auth.VerificationEmailFrom

	creditProductsByID := make(map[string]credit.ProductConfig)
	for _, pc := range appConfig.Credit.Products {
		creditProduct := credit.ProductConfig{
			Description: pc.Description,
			UseLimit:    pc.UseLimit,
			PriceUSD:    pc.PriceUSD,
		}

		if pc.Duration != nil {
			duration, err := time.ParseDuration(*pc.Duration)
			if err != nil {
				return Config{}, cerror.New(err, "failed to parse credit product duration", map[string]interface{}{
					"duration": pc.Duration,
				})
			}

			creditProduct.Duration = &duration
		}

		creditProductsByID[pc.ID] = creditProduct
	}

	return Config{
		HTTP: chttp.Config{
			Port:       appConfig.HTTP.Port,
			HealthPath: "/api/health",
		},
		SQL: csql.Config{
			Host:     appConfig.SQL.Host,
			Port:     appConfig.SQL.Port,
			Name:     appConfig.SQL.Name,
			User:     appConfig.SQL.User,
			Password: appConfig.SQL.Password,
		},
		Audience: audience.Config{
			BaseURL: baseURL,
		},
		Campaign: campaign.Config{
			BaseURL: baseURL,
		},
		AWSMailer: cmailer.AWSConfig{
			Region:          appConfig.AWS.Region,
			AccessKeyId:     appConfig.AWS.AccessKeyID,
			SecretAccessKey: appConfig.AWS.SecretAccessKey,
		},
		Profile: profile.Config{
			Passphrase: appConfig.Secrets.Passphrase,
		},
		AuthEmailOTP: authEmailOTP,
		Credit: credit.Config{
			Enabled:   appConfig.Credit.Enabled,
			StripeKey: appConfig.Credit.StripeKey,
			Products:  creditProductsByID,
		},
	}, nil
}
