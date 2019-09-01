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

// ValidateInit implements validation for a Bet to ensure all IDs exist before
// even trying to add a bet.
func (b Bet) ValidateInit() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.BetterID, validation.Required, validation.Min(1)),
		validation.Field(&b.CompetitionID, validation.Required, validation.Min(1)),
		validation.Field(&b.CompetitorID, validation.Required, validation.Min(1)),
	)
}

// Validate implements validation for a Bet.
func (b Bet) Validate(minScore, maxScore, maxPlacing int) error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Score, validation.Min(minScore), validation.Max(maxScore)),
		validation.Field(&b.Placing, validation.Min(1), validation.Max(maxPlacing)),
	)
}
