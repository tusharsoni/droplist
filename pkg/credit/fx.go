package credit

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,
	NewSQLRepo,

	NewRouter,
	NewGetProductsRoute,
	NewPurchasePackRoute,
	NewGetValidPacksRoute,
	NewCompletePackPurchaseRoute,
)
