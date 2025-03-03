package jwt

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"
	"lingua_exchange/internal/model"
	"lingua_exchange/pkg/strutil"
)

var (
	UserTokenExpireTime    = 24 * time.Hour     // / 15分钟过期时间
	RefreshTokenExpireTime = 7 * 24 * time.Hour // / 1天过期时间

	// cache prefix key, must end with a colon
	tokenCachePrefixKey = "access_token:"
	tokenUserKey        = "user_id:"
	authorizationKey    = "Authorization"
	refreshTokenKey     = "Refresh-Token"
	env                 = "env"
	platform            = "platform"
	deviceToken         = "deviceToken"
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

func ValidateWSToken(c *gin.Context) {
	wsURL := c.Request.URL.String()

	// 解析 URL
	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid URL"})
		return
	}
	queryParams := parsedURL.Query()
	authorization := queryParams.Get(authorizationKey)

	uid, err := validateToken(authorization, c)
	if err != nil {
		c.JSON(401, gin.H{"error": "refresh token expired, please re-login"})
		c.Abort()
		return
	}
	// 将授权令牌放入上下文
	c.Set(tokenUserKey, uid)

	// 调用下一个中间件/处理程序
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

// AuthMiddleware 自动Token验证和无感刷新中间件 开发环境
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ValidateAndRefreshTokens(c)
	}
}

func AuthWSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ValidateWSToken(c)
	}
}

func ParseWSUrlUserId(c *gin.Context) string {
	wsURL := c.Request.URL.String()

	// 解析 URL
	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid URL"})
		return ""
	}
	queryParams := parsedURL.Query()
	userId := queryParams.Get(tokenUserKey)
	return userId

}

func HeaderDevMode(c *gin.Context) (string, error) {
	env := c.Request.Header.Get(env)
	if env == "" {
		return "", errors.New("env cannot be empty")
	} else if env == "release" {
		return "release", nil
	} else {
		return "dev", nil
	}
}

func HeaderPlatform(c *gin.Context) string {
	platform := c.Request.Header.Get(platform)
	if platform == "" {
		return ""
	} else {
		return platform
	}
}

func HeaderDeviceToken(c *gin.Context) string {
	deviceToken := c.Request.Header.Get(deviceToken)
	if deviceToken == "" {
		return ""
	} else {
		return deviceToken
	}
}

func HeaderObtainUID(c *gin.Context) (int, error) {
	token := c.Request.Header.Get(authorizationKey)
	if token != "" {
		claims, err := jwt.ParseToken(token)
		if err == nil {
			return convertUID(claims.UID)
		}
	}
	return 0, errors.New("token expired")
}

func convertUID(uid string) (int, error) {
	toInt, err := strutil.StringToInt(uid)
	if err != nil {
		return 0, err
	}
	return toInt, err
}
