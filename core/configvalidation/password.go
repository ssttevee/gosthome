package cv

import (
	"bytes"
	"context"
	"encoding"
	"log/slog"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash []byte
}

// Validate implements validation.Validatable.
func (n *Password) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, n, validation.Field(&n.hash, validation.Required))
}

func (n *Password) Equal(other *Password) bool {
	nv := n.Valid()
	ov := other.Valid()
	if nv && ov {
		return bytes.Equal(n.hash, other.hash)
	}
	return !nv && !ov
}

func (n *Password) Valid() bool {
	if n == nil {
		return false
	}
	if len(n.hash) == 0 {
		return false
	}
	return true
}

func (n *Password) MarshalText() ([]byte, error) {
	return slices.Clone(n.hash), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *Password) UnmarshalText(text []byte) (err error) {
	_, err = bcrypt.Cost(text)
	if err != nil {
		n.hash, err = bcrypt.GenerateFromPassword(text, 10)
		if err != nil {
			return err
		}
		slog.Warn("Dont use plaintext password in config, please, store this password hash", "hash", string(n.hash))
	} else {
		n.hash = text
	}
	return nil
}

func ParsePassword(psk string) (*Password, error) {
	r := &Password{}
	err := r.UnmarshalText([]byte(psk))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (n *Password) Check(text string) bool {
	err := bcrypt.CompareHashAndPassword(n.hash, []byte(text))
	return err == nil
}

var _ encoding.TextUnmarshaler = (*Password)(nil)
var _ Validatable = (*Password)(nil)
