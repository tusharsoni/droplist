package credit

import (
	"droplist/pkg/ptr"
	"encoding/json"
	"time"
)

type Config struct {
	Enabled  bool
	Stripe   StripeConfig
	Products []ProductConfig
}

type StripeConfig struct {
	PublicKey string
	SecretKey string
}

type ProductConfig struct {
	ID          string
	Description string
	UseLimit    *int64
	Duration    *time.Duration
	PriceUSD    int64
}

func (p ProductConfig) MarshalJSON() ([]byte, error) {
	var j struct {
		ID          string
		Description string
		UseLimit    *int64
		DurationMS  *int64
		PriceUSD    int64
	}

	j.ID = p.ID
	j.Description = p.Description
	j.UseLimit = p.UseLimit
	j.PriceUSD = p.PriceUSD

	if p.Duration != nil {
		j.DurationMS = ptr.Int64(p.Duration.Milliseconds())
	}

	return json.Marshal(j)
}
