package campaign

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSQLQueue,
	NewSvc,

	NewRouter,
	NewCreateDraftCampaignRoute,
	NewPublishCampaignRoute,
	NewTestCampaignRoute,
)
