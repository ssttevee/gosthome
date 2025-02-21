package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	cv "github.com/gosthome/gosthome/core/configvalidation"
)

type Config struct {
	cid.IDConfig
	component.ConfigOf[Server, *Server]
	Address    string           `yaml:"address"`
	Port       uint16           `yaml:"port"`
	Password   *cv.Password     `yaml:"password"`
	Encryption ConfigEncryption `yaml:"encryption"`
}

func NewConfig() *Config {
	return &Config{
		Port: 6053,
	}
}

func (c *Config) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		c.IDConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(
			ctx, c,
			validation.Field(&c.Address),
		),
	)
}

var _ cv.Validatable = (*Config)(nil)

type ConfigEncryption struct {
	Key *frameshakers.ConfigNoisePSK `yaml:"key"`
}

// Validate implements validation.Validatable.
func (c *ConfigEncryption) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx, c,
		validation.Field(&c.Key),
	)
}

var _ cv.Validatable = (*ConfigEncryption)(nil)
