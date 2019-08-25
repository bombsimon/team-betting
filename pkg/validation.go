package pkg

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate implements validation for a Competition.
func (c Competition) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
	)
}

// Validate implements validation for a Competitor.
func (c Competitor) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
	)
}

// Validate implements validation for a Better.
func (b Better) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Email, validation.Required, is.Email),
	)
}

// Validate implements validation for a Bet.
func (b Bet) Validate(minScore, maxScore int) error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.BetterID, validation.Required, validation.Min(1)),
		validation.Field(&b.CompetitionID, validation.Required, validation.Min(1)),
		validation.Field(&b.CompetitorID, validation.Required, validation.Min(1)),
		validation.Field(&b.Score, validation.Min(minScore), validation.Max(maxScore)),
	)
}
