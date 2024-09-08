package token

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	Generate(payload jwt.MapClaims) (string, error)
	Verify(token string) (interface{}, error)
}

type Options struct {
	Alg        jwt.SigningMethod
	Secret     string
	PrivateKey *string
	PublicKey  *string
	Exp        time.Duration
	IgnoreExp  bool
}

func NewJwt(opt Options) Jwt {
	return &JwtImp{
		Alg:        opt.Alg,
		Secret:     opt.Secret,
		PrivateKey: opt.PrivateKey,
		PublicKey:  opt.PublicKey,
		Exp:        opt.Exp,
		IgnoreExp:  opt.IgnoreExp,
	}
}

type JwtImp struct {
	Alg        jwt.SigningMethod
	Secret     string
	PrivateKey *string
	PublicKey  *string
	Exp        time.Duration
	IgnoreExp  bool
}

func (p *JwtImp) Generate(payload jwt.MapClaims) (string, error) {
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(p.Exp).Unix()

	claims := jwt.NewWithClaims(p.Alg, payload)

	var key interface{}
	if p.Alg == jwt.SigningMethodRS256 {
		decodedPrivateKey, err := base64.StdEncoding.DecodeString(*p.PrivateKey)
		if err != nil {
			return "", fmt.Errorf("could not decode key: %w", err)
		}
		key, err = jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
		if err != nil {
			return "", fmt.Errorf("could not parse key: %w", err)
		}
	} else {
		key = p.Secret
	}
	token, err := claims.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *JwtImp) Verify(token string) (interface{}, error) {
	var key interface{}
	if p.Alg == jwt.SigningMethodRS256 {
		decodedPublicKey, err := base64.StdEncoding.DecodeString(*p.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("could not decode: %w", err)
		}
		key, err = jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
		if err != nil {
			return "", fmt.Errorf("validate: parse key: %w", err)
		}
	} else {
		key = p.Secret
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	exp := claims["exp"].(time.Time)

	if !p.IgnoreExp && time.Now().After(exp) {
		return nil, fmt.Errorf("validate: token expired")
	}
	return claims, nil
}
