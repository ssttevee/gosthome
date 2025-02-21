package cid

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	cv "github.com/gosthome/gosthome/core/configvalidation"
)

type IDConfig struct {
	ID string `yaml:"id"`
}

// ValidateWithContext implements validation.ValidatableWithContext.
func (i *IDConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i, validation.Field(&i.ID, cv.String(cv.Optional(cv.Name()))))
}

var _ cv.Validatable = (*IDConfig)(nil)
