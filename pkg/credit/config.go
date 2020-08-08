package credit

import "time"

type Config struct {
	Enabled   bool
	StripeKey string
	Products  map[string]ProductConfig
}

type ProductConfig struct {
	Description string
	UseLimit    *int64
	Duration    *time.Duration
	PriceUSD    int64
}
