package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tinh-tinh/tinhtinh/core"
)

type Service struct {
	Alg        jwt.SigningMethod
	Secret     *string
	PrivateKey *string
	PublicKey  *string
	Exp        time.Duration
	IgnoreExp  bool
}

const JWT core.Provide = "JWT"

func Register(opt Service) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		provider := core.NewProvider(module)
		provider.Set(JWT, module)

		return module
	}
}
func (j *Service) GenerateToken(payload jwt.MapClaims) (string, error) {
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(j.Exp).Unix()

	claims := jwt.NewWithClaims(j.Alg, payload)

	var key interface{}
	if j.Alg == jwt.SigningMethodRS256 {
		decodedPrivateKey, err := base64.StdEncoding.DecodeString(*j.PrivateKey)
		if err != nil {
			return "", fmt.Errorf("could not decode key: %w", err)
		}
		key, err = jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
		if err != nil {
			return "", fmt.Errorf("could not parse key: %w", err)
		}
	} else {
		key = j.Secret
	}
	token, err := claims.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *Service) VerifyToken(token string) (interface{}, error) {
	var key interface{}
	if j.Alg == jwt.SigningMethodRS256 {
		decodedPublicKey, err := base64.StdEncoding.DecodeString(*j.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("could not decode: %w", err)
		}
		key, err = jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
		if err != nil {
			return "", fmt.Errorf("validate: parse key: %w", err)
		}
	} else {
		key = j.Secret
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

	if !j.IgnoreExp && time.Now().After(exp) {
		return nil, fmt.Errorf("validate: token expired")
	}
	return claims, nil
}
