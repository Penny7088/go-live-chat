package jwt

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"
	"lingua_exchange/internal/model"
)

var (
	UserTokenExpireTime    = 15 * time.Second // / 15分钟过期时间
	RefreshTokenExpireTime = 24 * time.Hour   // / 1天过期时间

	// cache prefix key, must end with a colon
	tokenCachePrefixKey = "access_token:"
	tokenUserKey        = "user_id:"
	authorizationKey    = "Authorization"
	refreshTokenKey     = "Refresh-Token"
)

func GenerateTokens(userID uint64) (string, string, error) {
	accessToken, accessExp, err := createToken(userID, UserTokenExpireTime)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := createToken(userID, RefreshTokenExpireTime)
	if err != nil {
		return "", "", err
	}

	err = storeAccessTokenInRedis(accessToken, userID, accessExp)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// 创建Token并返回Token字符串和有效期
func createToken(userID uint64, duration time.Duration) (string, int64, error) {
	expiration := time.Now().Add(duration).Unix()
	token, err := jwt.GenerateToken(utils.Uint64ToStr(userID))
	return token, expiration, err
}

// 存储Access Token到Redis
func storeAccessTokenInRedis(token string, userID uint64, exp int64) error {
	ctx := context.Background()
	return model.GetRedisCli().Set(ctx, tokenCachePrefixKey+token, userID, time.Until(time.Unix(exp, 0))).Err()
}

// ValidateAndRefreshTokens Token验证并无感刷新
func ValidateAndRefreshTokens(c *gin.Context) {
	accessToken := c.GetHeader(authorizationKey)
	if accessToken == "" {
		c.JSON(401, gin.H{"error": "authorizationKey token required"})
		c.Abort()
		return
	}

	// 尝试验证Access Token
	userID, err := validateToken(accessToken, c)
	if err == nil {
		c.Set(tokenUserKey, userID)
		c.Next()
		return
	}

	// 若Access Token过期则检查Refresh Token
	refreshToken := c.GetHeader(refreshTokenKey)
	if refreshToken == "" {
		c.JSON(401, gin.H{"error": "refresh token required"})
		c.Abort()
		return
	}

	userID, err = validateToken(refreshToken, c)
	if err != nil {
		c.JSON(401, gin.H{"error": "refresh token expired, please re-login"})
		c.Abort()
		return
	}

	parseUint, err := strconv.ParseUint(userID, 10, 64)
	// Access Token无效，刷新Access Token
	newAccessToken, _, err := GenerateTokens(parseUint)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate new access token"})
		c.Abort()
		return
	}

	c.Header(authorizationKey, newAccessToken)
	c.Set(tokenUserKey, userID)
	c.Next()
}

// 验证Token有效性
func validateToken(tokenString string, c *gin.Context) (string, error) {
	token, err := jwt.ParseToken(tokenString)
	if err != nil {
		logger.Warn("ParseToken error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return "", err
	}

	// 检查Redis中Token是否有效
	ctx := context.Background()
	_, err = model.GetRedisCli().Get(ctx, tokenCachePrefixKey+tokenString).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("token expired")
	} else if err != nil {
		return "", err
	}

	return token.UID, nil
}

// AuthMiddleware 自动Token验证和无感刷新中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ValidateAndRefreshTokens(c)
	}
}
