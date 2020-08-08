package main

import (
	"droplist/pkg/audience"
	"droplist/pkg/campaign"
	"droplist/pkg/content"
	"droplist/pkg/credit"
	"droplist/pkg/profile"
	"droplist/pkg/web"

	"github.com/tusharsoni/copper/cauth"
	cauthemailotp "github.com/tusharsoni/copper/cauth/emailotp"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,
		csql.Fx,
		cmailer.AWSFx,
		cauth.Fx,
		cauthemailotp.Fx,

		fx.Provide(NewConfig),

		audience.Fx,
		campaign.Fx,
		content.Fx,
		profile.Fx,
		web.Fx,
		credit.Fx,

		fx.Invoke(
			cauth.RunMigrations,
			cauthemailotp.RunMigrations,
			audience.RunMigrations,
			campaign.RunMigrations,
			content.RunMigrations,
			profile.RunMigrations,
			credit.RunMigrations,
		),

		fx.Invoke(
			campaign.RegisterMailer,
		),
	)

	app.Run()
}
