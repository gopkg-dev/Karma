package jwtx

import (
	"context"
	"time"

	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/gopkg-dev/karma/errors"
)

// Auth 接口定义了 JWT 认证的方法
type Auther interface {
	// GenerateToken 生成一个包含给定 subject 的 JWT（JSON Web Token）
	GenerateToken(ctx context.Context, subject string) (TokenInfo, error)
	// DestroyToken 使 token 失效，从 token 存储中移除
	DestroyToken(ctx context.Context, accessToken string) error
	// ParseSubject 从给定的 access token 中解析出 subject（或用户标识）
	ParseSubject(ctx context.Context, accessToken string) (string, error)
	// Release 释放 JWTAuth 实例持有的任何资源
	Release(ctx context.Context) error
}

const defaultKey = "E4N6B7H9A2R5S3T1U8"

var (
	ErrMissingJwtToken        = errors.Unauthorized("JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized("keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized("Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("Wrong signing method")
)

type options struct {
	signingMethod jwtV4.SigningMethod
	signingKey    []byte
	keyFunc       jwtV4.Keyfunc
	expired       int
	tokenType     string
}

type Option func(*options)

func SetSigningMethod(method jwtV4.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func SetSigningKey(key string) Option {
	return func(o *options) {
		o.signingKey = []byte(key)
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func New(store Store, opts ...Option) Auther {
	o := options{
		tokenType:     "Bearer",
		expired:       7200,
		signingMethod: jwtV4.SigningMethodHS512,
		signingKey:    []byte(defaultKey),
	}
	for _, opt := range opts {
		opt(&o)
	}
	o.keyFunc = func(t *jwtV4.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtV4.SigningMethodHMAC); !ok {
			return nil, ErrMissingKeyFunc
		}
		return o.signingKey, nil
	}
	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

type JWTAuth struct {
	opts  *options
	store Store
}

func (a *JWTAuth) GenerateToken(ctx context.Context, subject string) (TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second)
	token := jwtV4.NewWithClaims(a.opts.signingMethod, &jwtV4.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwtV4.NewNumericDate(expiresAt),
		NotBefore: jwtV4.NewNumericDate(now),
		IssuedAt:  jwtV4.NewNumericDate(now),
		ID:        uuid.NewString(),
	})
	jwtToken, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}
	info := &tokenInfo{
		ExpiresAt:   expiresAt.Unix(),
		TokenType:   a.opts.tokenType,
		AccessToken: jwtToken,
	}
	return info, nil
}

func (a *JWTAuth) parseToken(jwtToken string) (*jwtV4.RegisteredClaims, error) {
	token, err := jwtV4.ParseWithClaims(jwtToken, &jwtV4.RegisteredClaims{}, a.opts.keyFunc)
	if err != nil {
		var ve *jwtV4.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwtV4.ValidationErrorMalformed != 0 {
				return nil, ErrTokenInvalid
			} else if ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0 {
				return nil, ErrTokenExpired
			} else {
				return nil, ErrTokenParseFail
			}
		}
		return nil, err
	} else if !token.Valid {
		return nil, ErrTokenInvalid
	} else if token.Method != a.opts.signingMethod {
		return nil, ErrUnSupportSigningMethod
	}
	return token.Claims.(*jwtV4.RegisteredClaims), nil
}

func (a *JWTAuth) callStore(fn func(Store) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

func (a *JWTAuth) DestroyToken(ctx context.Context, jwtToken string) error {
	claims, err := a.parseToken(jwtToken)
	if err != nil {
		return err
	}

	return a.callStore(func(store Store) error {
		expired := claims.ExpiresAt.Sub(time.Now())
		return store.Set(ctx, jwtToken, expired)
	})
}

func (a *JWTAuth) ParseSubject(ctx context.Context, jwtToken string) (string, error) {
	if jwtToken == "" {
		return "", ErrMissingJwtToken
	}

	claims, err := a.parseToken(jwtToken)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(store Store) error {
		if exists, err := store.Check(ctx, jwtToken); err != nil {
			return err
		} else if exists {
			return ErrTokenInvalid
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (a *JWTAuth) Release(ctx context.Context) error {
	return a.callStore(func(store Store) error {
		return store.Close(ctx)
	})
}
