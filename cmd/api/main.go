package main

import (
	"shoot/pkg/audience"
	"shoot/pkg/campaign"

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

		fx.Provide(NewConfig),

		audience.Fx,
		campaign.Fx,

		fx.Invoke(
			audience.RunMigrations,
			campaign.RunMigrations,
		),

		fx.Invoke(
			campaign.RegisterMailer,
		),
	)

	app.Run()
}
