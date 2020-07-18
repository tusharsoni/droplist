package main

import (
	"shoot/pkg/audience"
	"shoot/pkg/campaign"
	"shoot/pkg/content"

	"github.com/tusharsoni/copper/cauth"
	cauthanonymous "github.com/tusharsoni/copper/cauth/anonymous"

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
		cmailer.LoggerFX,
		cauth.Fx,
		cauthanonymous.Fx,

		fx.Provide(NewConfig),

		audience.Fx,
		campaign.Fx,
		content.Fx,

		fx.Invoke(
			cauth.RunMigrations,
			audience.RunMigrations,
			campaign.RunMigrations,
			content.RunMigrations,
		),

		fx.Invoke(
			campaign.RegisterMailer,
		),
	)

	app.Run()
}
