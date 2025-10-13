package handler

import (
	"net/http"
	"strings"
	"time"
	"user-auth-system/internal/app/ds"
	"user-auth-system/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) LoginUserAPI(ctx *gin.Context) {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	// Проверяем тело запроса
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Получаем пользователя по логину
	user, err := h.Repository.GetUserByLogin(body.Login)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
		return
	}

	// Проверяем пароль
	if user.Password != generateHashString(body.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "неверный пароль"})
		return
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(h.Config.JWT.AccessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bitop-admin",
		},
		UserID: user.ID,
		Role:   user.Role,
	})

	strToken, err := token.SignedString([]byte(h.Config.JWT.AccessSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сгенерировать токен"})
		return
	}

	// Возвращаем JSON с токеном и данными пользователя
	ctx.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"login":       user.Login,
		"role":        user.Role,
		"accessToken": strToken,
	})
}

// LogoutUserAPI сохраняет JWT в блеклисте Redis
func (h *Handler) LogoutUserAPI(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwtStr := authHeader[len(prefix):]

	// Парсинг токена для проверки
	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.AccessSecret), nil
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Записываем в блеклист redis с TTL равным оставшемуся сроку жизни токена (берём ExpiresIn из конфига)
	if h.Redis == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), jwtStr, h.Config.JWT.AccessTokenTTL); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) WithAuthCheck(allowedRoles ...role.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		jwtStr := authHeader[len(prefix):]

		// Проверяем блеклист Redis
		if h.Redis != nil {
			if err := h.Redis.CheckJWTInBlacklist(ctx.Request.Context(), jwtStr); err == nil {
				// токен в блеклисте
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// Разбираем токен
		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.AccessSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Проверяем роль, если переданы allowedRoles
		if len(allowedRoles) > 0 {
			allowed := false
			for _, r := range allowedRoles {
				if claims.Role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// Сохраняем информацию о пользователе в контексте
		ctx.Set("userID", claims.UserID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
