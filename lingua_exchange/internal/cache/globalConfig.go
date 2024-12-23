package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/model"
)

const (
	globalConfigKey            = "globalConfig:"
	verificationCode           = "verificationCode:"
	verificationCodeExpireTime = 5 * 60 * time.Second
	VCodeSignUpType            = "signUp:"
	VCodeForgetType            = "forget:"
)

type GlobalConfigCache interface {
	SetVerificationCode(ctx context.Context, email string, code string, codeType string) error
	GetVerificationCode(ctx context.Context, email string, codeType string) (string, error)
}

type globalConfigCache struct {
	cache cache.Cache
}

func NewGlobalConfigCache(cacheType *model.CacheType) GlobalConfigCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GlobalConfigModel{}
		})
		return &globalConfigCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GlobalConfigModel{}
		})
		return &globalConfigCache{cache: c}
	}

	return nil
}

func (g globalConfigCache) GetVerificationCodeKey(email string, codeType string) string {
	return verificationCode + codeType + email
}

func (g globalConfigCache) SetVerificationCode(ctx context.Context, email string, code string, codeType string) error {
	if email == "" || code == "" {
		return nil
	}
	codeKey := g.GetVerificationCodeKey(email, codeType)
	err := model.GetRedisCli().Set(ctx, codeKey, code, verificationCodeExpireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (g globalConfigCache) GetVerificationCode(ctx context.Context, email string, codeType string) (string, error) {
	if email == "" {
		return "", nil
	}
	codeKey := g.GetVerificationCodeKey(email, codeType)
	code, err := model.GetRedisCli().Get(ctx, codeKey).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}
