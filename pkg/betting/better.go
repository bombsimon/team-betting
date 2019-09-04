package betting

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bombsimon/team-betting/pkg"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/pkg/errors"
)

const (
	// HMACSecret is the secret to use for JWTs
	HMACSecret = "s3cr3t"
)

// SendSignInEmail will send an email to help user sign in.
func (s *Service) SendSignInEmail(ctx context.Context, email string) error {
	var better pkg.Better

	if err := s.DB.Gorm.Where("email = ?", email).Find(&better).Error; err != nil {
		return errors.Wrapf(pkg.ErrNotFound, "could not find better with email %s", email)
	}

	better.LinkID = null.StringFrom(uuid.New().String())
	better.LinkSentAt = null.TimeFrom(time.Now())
	better.Confirmed = true

	if err := s.DB.Gorm.Save(&better).Error; err != nil {
		return errors.Wrap(err, "could not save reset link")
	}

	linkData := map[string]interface{}{
		"email":   better.Email,
		"link_id": better.LinkID.String,
		"eat":     better.LinkSentAt.Time.Add(2 * time.Hour),
	}

	jsonLinkData, _ := json.Marshal(linkData)
	b64Signin := base64.StdEncoding.EncodeToString(jsonLinkData)

	message := []string{
		"<h1>Here's your sign in link!</h1>",
		fmt.Sprintf(
			"<p><a href=\"http://localhost:3000?signin=%s\">Click here to sign in.</a></p>",
			b64Signin,
		),
		"<p>If you did not request this email, just throw it away.</p>",
		"<p>Happy betting!</p>",
	}

	return s.MailService.SendMail(&pkg.MailContent{
		From:    "no-reply@sawert.se",
		To:      email,
		Subject: "Your sign in link!",
		Body:    strings.Join(message, "\n"),
	})
}

// SignInFromEmail will parse email sign in data and return a JWT if valid.
func (s *Service) SignInFromEmail(ctx context.Context, email, hash string) (string, error) {
	var better pkg.Better

	if err := s.DB.Gorm.Where("email = ?", email).Find(&better).Error; err != nil {
		return "", errors.Wrapf(err, "could not find better with email %s", email)
	}

	// TODO: Check if hash is ok.

	return s.JWTForBetter(ctx, &better)
}

// JWTForBetter will create a JWT for the passed better.
func (s *Service) JWTForBetter(ctx context.Context, better *pkg.Better) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    better.ID,
		"name":  better.Name,
		"email": better.Email,
		"image": better.Image,
		"nbf":   time.Now(),
		"exp":   time.Now().Add(1 * time.Hour),
	})

	return token.SignedString([]byte(HMACSecret))
}

// BetterFromJWT will parse a JWT and return the better it's signed for.
func (s *Service) BetterFromJWT(ctx context.Context, tokenString string) (*pkg.Better, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(HMACSecret), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "could not parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("could not get token claims")
	}

	return &pkg.Better{
		ID:    int(claims["id"].(float64)),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
		Image: null.StringFrom(claims["image"].(string)),
	}, nil
}
