package config

import (
	"github.com/BurntSushi/toml"
	"github.com/tusharsoni/copper/cerror"
)

type Config struct {
	HTTP struct {
		Port    uint   `toml:"port"`
		BaseURL string `toml:"base_url"`
	} `toml:"http"`

	SQL struct {
		Host     string `toml:"host"`
		Port     uint   `toml:"port"`
		Name     string `toml:"name"`
		User     string `toml:"user"`
		Password string `toml:"password"`
	} `toml:"sql"`

	AWS struct {
		Region          string `toml:"region"`
		AccessKeyID     string `toml:"access_key_id"`
		SecretAccessKey string `toml:"secret_access_key"`
	}

	Secrets struct {
		Passphrase string `toml:"passphrase"`
	} `toml:"secrets"`

	Auth struct {
		VerificationEmailFrom string `toml:"verification_email_from"`
	} `toml:"auth"`

	Credit struct {
		Enabled  bool `toml:"enabled"`
		Products []struct {
			ID          string  `toml:"id"`
			Description string  `toml:"description"`
			UseLimit    *int64  `toml:"use_limit"`
			Duration    *string `toml:"duration"`
			PriceUSD    int64   `toml:"price_usd"`
		} `toml:"products"`
	} `toml:"credit"`

	Stripe struct {
		PublicKey string `toml:"public_key"`
		SecretKey string `toml:"secret_key"`
	} `toml:"stripe"`
}

func Read(path string) (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, cerror.New(err, "failed to read config", map[string]interface{}{
			"path": path,
		})
	}

	return &config, nil
}
