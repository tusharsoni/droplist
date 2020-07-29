package config

import (
	"github.com/BurntSushi/toml"
	"github.com/tusharsoni/copper/cerror"
)

type Config struct {
	HTTP struct {
		Port    uint
		BaseURL string
	}

	SQL struct {
		Host     string
		Port     uint
		Name     string
		User     string
		Password string
	}

	AWS struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
	}

	Secrets struct {
		Passphrase string
	}

	Auth struct {
		VerificationEmailFrom string
	}
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
