package campaign

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewRouter,
	NewCreateDraftCampaignRoute,
)
