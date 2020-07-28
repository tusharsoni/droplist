package main

import (
	"droplist/pkg/audience"
	"droplist/pkg/campaign"
	"droplist/pkg/content"
	"droplist/pkg/profile"
	"droplist/pkg/web"

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
		cmailer.AWSFx,
		cauth.Fx,
		cauthanonymous.Fx,

		fx.Provide(NewConfig),

		audience.Fx,
		campaign.Fx,
		content.Fx,
		profile.Fx,
		web.Fx,

		fx.Invoke(
			cauth.RunMigrations,
			audience.RunMigrations,
			campaign.RunMigrations,
			content.RunMigrations,
			profile.RunMigrations,
		),

		fx.Invoke(
			campaign.RegisterMailer,
		),
	)

	app.Run()
}
