package campaign

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSQLQueue,
	NewSvc,

	NewRouter,
	NewGetCampaignRoute,
	NewCreateDraftCampaignRoute,
	NewPublishCampaignRoute,
	NewTestCampaignRoute,
	NewOpenEventImageRoute,
	NewClickEventRoute,
)
